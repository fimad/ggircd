package irc

import (
	"strings"
)

func (h *UserHandler) handleCmdNames(state State, user *User, conn Connection, msg Message) Handler {
	if len(msg.Params) == 0 {
		h.handleCmdNamesAll(state, user, msg)
		return h
	}
	h.handleCmdNamesChannel(state, user, msg)
	return h
}

func (h *UserHandler) handleCmdNamesAll(state State, user *User, msg Message) {
	state.ForChannels(func(ch *Channel) {
		if ch.Mode[ChannelModePrivate] && ch.Users[user] {
			return
		}

		if ch.Mode[ChannelModeSecret] && !ch.Users[user] {
			return
		}

		sendNames(state, user, ch)
	})
}

func (h *UserHandler) handleCmdNamesChannel(state State, user *User, msg Message) {
	names := strings.Split(msg.Params[0], ",")
	for _, name := range names {
		channel := state.GetChannel(name)
		if channel == nil {
			break
		}

		if channel.Mode[ChannelModePrivate] && !channel.Users[user] {
			break
		}

		sendNames(state, user, channel)
	}
}
