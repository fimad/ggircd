package irc

import (
	"log"
)

const (
	ChannelModeOp        = "o"
	ChannelModePrivate   = "p"
	ChannelModeSecret    = "s"
	ChannelModeInvite    = "i"
	ChannelModeTopicOp   = "t"
	ChannelModeNoOutside = "n"
	ChannelModeModerated = "m"
	ChannelModeUserLimit = "l"
	ChannelModeBanMask   = "b"
	ChannelModeVoice     = "v"
	ChannelModeKey       = "k"

	UserModeAway      = "a"
	UserModeInvisible = "i"
)

type Mode map[string]bool

var ChannelModes = Mode{
	ChannelModeOp:        true,
	ChannelModePrivate:   true,
	ChannelModeSecret:    true,
	ChannelModeInvite:    true,
	ChannelModeTopicOp:   true,
	ChannelModeNoOutside: true,
	ChannelModeModerated: true,
	ChannelModeUserLimit: true,
	ChannelModeBanMask:   true,
	ChannelModeVoice:     true,
	ChannelModeKey:       true,
}

var ChannelPosParamModes = Mode{
	ChannelModeOp:        true,
	ChannelModeUserLimit: true,
	ChannelModeBanMask:   true,
	ChannelModeVoice:     true,
	ChannelModeKey:       true,
}

var ChannelNegParamModes = Mode{
	ChannelModeOp:      true,
	ChannelModeBanMask: true,
	ChannelModeVoice:   true,
}

var UserModes = Mode{
	UserModeAway:      true,
	UserModeInvisible: true,
}

var UserModesSettable = Mode{
	UserModeInvisible: true,
}

var UserPosParamModes = Mode{}
var UserNegParamModes = Mode{}

// ParseMode takes a mode map containing all of the valid modes and a string
// where each character is a mode flag. It returns a mode map of containing all
// of the valid modes that were present in the given line.
func ParseMode(valid Mode, line string) Mode {
	mode := make(Mode)
	for _, r := range line {
		f := string(r)
		if valid[f] {
			mode[f] = true
		} else {
			log.Printf("Unknown flag: %s", f)
		}
	}
	return mode
}

// ParseModeDiff takes a set of valid modes, modes that take parameters and an
// array of strings and returns a mapping from modes that are present in the
// mode line to their parameters.
//
// NOTE: Modes that do not take parameters will map to the empty string.
func ParseModeDiff(valid Mode, posTakesParam Mode, negTakesParam Mode, unknownModeMessage Message, line []string) (pos map[string][]string, neg map[string][]string, errs []Message) {
	pos = make(map[string][]string)
	neg = make(map[string][]string)

	if len(line) < 1 {
		return nil, nil, []Message{ErrorNeedMoreParams}
	}

	other := neg
	otherTakes := negTakesParam
	curr := pos
	currTakes := posTakesParam
	params := line[1:]
	for _, char := range line[0] {
		mode := string(char)

		if mode == "+" {
			other = neg
			otherTakes = negTakesParam
			curr = pos
			currTakes = posTakesParam
			continue
		} else if mode == "-" {
			other = pos
			otherTakes = posTakesParam
			curr = neg
			currTakes = negTakesParam
			continue
		}

		if !valid[mode] {
			errs = append(errs, unknownModeMessage.WithParams(mode))
			continue
		}

		if currTakes == nil || !currTakes[mode] {
			curr[mode] = append(curr[mode], "")
			delete(other, mode)
			continue
		}

		if len(params) == 0 {
			errs = []Message{ErrorNeedMoreParams}
			break
		}

		if otherTakes == nil || !otherTakes[mode] {
			delete(other, mode)
		}

		curr[mode] = append(curr[mode], params[0])
		params = params[1:]
	}

	if errs != nil {
		pos, neg = nil, nil
	}

	return pos, neg, errs
}
