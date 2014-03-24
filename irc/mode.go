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
