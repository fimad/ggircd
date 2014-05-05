package irc

// sendNames sends the messages associated with a NAMES request.
func sendNames(state state, n *user, channels ...*channel) {
	for _, channel := range channels {
		channel.forUsers(func(u *user) {
			// Don't display this user if they are invisible and the querying user is
			// not in the same channel as them.
			if u.mode[userModeInvisible] && !channel.users[n] {
				return
			}

			params := make([]string, 2)

			if channel.mode[channelModeSecret] {
				params[0] = "@"
			} else if channel.mode[channelModePrivate] {
				params[0] = "*"
			} else {
				params[0] = "="
			}

			params[1] = channel.name

			nick := u.nick
			if channel.ops[u] {
				nick = "@" + nick
			} else if channel.voices[n] {
				nick = "+" + nick
			}
			sendNumericTrailing(state, n, replyNamReply, nick, params...)
		})
		sendNumericTrailing(state, n, replyEndOfNames, "End NAMES", channel.name)
	}
}
