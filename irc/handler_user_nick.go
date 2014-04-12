package irc

func (h *UserHandler) handleCmdNick(state State, user *User, conn Connection, msg Message) Handler {
	if len(msg.Params) != 1 {
		sendNumeric(state, user, ErrorNoNicknameGiven)
		return h
	}

	if !state.SetNick(user, msg.Params[0]) {
		sendNumeric(state, user, ErrorNicknameInUse, msg.Params[0])
	}
	return h
}
