package irc

func (h *userHandler) handleCmdQuit(state state, user *user, conn connection, msg message) handler {
	state.removeUser(user)
	conn.kill()
	return nullHandler{}
}
