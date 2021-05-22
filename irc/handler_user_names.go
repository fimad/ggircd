package irc

import (
	"strings"
)

func (h *userHandler) handleCmdNames(state state, user *user, conn connection, msg message) handler {
	if len(msg.params) == 0 {
		h.handleCmdNamesAll(state, user, msg)
		return h
	}
	h.handleCmdNamesChannel(state, user, msg)
	return h
}

func (h *userHandler) handleCmdNamesAll(state state, user *user, msg message) {
	state.forChannels(func(ch *channel) {
		isHidden := ch.mode[channelModePrivate] || ch.mode[channelModeSecret]
		if isHidden && !ch.users[user] {
			return
		}

		sendNames(state, user, ch)
	})
}

func (h *userHandler) handleCmdNamesChannel(state state, user *user, msg message) {
	names := strings.Split(msg.params[0], ",")
	for _, name := range names {
		ch := state.getChannel(name)
		if ch == nil {
			break
		}

		isHidden := ch.mode[channelModePrivate] || ch.mode[channelModeSecret]
		if isHidden && !ch.users[user] {
			break
		}

		sendNames(state, user, ch)
	}
}
