package irc

func (h *userHandler) handleCmdPing(state state, user *user, conn connection, msg message) handler {
	name := state.getConfig().Name
	conn.send(cmdPong.withPrefix(name).withParams(name).withTrailing(name))
	return h
}
