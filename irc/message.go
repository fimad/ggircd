package irc

import (
  "strings"
)

type Message struct {
  Prefix  string
  Command string
  Params  []string

  // The Relay that this message originated from.
  Relay *Relay

  // Only used in Dispatcher to Relay messages. If true, the Relay will shut
  // down after sending this message.
  ShouldKill bool
}

// ToString serializes a Message to an IRC protocol compatible string.
func (m Message) ToString() (string, bool) {
  if m.Command == "" {
    return "", false
  }

  var msg string
  if len(m.Prefix) > 0 {
    msg = ":" + m.Prefix + " "
  }

  msg += m.Command

  for i := 0; i < len(m.Params); i++ {
    param := m.Params[i]
    if strings.Index(param, " ") == -1 {
      msg += " " + param
      continue
    }

    // Only the last parameter is allowed to contain a space.
    if i != len(m.Params)-1 {
      return "", false
    }
    msg += " :" + param
  }

  msg += "\x0d\x0a"

  return msg, true
}
