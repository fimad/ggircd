package irc

func (h *userHandler) handleCmdMotd(state state, user *user, conn connection, msg message) handler {
	sendMOTD(state, user)
	return h
}
