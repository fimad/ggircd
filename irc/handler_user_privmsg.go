package irc

func (h *UserHandler) handleCmdPrivMsg(state State, user *User, conn Connection, msg Message) Handler {
	if len(msg.Params) < 1 {
		sendNumeric(state, user, ErrorNoRecipient)
		return h
	}

	if msg.Trailing == "" {
		sendNumeric(state, user, ErrorNoTextToSend)
		return h
	}

	target := msg.Params[0]
	msg.Prefix = user.Prefix()

	if target[0] != '#' && target[0] != '&' {
		targetUser := state.GetUser(target)
		if targetUser == nil {
			sendNumeric(state, user, ErrorNoSuchNick, target)
			return h
		}
		targetUser.Send(msg)
		return h
	}

	channel := state.GetChannel(target)

	if channel == nil || !channel.CanPrivMsg(user) {
		sendNumeric(state, user, ErrorCannotSendToChan, target)
		return h
	}

	channel.ForUsers(func (u *User) {
		if u != user {
			u.Send(msg)
		}
	})
	return h
}
