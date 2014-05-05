package irc

import (
	"strings"
)

func (h *userHandler) handleCmdKick(state state, user *user, conn connection, msg message) handler {
	if len(msg.params) < 2 {
		sendNumeric(state, user, errorNeedMoreParams)
		return h
	}

	channelNames := strings.Split(msg.params[0], ",")
	nickNames := strings.Split(msg.params[1], ",")

	if len(channelNames) != 1 && len(nickNames) != len(channelNames) {
		sendNumeric(state, user, errorNeedMoreParams)
		return h
	}

	// A helper for kicking and sending a message.
	kickUser := func(ch *channel, nickName string) {
		if !ch.users[user] {
			sendNumeric(state, user, errorNotOnChannel, ch.name)
			return
		}

		if !ch.ops[user] {
			sendNumeric(state, user, errorChanOPrivsNeeded, ch.name)
			return
		}

		nick := state.getUser(nickName)
		if nick == nil {
			sendNumeric(state, user, errorUserNotInChannel, nickName, ch.name)
			return
		}

		if !ch.users[nick] {
			sendNumeric(state, user, errorUserNotInChannel, nick.nick, ch.name)
			return
		}

		ch.send(cmdKick.
			withPrefix(user.prefix()).
			withParams(ch.name, nick.nick))
		state.removeFromChannel(ch, nick)
	}

	// If only one channel is given, then kick all of the users given from that
	// one channel.
	if len(channelNames) == 1 {
		channel := state.getChannel(channelNames[0])
		if channel == nil {
			sendNumeric(state, user, errorNoSuchChannel, channelNames[0])
			return h
		}

		for _, nickName := range nickNames {
			kickUser(channel, nickName)
		}
		return h
	}

	// If there are equal number of nick names and channels, then step through
	// them in parallel treating each as a single kick.
	for i, nickName := range nickNames {
		channel := state.getChannel(channelNames[i])
		if channel == nil {
			sendNumeric(state, user, errorNoSuchChannel, channelNames[i])
			continue
		}
		kickUser(channel, nickName)
	}
	return h
}
