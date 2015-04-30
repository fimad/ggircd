package irc

func (h *userHandler) handleCmdTopic(state state, user *user, conn connection, msg message) handler {
	if len(msg.params) < 1 {
		sendNumeric(state, user, errorNeedMoreParams)
		return h
	}

	name := msg.params[0]
	channel := state.getChannel(name)
	if channel == nil {
		sendNumeric(state, user, errorNoSuchChannel, name)
		return h
	}

	var topic = msg.laxTrailing(1)
	if topic == "" {
		sendTopic(state, user, channel)
		return h
	}

	if !channel.users[user] {
		sendNumeric(state, user, errorNotOnChannel, name)
		return h
	}

	if channel.mode[channelModeTopicOp] && !channel.ops[user] {
		sendNumeric(state, user, errorChanOPrivsNeeded, name)
		return h
	}

	msg.prefix = user.prefix()
	channel.topic = topic
	channel.send(msg)
	return h
}
