package irc

func (h *userHandler) handleCmdNick(state state, user *user, conn connection, msg message) handler {
	if len(msg.params) != 1 {
		sendNumeric(state, user, errorNoNicknameGiven)
		return h
	}

	if !state.setNick(user, msg.params[0]) {
		sendNumeric(state, user, errorNicknameInUse, msg.params[0])
	}
	return h
}
