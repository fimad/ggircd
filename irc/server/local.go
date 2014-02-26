package server

import (
  "log"
  "fmt"
  "net"

  "github.com/fimad/ggircd/irc/config"
  "github.com/fimad/ggircd/irc/message"
)

// Local implements the Server interface and corresponds to the locally running
// irc server.
type Local struct {
  listener net.Listener
  messages chan message.Message
}

func NewLocal(cfg config.Config) Local {
  ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
  if err != nil {
    log.Fatalf("Could not create server: %v", err)
  }

  msg := make(chan message.Message)

  return Local{
    listener: ln,
    messages: msg,
  }
}

func (l *Local) Loop() {
  go l.dispatchLoop()
  l.acceptLoop()
}

// acceptLoop waits for a new connection, and calls handleConn in a new
// goroutine for each received connection.
func (l *Local) acceptLoop() {
  for {
    conn, err := l.listener.Accept()
    if err != nil {
      log.Printf("Could not accept new connection: ", err)
      continue
    }
    go handleConn(conn, l.messages)
  }
}

// handleConn is the entry point for a new connection. It handles generic
// connection setup that is required before the connection can be classified as
// a server or a client.
func handleConn(conn net.Conn, messages chan<- message.Message) {
  parser := message.NewParser(conn)
  for {
    msg, ok := parser()
    if !ok {
      log.Print("Could not parse...")
      break
    }
    messages <- msg
  }
}

// dispatchLoop continuously reads from the message channel and handles one
// message at a time.
func (l *Local) dispatchLoop() {
  for {
    msg := <-l.messages
    log.Printf("Got new message. prefix = '%s', command = '%s', params = %q",
      msg.Prefix, msg.Command, msg.Params); }
}
