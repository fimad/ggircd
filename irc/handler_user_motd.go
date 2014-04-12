package irc

func (h *UserHandler) handleCmdMotd(state State, user *User, conn Connection, msg Message) Handler {
	sendMOTD(state, user)
	return h
}
