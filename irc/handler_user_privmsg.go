package irc

func (h *userHandler) handleCmdPrivMsg(state state, u *user, conn connection, msg message) handler {
	if len(msg.params) < 1 {
		sendNumeric(state, u, errorNoRecipient)
		return h
	}

	msgContents := msg.laxTrailing(1)
	if msgContents == "" {
		sendNumeric(state, u, errorNoTextToSend)
		return h
	}

	target := msg.params[0]

	// Construct a new PRIVMSG command from the users. This ensures that the
	// message is normalized and in a format that is recognizable by all IRC
	// clients.
	msgToSend := msg.
		withPrefix(u.prefix()).
		withParams(target).
		withTrailing(msgContents)

	if target[0] != '#' && target[0] != '&' {
		targetUser := state.getUser(target)
		if targetUser == nil {
			sendNumeric(state, u, errorNoSuchNick, target)
			return h
		}
		targetUser.send(msgToSend)

		if targetUser.mode[userModeAway] {
			u.send(cmdPrivMsg.
				withPrefix(targetUser.prefix()).
				withParams(u.nick).
				withTrailing(targetUser.awayMessage))
		}
		return h
	}

	channel := state.getChannel(target)

	if channel == nil || !channel.canPrivMsg(u) {
		sendNumeric(state, u, errorCannotSendToChan, target)
		return h
	}

	channel.forUsers(func(n *user) {
		if n != u {
			n.send(msgToSend)
		}
	})
	return h
}
