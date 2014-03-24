package irc

import (
  "strings"
)

type Message struct {
  Prefix   string
  Command  string
  Params   []string
  Trailing string
}

// WithParams creates a new copy of a message with the given parameters.
func (m Message) WithParams(params ...string) Message {
  m.Params = params
  return m
}

// WithParams creates a new copy of a message with the given parameters.
func (m Message) WithTrailing(trailing string) Message {
  m.Trailing = trailing
  return m
}

// WithPrefix creates a new copy of a message with the given prefix.
func (m Message) WithPrefix(prefix string) Message {
  m.Prefix = prefix
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

  for i := 0; i < len(m.Params); i++ {
    param := m.Params[i]
    if strings.Index(param, " ") != -1 {
      return "", false
    }
    msg += " " + param
  }

  if m.Trailing != "" {
    msg += " :" + m.Trailing
  }

  msg += "\x0d\x0a"

  return msg, true
}
