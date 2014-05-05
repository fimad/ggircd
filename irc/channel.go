package irc

type channel struct {
	name string

	mode mode

	topic string

	limit int
	key   string

	banNick string
	banUser string
	banHost string

	users  map[*user]bool
	ops    map[*user]bool
	voices map[*user]bool

	invited map[*user]bool
}

func (ch channel) send(msg message) {
	for user := range ch.users {
		user.sink.send(msg)
	}
}

// invite marks a user as invited to the channel. This only has an effect if the
// channel is invite only.
func (ch channel) invite(user *user) {
	if !ch.mode[channelModeInvite] {
		return
	}
	ch.invited[user] = true
}

// forUsers iterates over all of the users in the channel and passes a pointer
// to each to the supplied callback.
func (ch channel) forUsers(callback func(*user)) {
	for u := range ch.users {
		callback(u)
	}
}

// isBanned takes a user and returns a boolean indicating if that user is banned
// in this channel.
func (ch channel) isBanned(user *user) bool {
	// TODO(will): Actually implement banned users.
	return false
}

// canPrivMsg returns a boolean indicating whether or not the given user has
// permission to message the channel.
func (ch channel) canPrivMsg(user *user) bool {
	if ch.mode[channelModeNoOutside] && !ch.users[user] {
		return false
	}

	if ch.mode[channelModeModerated] && !ch.voices[user] && !ch.ops[user] {
		return false
	}

	return !ch.isBanned(user)
}
