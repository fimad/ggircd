package irc

func (h *userHandler) handleCmdNotice(state state, user *user, conn connection, msg message) handler {
	return h.handleCmdPrivMsg(state, user, conn, msg)
}
