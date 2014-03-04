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

// WithParams creates a new copy of a message with the given parameters.
func (m Message) WithParams(params ...string) Message {
  m.Params = params
  return m
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

  for i := 0; i < len(m.Params)-1; i++ {
    param := m.Params[i]
    if strings.Index(param, " ") != -1 {
      return "", false
    }
    msg += " " + param
  }
  // Always prefix the last parameter with a ':'
  msg += " :" + m.Params[len(m.Params)-1]

  msg += "\x0d\x0a"

  return msg, true
}
