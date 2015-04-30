package irc

import (
	"strings"
)

func (h *userHandler) handleCmdPart(state state, user *user, conn connection, msg message) handler {
	if len(msg.params) == 0 {
		sendNumeric(state, user, errorNeedMoreParams)
		return h
	}

	reason := msg.laxTrailing(1)
	channels := strings.Split(msg.params[0], ",")
	for i := 0; i < len(channels); i++ {
		name := channels[i]
		channel := state.getChannel(name)

		if channel == nil {
			sendNumeric(state, user, errorNoSuchChannel, name)
			continue
		}

		if !channel.users[user] {
			sendNumeric(state, user, errorNotOnChannel, name)
			continue
		}

		state.partChannel(channel, user, reason)
	}
	return h
}
