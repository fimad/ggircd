package irc

// sendNames sends the messages associated with a NAMES request.
func sendNames(state State, user *User, channels ...*Channel) {
	for _, channel := range channels {
		channel.ForUsers(func(u *User) {
			// Don't display this user if they are invisible and the querying user is
			// not in the same channel as them.
			if u.Mode[UserModeInvisible] && !channel.Users[user] {
				return
			}

			params := make([]string, 2)

			if channel.Mode[ChannelModeSecret] {
				params[0] = "@"
			} else if channel.Mode[ChannelModePrivate] {
				params[0] = "*"
			} else {
				params[0] = "="
			}

			params[1] = channel.Name

			nick := u.Nick
			if channel.Ops[u] {
				nick = "@" + nick
			} else if channel.Voices[user] {
				nick = "+" + nick
			}
			sendNumericTrailing(state, user, ReplyNamReply, nick, params...)
		})
		sendNumericTrailing(state, user, ReplyEndOfNames, "End NAMES", channel.Name)
	}
}
