package irc

import (
  "strings"
)

type Message struct {
  Prefix  string
  Command string
  Params  []string

  // The connection that this message originated from
  // Sender ...

  // The Client responsible for sending this message. Nil if not send by a
  // client.
  Client Client

  // The Server responsible for sending this message. Nil if not send by a
  // server.
  //Server Server

  // IsNewConn indicates that this message is coming from a Server or Client
  // that is not yet registered with the Local server. This field is only set by
  // handleNewConn.
  IsNewConn bool
}

// ToString serializes a Message to an IRC protocol compatible string.
func (m Message) ToString() (string, bool) {
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
