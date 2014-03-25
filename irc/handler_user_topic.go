package irc

func (h *UserHandler) handleCmdTopic(state State, user *User, conn Connection, msg Message) Handler {
	if len(msg.Params) < 1 {
		sendNumeric(state, user, ErrorNeedMoreParams)
		return h
	}

	name := msg.Params[0]
	channel := state.GetChannel(name)
	if channel == nil {
		sendNumeric(state, user, ErrorNoSuchChannel, name)
		return h
	}

	if msg.Trailing == "" {
		sendTopic(state, user, channel)
		return h
	}

	if !channel.Users[user] {
		sendNumeric(state, user, ErrorNotOnChannel, name)
		return h
	}

	if channel.Mode[ChannelModeTopicOp] && !channel.Ops[user] {
		sendNumeric(state, user, ErrorChanOPrivsNeeded, name)
		return h
	}

	msg.Prefix = user.Prefix()
	channel.Topic = msg.Trailing
	channel.Send(msg)
	return h
}
