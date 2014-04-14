package irc

func (h *UserHandler) handleCmdWho(state State, user *User, conn Connection, msg Message) Handler {
	sendWho := func(ch *Channel, nick *User) {
		channelName := "*"
		flag := "H" // TODO(will): Update to support away.
		if ch != nil {
			channelName = ch.Name
			if ch.Ops[nick] {
				flag += "@"
			} else if ch.Voices[nick] {
				flag += "+"
			}

			// Don't bother reporting this user if the request is op only and this
			// user is not an op.
			if len(msg.Params) > 1 && msg.Params[1][0] == 'o' && !ch.Ops[nick] {
				return
			}
		}

		sendNumericTrailing(state, user,
			ReplyWhoReply,
			"0 "+nick.RealName,
			channelName,
			nick.User,
			nick.Host,
			state.GetConfig().Name,
			nick.Nick,
			flag)
	}

	// Always send an end of who message.
	defer sendNumeric(state, user, ReplyEndOfWho)

	if len(msg.Params) == 0 {
		return h
	}

	query := msg.Params[0]

	// Handle queries about specific users.
	if query[0] != '#' && query[0] != '&' {
		nick := state.GetUser(query)
		if nick != nil {
			sendWho(nil, nick)
		}
	}

	// Handle queries about specific channels
	ch := state.GetChannel(query)
	if ch == nil {
		return h
	}
	ch.ForUsers(func(nick *User) { sendWho(ch, nick) })

	return h
}
