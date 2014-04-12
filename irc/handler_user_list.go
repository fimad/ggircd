package irc

import (
	"strconv"
	"strings"
)

const endOfListMessage = "End of /LIST"

func (h *UserHandler) handleCmdList(state State, user *User, conn Connection, msg Message) Handler {
	sendListElem := func(ch *Channel) {
		size := strconv.Itoa(len(ch.Users))
		sendNumericTrailing(state, user, ReplyList, ch.Topic, ch.Name, size)
	}

	if len(msg.Params) == 0 {
		state.ForChannels(sendListElem)
		sendNumericTrailing(state, user, ReplyListEnd, endOfListMessage)
		return h
	}

	channels := strings.Split(msg.Params[0], ",")

	for _, name := range channels {
		ch := state.GetChannel(name)
		if ch == nil {
			continue
		}
		sendListElem(ch)
	}

	sendNumericTrailing(state, user, ReplyListEnd, endOfListMessage)
	return h
}
