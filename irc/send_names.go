package irc

// sendNames sends the messages associated with a NAMES request.
func sendNames(state State, user *User, channels ...*Channel) {
	for _, channel := range channels {
		params := make([]string, 3)

		if channel.Mode[ChannelModeSecret] {
			params[0] = "@"
		} else if channel.Mode[ChannelModePrivate] {
			params[0] = "*"
		} else {
			params[0] = "="
		}

		params[1] = channel.Name

		for u := range channel.Users {
			nick := user.Nick
			if channel.Ops[u] {
				nick = "@" + nick
			} else if channel.Voices[user] {
				nick = "+" + nick
			}
			params[2] = nick
			sendNumeric(state, user, ReplyNamReply, params...)
		}
		sendNumericTrailing(state, user, ReplyEndOfNames, "End NAMES", channel.Name)
	}
}
