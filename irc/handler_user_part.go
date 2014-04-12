package irc

import (
	"strings"
)

func (h *UserHandler) handleCmdPart(state State, user *User, conn Connection, msg Message) Handler {
	if len(msg.Params) == 0 {
		sendNumeric(state, user, ErrorNeedMoreParams)
		return h
	}

	reason := "PARTing"
	if len(msg.Params) > 1 {
		reason = msg.Params[1]
	}

	channels := strings.Split(msg.Params[0], ",")
	for i := 0; i < len(channels); i++ {
		name := channels[i]
		channel := state.GetChannel(name)

		if channel == nil {
			sendNumeric(state, user, ErrorNoSuchChannel, name)
			continue
		}

		if !channel.Users[user] {
			sendNumeric(state, user, ErrorNotOnChannel, name)
			continue
		}

		state.PartChannel(channel, user, reason)
	}
	return h
}
