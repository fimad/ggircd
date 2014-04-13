package irc

func (h *UserHandler) handleCmdNotice(state State, user *User, conn Connection, msg Message) Handler {
	return h.handleCmdPrivMsg(state, user, conn, msg)
}
