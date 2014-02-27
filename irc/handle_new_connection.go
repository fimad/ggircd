package irc

import (
  "log"
  "net"
)

const (
  statePass       = 0
  stateNickServer = 1
)

// handleNewConn is the entry point for a new connection. It handles generic
// connection setup that is required before the connection can be classified as
// a server or a client.
// TODO(will): The password message is currently treated as a NOP.
func handleNewConn(conn net.Conn, messages chan<- Message, cfg Config) {
  parser := NewMessageParser(conn)
  state := statePass
  //password := ""

  for {
    msg, ok := parser()
    if !ok {
      log.Print("Could not parse...")
      continue
    }

    switch state {
    case statePass:
      if msg.Command == "PASS" {
        if len(msg.Params) < 1 {
          SendError(conn, ErrorNeedMoreParams)
          continue
        }
        //password = msg.Params[1]
        state = stateNickServer
        continue
      }
      fallthrough  // The PASS command is optional.

    case stateNickServer:
      if msg.Command == "NICK" {
        msg.Client = &remoteClient {
          conn: conn,
          inbox: make(chan Message),
          outbox: messages,
        }
        msg.IsNewConn = true
        messages <- msg
        msg.Client.Loop()

        // Should not get here, but if it does we should break out of the loop.
        break
      }

      if msg.Command == "SERVER" {
        //TODO(will): fill out server interactions
      }
    }
  }
}
