package irc

// State represents the state of this IRC server. State is not safe for
// concurrent access.
type State interface {
	// GetConfig returns the IRC server's configuration.
	GetConfig() Config

	// ForChannels iterates over all of the registered channels and passes a
	// pointer to each to the supplied callback.
	ForChannels(func(*Channel))

	// ForUsers iterates over all of the registered users and passes a pointer to
	// each to the supplied callback.
	ForUsers(func(*User))

	// GetChannel returns a pointer to the channel struct corresponding to the
	// given channel name.
	GetChannel(string) *Channel

	// GetUser returns a pointer to the user struct corresponding to the given
	// nickname.
	GetUser(string) *User

	// NewUser creates a new user with the given nickname and the appropriate
	// default values.
	NewUser(string) *User

	// RemoveUser removes a user from this IRC server. In addition, it forces the
	// user to part from all channels that they are in.
	RemoveUser(*User)

	// SetNick updates the nickname for the given user. If there is a user with
	// the given nickname then this method does nothing and returns false.
	SetNick(*User, string) bool

	// NewChannel creates a new channel with the given name and the appropriate
	// default values.
	NewChannel(string) *Channel

	// RecycleChannel removes a channel if there are no more joined users.
	RecycleChannel(*Channel)

	// JoinChannel adds a user to a channel. It does not perform any permissions
	// checking, it only updates pointers.
	JoinChannel(*Channel, *User)

	// PartChannel removes a user from this channel. It sends a parting message to
	// all remaining members of the channel, and removes the channel if there are
	// no remaining users.
	PartChannel(*Channel, *User, string)

	// RemoveFromChannel silently removes a user from the given channel. It does
	// not send any messages to the channel or user. The channel will also be
	// reaped if there are no active users left.
	RemoveFromChannel(*Channel, *User)
}

// stateImpl is a concrete implementation of the State interface.
type stateImpl struct {
	config   Config
	channels map[string]*Channel
	users    map[string]*User
}

func NewState(config Config) State {
	return &stateImpl{
		config:   config,
		channels: make(map[string]*Channel),
		users:    make(map[string]*User),
	}
}

func (s stateImpl) GetConfig() Config {
	return s.config
}

func (s stateImpl) ForChannels(callback func(*Channel)) {
	for _, ch := range s.channels {
		callback(ch)
	}
}

func (s stateImpl) ForUsers(callback func(*User)) {
	for _, u := range s.users {
		callback(u)
	}
}

func (s stateImpl) GetChannel(name string) *Channel {
	return s.channels[Lowercase(name)]
}

func (s stateImpl) GetUser(nick string) *User {
	return s.users[Lowercase(nick)]
}

func (s *stateImpl) NewUser(nick string) *User {
	nick = Lowercase(nick)
	if s.users[nick] != nil {
		return nil
	}

	u := &User{
		Nick:     nick,
		Channels: make(map[*Channel]bool),
	}
	s.users[nick] = u
	return u
}

func (s *stateImpl) SetNick(user *User, nick string) bool {
	nick = Lowercase(nick)
	if s.users[nick] != nil {
		return false
	}

	user.ForChannels(func(ch *Channel) {
		ch.Send(CmdNick.WithPrefix(user.Prefix()).WithParams(nick))
	})

	delete(s.users, user.Nick)
	s.users[nick] = user
	user.Nick = nick
	return true
}

func (s *stateImpl) RemoveUser(user *User) {
	user.ForChannels(func(ch *Channel) {
		s.PartChannel(ch, user, "QUITing")
	})
	delete(s.users, user.Nick)
}

func (s *stateImpl) NewChannel(name string) *Channel {
	name = Lowercase(name)
	if s.channels[name] != nil {
		return nil
	}

	if name[0] != '#' && name[0] != '&' {
		return nil
	}

	ch := &Channel{
		Name:    name,
		Mode:    ParseMode(ChannelModes, s.GetConfig().DefaultChannelMode),
		Users:   make(map[*User]bool),
		Ops:     make(map[*User]bool),
		Voices:  make(map[*User]bool),
		Invited: make(map[*User]bool),
	}
	s.channels[name] = ch
	return ch
}

func (s *stateImpl) RecycleChannel(channel *Channel) {
	if channel == nil || len(channel.Users) != 0 {
		return
	}
	delete(s.channels, channel.Name)
}

func (s *stateImpl) JoinChannel(channel *Channel, user *User) {
	joinMsg := CmdJoin.WithPrefix(user.Prefix()).WithParams(channel.Name)
	channel.Send(joinMsg)

	channel.Users[user] = true
	user.Channels[channel] = true

	if len(channel.Users) == 1 {
		channel.Ops[user] = true
	}

	sendTopic(s, user, channel)
	sendNames(s, user, channel)
}

func (s *stateImpl) PartChannel(ch *Channel, user *User, reason string) {
	s.RemoveFromChannel(ch, user)
	ch.Send(CmdPart.
		WithPrefix(user.Prefix()).
		WithParams(ch.Name).
		WithTrailing(reason))
}

func (s *stateImpl) RemoveFromChannel(ch *Channel, user *User) {
	delete(user.Channels, ch)

	if !ch.Users[user] {
		return
	}

	delete(ch.Users, user)
	delete(ch.Voices, user)
	delete(ch.Ops, user)
	delete(ch.Invited, user)

	s.RecycleChannel(ch)
}
