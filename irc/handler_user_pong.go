package irc

func (h *userHandler) handleCmdPong(state state, user *user, conn connection, msg message) handler {
	h.gotPong <- struct{}{}
	return h
}
