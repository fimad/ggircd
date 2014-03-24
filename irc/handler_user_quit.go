package irc

func (h *UserHandler) handleCmdQuit(state State, user *User, conn Connection, msg Message) Handler {
	conn.Kill()
	state.RemoveUser(user)
	return NullHandler{}
}
