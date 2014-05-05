package irc

import (
	"strconv"
)

func (h *UserHandler) handleCmdMode(state State, user *User, _ Connection, msg Message) Handler {
	if len(msg.Params) < 1 {
		sendNumeric(state, user, ErrorNeedMoreParams)
		return h
	}

	if msg.Params[0][0] == '#' || msg.Params[0][0] == '&' {
		h.handleCmdModeChannel(state, user, msg)
	} else {
		h.handleCmdModeUser(state, user, msg)
	}
	return h
}

func (h *UserHandler) handleCmdModeChannel(state State, user *User, msg Message) {
	channel := state.GetChannel(msg.Params[0])
	if channel == nil {
		sendNumeric(state, user, ErrorNoSuchChannel, msg.Params[0])
		return
	}

	if len(msg.Params) == 1 {
		sendChannelMode(state, user, channel)
		return
	}

	if !channel.Ops[user] {
		sendNumeric(state, user, ErrorChanOPrivsNeeded, channel.Name)
		return
	}

	pos, neg, errs := ParseModeDiff(
		ChannelModes,
		ChannelPosParamModes,
		ChannelNegParamModes,
		ErrorUnknownMode,
		msg.Params[1:])

	if errs != nil {
		for _, err := range errs {
			sendNumeric(state, user, err, err.Params...)
		}
		return
	}

	// Perform a dry run first
	ok := true
	dryRun := func(curr map[string][]string) {
		for mode, values := range curr {
			for _, value := range values {
				switch mode {
				case ChannelModeOp:
					fallthrough
				case ChannelModeVoice:
					if state.GetUser(value) == nil {
						sendNumeric(state, user, ErrorNoSuchNick, value)
						ok = false
					}
				case ChannelModeUserLimit:
					_, err := strconv.Atoi(value)
					if err != nil {
						sendNumericTrailing(state, user, ErrorUnknownMode, "Not a number", value)
						ok = false
					}
					// TODO(will): Handle ban masks...
				}
			}
		}
	}

	dryRun(pos)
	dryRun(neg)
	if !ok {
		return
	}

	wetRun := func(curr map[string][]string, affinity bool) {
		for mode, values := range curr {
			for _, value := range values {
				switch mode {
				case ChannelModeOp:
					channel.Ops[state.GetUser(value)] = affinity

				case ChannelModeUserLimit:
					channel.Mode[mode] = affinity
					if !affinity {
						break
					}
					limit, _ := strconv.Atoi(value)
					channel.Limit = limit

				case ChannelModeBanMask:
					// TODO(will): Handle ban masks.

				case ChannelModeVoice:
					channel.Voices[state.GetUser(value)] = affinity

				case ChannelModeKey:
					channel.Mode[mode] = affinity
					if !affinity {
						break
					}
					channel.Key = value

				default:
					channel.Mode[mode] = affinity
				}
			}
		}
	}
	wetRun(pos, true)
	wetRun(neg, false)

	msg.Prefix = user.Prefix()
	channel.Send(msg)
}

func (h *UserHandler) handleCmdModeUser(state State, user *User, msg Message) {
	if len(msg.Params) < 2 {
		sendNumeric(state, user, ErrorNeedMoreParams)
		return
	}

	if Lowercase(msg.Params[0]) != Lowercase(user.Nick) {
		sendNumeric(state, user, ErrorUsersDontMatch)
		return
	}

	pos, neg, errs := ParseModeDiff(
		UserModesSettable,
		UserPosParamModes,
		UserNegParamModes,
		ErrorUModeUnknownFlag,
		msg.Params[1:])

	if errs != nil {
		for _, err := range errs {
			sendNumeric(state, user, err, err.Params...)
		}
		return
	}

	wetRun := func(diff map[string][]string, affinity bool) {
		for mode := range diff {
			user.Mode[mode] = affinity
		}
	}
	wetRun(pos, true)
	wetRun(neg, false)

	modeLine := ""
	for mode, isSet := range user.Mode {
		if isSet {
			modeLine = modeLine + mode
		}
	}
	sendNumeric(state, user, ReplyUModeIs, modeLine)
}
