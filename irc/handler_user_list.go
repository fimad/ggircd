package irc

import (
	"strconv"
	"strings"
)

const endOfListMessage = "End of /LIST"

func (h *userHandler) handleCmdList(state state, user *user, conn connection, msg message) handler {
	sendListElem := func(ch *channel) {
		size := strconv.Itoa(len(ch.users))
		sendNumericTrailing(state, user, replyList, ch.topic, ch.name, size)
	}

	if len(msg.params) == 0 {
		state.forChannels(sendListElem)
		sendNumericTrailing(state, user, replyListEnd, endOfListMessage)
		return h
	}

	channels := strings.Split(msg.params[0], ",")

	for _, name := range channels {
		ch := state.getChannel(name)
		if ch == nil {
			continue
		}
		sendListElem(ch)
	}

	sendNumericTrailing(state, user, replyListEnd, endOfListMessage)
	return h
}
