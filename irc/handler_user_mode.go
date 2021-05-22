package irc

import (
	"strconv"
)

func (h *userHandler) handleCmdMode(state state, user *user, _ connection, msg message) handler {
	if len(msg.params) < 1 {
		sendNumeric(state, user, errorNeedMoreParams)
		return h
	}

	if msg.params[0][0] == '#' || msg.params[0][0] == '&' {
		h.handleCmdModeChannel(state, user, msg)
	} else {
		h.handleCmdModeUser(state, user, msg)
	}
	return h
}

func (h *userHandler) handleCmdModeChannel(state state, user *user, msg message) {
	channel := state.getChannel(msg.params[0])
	if channel == nil {
		sendNumeric(state, user, errorNoSuchChannel, msg.params[0])
		return
	}

	if len(msg.params) == 1 {
		sendChannelMode(state, user, channel)
		return
	}

	if !channel.ops[user] {
		sendNumeric(state, user, errorChanOPrivsNeeded, channel.name)
		return
	}

	pos, neg, errs := parseModeDiff(
		channelModes,
		channelPosParamModes,
		channelNegParamModes,
		errorUnknownMode,
		msg.params[1:])

	if errs != nil {
		for _, err := range errs {
			sendNumeric(state, user, err, err.params...)
		}
		return
	}

	// Perform a dry run first
	ok := true
	dryRun := func(curr map[string][]string, affinity bool) {
		for mode, values := range curr {
			for _, value := range values {
				switch mode {
				case channelModeOp:
					fallthrough
				case channelModeVoice:
					if state.getUser(value) == nil {
						sendNumeric(state, user, errorNoSuchNick, value)
						ok = false
					}
				case channelModeUserLimit:
					_, err := strconv.Atoi(value)
					if err != nil {
						sendNumericTrailing(
							state, user, errorUnknownMode, "Not a number", value)
						ok = false
					}
					// TODO(will): Handle ban masks...
				}
			}
		}
	}

	dryRun(pos, true)
	dryRun(neg, false)
	if !ok {
		return
	}

	wetRun := func(curr map[string][]string, affinity bool) {
		for mode, values := range curr {
			for _, value := range values {
				switch mode {
				case channelModeOp:
					channel.ops[state.getUser(value)] = affinity

				case channelModeUserLimit:
					channel.mode[mode] = affinity
					if !affinity {
						break
					}
					limit, _ := strconv.Atoi(value)
					channel.limit = limit

				case channelModeBanMask:
					// TODO(will): Handle ban masks.

				case channelModeVoice:
					channel.voices[state.getUser(value)] = affinity

				case channelModeKey:
					channel.mode[mode] = affinity
					if !affinity {
						break
					}
					channel.key = value

				case channelModePrivate:
					// The ordering of modes in a mode command is unspecified since they
					// are stored in a hash map. Give priority to secret mode sets if
					// there private and secret both appear in the same mode setting.
					if affinity && pos[channelModeSecret] != nil {
						break
					}
					channel.mode[mode] = affinity
					if affinity {
						channel.mode[channelModeSecret] = false
					}

				case channelModeSecret:
					channel.mode[mode] = affinity
					if affinity {
						channel.mode[channelModePrivate] = false
					}

				default:
					channel.mode[mode] = affinity
				}
			}
		}
	}
	wetRun(pos, true)
	wetRun(neg, false)

	msg.prefix = user.prefix()
	channel.send(msg)
}

func (h *userHandler) handleCmdModeUser(state state, user *user, msg message) {
	if len(msg.params) < 2 {
		sendNumeric(state, user, errorNeedMoreParams)
		return
	}

	if lowercase(msg.params[0]) != lowercase(user.nick) {
		sendNumeric(state, user, errorUsersDontMatch)
		return
	}

	pos, neg, errs := parseModeDiff(
		userModesSettable,
		userPosParamModes,
		userNegParamModes,
		errorUModeUnknownFlag,
		msg.params[1:])

	if errs != nil {
		for _, err := range errs {
			sendNumeric(state, user, err, err.params...)
		}
		return
	}

	wetRun := func(diff map[string][]string, affinity bool) {
		for mode := range diff {
			user.mode[mode] = affinity
		}
	}
	wetRun(pos, true)
	wetRun(neg, false)

	modeLine := ""
	for mode, isSet := range user.mode {
		if isSet {
			modeLine = modeLine + mode
		}
	}
	sendNumeric(state, user, replyUModeIs, modeLine)
}
