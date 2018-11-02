package irc

import (
	"github.com/prometheus/client_golang/prometheus"
)

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

func isValidNick(nick string) bool {
	return len(nick) <= 9
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

	if !isValidNick(nickLower) {
		return nil
	}

	logf(debug, "Adding new user %s", nick)

	u := &user{
		nick:     nick,
		channels: make(map[*channel]bool),
		mode:     parseMode(userModes, s.getConfig().DefaultUserMode),
	}
	s.users[nickLower] = u
	return u
}

func (s *stateImpl) setNick(user *user, nick string) bool {
	logf(debug, "User requesting nick change from %s to %s.", user.nick, nick)

	lowerNick := lowercase(nick)
	if n := s.users[lowerNick]; n != nil && n != user {
		return false
	}

	if !isValidNick(lowerNick) {
		return false
	}

	if len(user.channels) == 0 {
		// Always send the user requesting a nick change a confirmation message
		// since some clients will be left in a strange state if they do not receive
		// either a success or error message.
		user.send(cmdNick.withPrefix(user.prefix()).withParams(nick))
	} else {
		user.forChannels(func(ch *channel) {
			ch.send(cmdNick.withPrefix(user.prefix()).withParams(nick))
		})
	}

	delete(s.users, lowercase(user.nick))
	s.users[lowerNick] = user
	user.nick = nick
	return true
}

func (s *stateImpl) removeUser(user *user) {
	logf(debug, "Removing user %s", user.nick)

	user.forChannels(func(ch *channel) {
		s.partChannel(ch, user, "QUITing")
	})

	nickLower := lowercase(user.nick)
	delete(s.users, nickLower)
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
	logf(debug, "Recycling channel %+v", channel)

	if channel == nil || len(channel.users) != 0 {
		return
	}
	delete(s.channels, channel.name)
}

func (s *stateImpl) joinChannel(channel *channel, user *user) {
	// Don't add a user to a channel more than once.
	if channel.users[user] {
		return
	}

	nicks_in_channel.With(
		prometheus.Labels{
			"channel": channel.name,
		},
	).Inc()

	logf(debug, "Adding %+v to %+v", user, channel)

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

func (s *stateImpl) partChannel(channel *channel, user *user, reason string) {
	nicks_in_channel.With(
		prometheus.Labels{
			"channel": channel.name,
		},
	).Dec()

	channel.send(cmdPart.
		withPrefix(user.prefix()).
		withParams(channel.name).
		withTrailing(reason))
	s.removeFromChannel(channel, user)
}

func (s *stateImpl) removeFromChannel(channel *channel, user *user) {
	logf(debug, "Removing %+v from %+v", user, channel)

	delete(user.channels, channel)

	if !channel.users[user] {
		return
	}

	delete(channel.users, user)
	delete(channel.voices, user)
	delete(channel.ops, user)
	delete(channel.invited, user)

	s.recycleChannel(channel)
}
