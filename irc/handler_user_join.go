package irc

import (
	"strings"
)

func (h *userHandler) handleCmdJoin(state state, user *user, conn connection, msg message) handler {
	if len(msg.params) == 0 {
		sendNumeric(state, user, errorNeedMoreParams)
		return h
	}
	channels := strings.Split(msg.params[0], ",")

	var keys []string
	if len(msg.params) > 1 {
		keys = strings.Split(msg.params[1], ",")
	}

	for i := 0; i < len(channels); i++ {
		name := channels[i]
		channel := state.getChannel(name)
		if channel == nil {
			channel = state.newChannel(name)
			defer state.recycleChannel(channel)
		}

		if channel == nil {
			sendNumeric(state, user, errorNoSuchChannel, name)
			continue
		}

		if channel.mode[channelModeInvite] && !channel.invited[user] {
			sendNumeric(state, user, errorInviteOnlyChan, name)
			continue
		}

		if channel.mode[channelModeKey] &&
			(keys == nil || len(keys) <= i || keys[i] != channel.key) {
			sendNumeric(state, user, errorBadChannelKey, name)
			continue
		}

		if channel.mode[channelModeUserLimit] &&
			len(channel.users) >= channel.limit {
			sendNumeric(state, user, errorChannelIsFull, name)
			continue
		}

		if channel.isBanned(user) {
			sendNumeric(state, user, errorBannedFromChan, name)
			continue
		}

		state.joinChannel(channel, user)
	}

	return h
}
