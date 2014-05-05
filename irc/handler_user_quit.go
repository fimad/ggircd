package irc

func (h *userHandler) handleCmdQuit(state state, user *user, conn connection, msg message) handler {
	conn.kill()
	state.removeUser(user)
	return nullHandler{}
}
