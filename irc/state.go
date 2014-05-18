package irc

// state represents the state of this IRC server. State is not safe for
// concurrent access.
type state interface {
	// GetConfig returns the IRC server's configuration.
	getConfig() Config

	// forChannels iterates over all of the registered channels and passes a
	// pointer to each to the supplied callback.
	forChannels(func(*channel))

	// forUsers iterates over all of the registered users and passes a pointer to
	// each to the supplied callback.
	forUsers(func(*user))

	// getChannel returns a pointer to the channel struct corresponding to the
	// given channel name.
	getChannel(string) *channel

	// getUser returns a pointer to the user struct corresponding to the given
	// nickname.
	getUser(string) *user

	// newUser creates a new user with the given nickname and the appropriate
	// default values.
	newUser(string) *user

	// removeUser removes a user from this IRC server. In addition, it forces the
	// user to part from all channels that they are in.
	removeUser(*user)

	// setNick updates the nickname for the given user. If there is a user with
	// the given nickname then this method does nothing and returns false.
	setNick(*user, string) bool

	// newChannel creates a new channel with the given name and the appropriate
	// default values.
	newChannel(string) *channel

	// recycleChannel removes a channel if there are no more joined users.
	recycleChannel(*channel)

	// joinChannel adds a user to a channel. It does not perform any permissions
	// checking, it only updates pointers.
	joinChannel(*channel, *user)

	// partChannel removes a user from this channel. It sends a parting message to
	// all remaining members of the channel, and removes the channel if there are
	// no remaining users.
	partChannel(*channel, *user, string)

	// removeFromChannel silently removes a user from the given channel. It does
	// not send any messages to the channel or user. The channel will also be
	// reaped if there are no active users left.
	removeFromChannel(*channel, *user)
}

// stateImpl is a concrete implementation of the State interface.
type stateImpl struct {
	config   Config
	channels map[string]*channel
	users    map[string]*user
}

func newState(config Config) state {
	return &stateImpl{
		config:   config,
		channels: make(map[string]*channel),
		users:    make(map[string]*user),
	}
}

func (s stateImpl) getConfig() Config {
	return s.config
}

func (s stateImpl) forChannels(callback func(*channel)) {
	for _, ch := range s.channels {
		callback(ch)
	}
}

func (s stateImpl) forUsers(callback func(*user)) {
	for _, u := range s.users {
		callback(u)
	}
}

func (s stateImpl) getChannel(name string) *channel {
	return s.channels[lowercase(name)]
}

func (s stateImpl) getUser(nick string) *user {
	return s.users[lowercase(nick)]
}

func (s *stateImpl) newUser(nick string) *user {
	nickLower := lowercase(nick)
	if s.users[nickLower] != nil {
		return nil
	}

	u := &user{
		nick:     nick,
		channels: make(map[*channel]bool),
		mode:     parseMode(userModes, s.getConfig().DefaultUserMode),
	}
	s.users[nickLower] = u
	return u
}

func (s *stateImpl) setNick(user *user, nick string) bool {
	lowerNick := lowercase(nick)
	if n := s.users[lowerNick]; n != nil && n != user {
		return false
	}

	user.forChannels(func(ch *channel) {
		ch.send(cmdNick.withPrefix(user.prefix()).withParams(nick))
	})

	delete(s.users, user.nick)
	s.users[lowerNick] = user
	user.nick = nick
	return true
}

func (s *stateImpl) removeUser(user *user) {
	user.forChannels(func(ch *channel) {
		s.partChannel(ch, user, "QUITing")
	})
	delete(s.users, user.nick)
}

func (s *stateImpl) newChannel(name string) *channel {
	name = lowercase(name)
	if s.channels[name] != nil {
		return nil
	}

	if name[0] != '#' && name[0] != '&' {
		return nil
	}

	ch := &channel{
		name:    name,
		mode:    parseMode(channelModes, s.getConfig().DefaultChannelMode),
		users:   make(map[*user]bool),
		ops:     make(map[*user]bool),
		voices:  make(map[*user]bool),
		invited: make(map[*user]bool),
	}
	s.channels[name] = ch
	return ch
}

func (s *stateImpl) recycleChannel(channel *channel) {
	if channel == nil || len(channel.users) != 0 {
		return
	}
	delete(s.channels, channel.name)
}

func (s *stateImpl) joinChannel(channel *channel, user *user) {
	channel.users[user] = true
	user.channels[channel] = true

	if len(channel.users) == 1 {
		channel.ops[user] = true
	}

	joinMsg := cmdJoin.withPrefix(user.prefix()).withParams(channel.name)
	channel.send(joinMsg)

	sendTopic(s, user, channel)
	sendNames(s, user, channel)
}

func (s *stateImpl) partChannel(ch *channel, user *user, reason string) {
	ch.send(cmdPart.
		withPrefix(user.prefix()).
		withParams(ch.name).
		withTrailing(reason))
	s.removeFromChannel(ch, user)
}

func (s *stateImpl) removeFromChannel(ch *channel, user *user) {
	delete(user.channels, ch)

	if !ch.users[user] {
		return
	}

	delete(ch.users, user)
	delete(ch.voices, user)
	delete(ch.ops, user)
	delete(ch.invited, user)

	s.recycleChannel(ch)
}
