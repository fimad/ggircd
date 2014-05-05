package irc

const (
	channelModeOp        = "o"
	channelModePrivate   = "p"
	channelModeSecret    = "s"
	channelModeInvite    = "i"
	channelModeTopicOp   = "t"
	channelModeNoOutside = "n"
	channelModeModerated = "m"
	channelModeUserLimit = "l"
	channelModeBanMask   = "b"
	channelModeVoice     = "v"
	channelModeKey       = "k"

	userModeAway      = "a"
	userModeInvisible = "i"
)

type mode map[string]bool

var channelModes = mode{
	channelModeOp:        true,
	channelModePrivate:   true,
	channelModeSecret:    true,
	channelModeInvite:    true,
	channelModeTopicOp:   true,
	channelModeNoOutside: true,
	channelModeModerated: true,
	channelModeUserLimit: true,
	channelModeBanMask:   true,
	channelModeVoice:     true,
	channelModeKey:       true,
}

var channelPosParamModes = mode{
	channelModeOp:        true,
	channelModeUserLimit: true,
	channelModeBanMask:   true,
	channelModeVoice:     true,
	channelModeKey:       true,
}

var channelNegParamModes = mode{
	channelModeOp:      true,
	channelModeBanMask: true,
	channelModeVoice:   true,
}

var userModes = mode{
	userModeAway:      true,
	userModeInvisible: true,
}

var userModesSettable = mode{
	userModeInvisible: true,
}

var userPosParamModes = mode{}
var userNegParamModes = mode{}

// parseMode takes a mode map containing all of the valid modes and a string
// where each character is a mode flag. It returns a mode map of containing all
// of the valid modes that were present in the given line.
func parseMode(valid mode, line string) mode {
	mode := make(mode)
	for _, r := range line {
		f := string(r)
		if valid[f] {
			mode[f] = true
		} else {
			logf(debug, "Unknown flag: %s", f)
		}
	}
	return mode
}

// parseModeDiff takes a set of valid modes, modes that take parameters and an
// array of strings and returns a mapping from modes that are present in the
// mode line to their parameters.
//
// NOTE: Modes that do not take parameters will map to the empty string.
func parseModeDiff(valid mode, posTakesParam mode, negTakesParam mode, unknownModeMessage message, line []string) (pos map[string][]string, neg map[string][]string, errs []message) {
	pos = make(map[string][]string)
	neg = make(map[string][]string)

	if len(line) < 1 {
		return nil, nil, []message{errorNeedMoreParams}
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
			errs = append(errs, unknownModeMessage.withParams(mode))
			continue
		}

		if currTakes == nil || !currTakes[mode] {
			curr[mode] = append(curr[mode], "")
			delete(other, mode)
			continue
		}

		if len(params) == 0 {
			errs = []message{errorNeedMoreParams}
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
