package irc

func (h *userHandler) handleCmdWho(state state, u *user, conn connection, msg message) handler {
	sendWho := func(ch *channel, nick *user) {
		channelName := "*"
		flag := "H" // TODO(will): Update to support away.
		if ch != nil {
			channelName = ch.name
			if ch.ops[nick] {
				flag += "@"
			} else if ch.voices[nick] {
				flag += "+"
			}

			// Don't bother reporting this user if the request is op only and this
			// user is not an op.
			if len(msg.params) > 1 && msg.params[1][0] == 'o' && !ch.ops[nick] {
				return
			}
		}

		sendNumericTrailing(state, u,
			replyWhoReply,
			"0 "+nick.realName,
			channelName,
			nick.user,
			nick.host,
			state.getConfig().Name,
			nick.nick,
			flag)
	}

	// Always send an end of who message.
	defer sendNumeric(state, u, replyEndOfWho)

	if len(msg.params) == 0 {
		return h
	}

	query := msg.params[0]

	// Handle queries about specific users.
	if query[0] != '#' && query[0] != '&' {
		nick := state.getUser(query)
		if nick != nil {
			sendWho(nil, nick)
		}
	}

	// Handle queries about specific channels
	ch := state.getChannel(query)
	if ch == nil {
		return h
	}
	isSecretOrPrivate := ch.mode[channelModePrivate] || ch.mode[channelModeSecret]
	if !ch.users[u] && isSecretOrPrivate {
		return h
	}
	ch.forUsers(func(nick *user) { sendWho(ch, nick) })

	return h
}
