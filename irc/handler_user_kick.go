package irc

import (
	"strings"
)

func (h *UserHandler) handleCmdKick(state State, user *User, conn Connection, msg Message) Handler {
	if len(msg.Params) < 2 {
		sendNumeric(state, user, ErrorNeedMoreParams)
		return h
	}

	channelNames := strings.Split(msg.Params[0], ",")
	nickNames := strings.Split(msg.Params[1], ",")

	if len(channelNames) != 1 && len(nickNames) != len(channelNames) {
		sendNumeric(state, user, ErrorNeedMoreParams)
		return h
	}

	// A helper for kicking and sending a message.
	kickUser := func(ch *Channel, nickName string) {
		if !ch.Users[user] {
			sendNumeric(state, user, ErrorNotOnChannel, ch.Name)
			return
		}

		if !ch.Ops[user] {
			sendNumeric(state, user, ErrorChanOPrivsNeeded, ch.Name)
			return
		}

		nick := state.GetUser(nickName)
		if nick == nil {
			sendNumeric(state, user, ErrorUserNotInChannel, nickName, ch.Name)
			return
		}

		if !ch.Users[nick] {
			sendNumeric(state, user, ErrorUserNotInChannel, nick.Nick, ch.Name)
			return
		}

		ch.Send(CmdKick.
			WithPrefix(user.Prefix()).
			WithParams(ch.Name, nick.Nick))
		state.RemoveFromChannel(ch, nick)
	}

	// If only one channel is given, then kick all of the users given from that
	// one channel.
	if len(channelNames) == 1 {
		channel := state.GetChannel(channelNames[0])
		if channel == nil {
			sendNumeric(state, user, ErrorNoSuchChannel, channelNames[0])
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
		channel := state.GetChannel(channelNames[i])
		if channel == nil {
			sendNumeric(state, user, ErrorNoSuchChannel, channelNames[i])
			continue
		}
		kickUser(channel, nickName)
	}
	return h
}
