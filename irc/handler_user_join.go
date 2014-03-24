package irc

import (
	"strings"
)

func (h *UserHandler) handleCmdJoin(state State, user *User, conn Connection, msg Message) Handler {
	if len(msg.Params) == 0 {
		sendNumeric(state, user, ErrorNeedMoreParams)
		return h
	}
	channels := strings.Split(msg.Params[0], ",")

	var keys []string
	if len(msg.Params) > 1 {
		keys = strings.Split(msg.Params[1], ",")
	}

	for i := 0; i < len(channels); i++ {
		name := channels[i]
		channel := state.GetChannel(name)
		if channel == nil {
			channel = state.NewChannel(name)
			defer state.RecycleChannel(channel)
		}

		if channel == nil {
			sendNumeric(state, user, ErrorNoSuchChannel, name)
			continue
		}

		if channel.Mode[ChannelModeInvite] {
			sendNumeric(state, user, ErrorInviteOnlyChan, name)
			continue
		}

		if channel.Mode[ChannelModeKey] &&
			(keys == nil || len(keys) <= i || keys[i] != channel.Key) {
			sendNumeric(state, user, ErrorBadChannelKey, name)
			continue
		}

		if channel.Mode[ChannelModeUserLimit] &&
			len(channel.Users) >= channel.Limit {
			sendNumeric(state, user, ErrorChannelIsFull, name)
			continue
		}

		if channel.IsBanned(user) {
			sendNumeric(state, user, ErrorBannedFromChan, name)
			continue
		}

		state.JoinChannel(channel, user)
	}

	return h
}
