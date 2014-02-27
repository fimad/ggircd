package irc

import (
  "fmt"
  "log"
  "net"
)

// LocalServer implements the Server interface and corresponds to the locally running
// irc server.
type LocalServer struct {
  config   Config
  listener net.Listener
  messages chan Message
}

func NewLocalServer(cfg Config) LocalServer {
  ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
  if err != nil {
    log.Fatalf("Could not create server: %v", err)
  }

  msg := make(chan Message)

  return LocalServer{
    config:   cfg,
    listener: ln,
    messages: msg,
  }
}

// Loop is the entry point for the local server. This method does not return.
func (l *LocalServer) Loop() {
  go l.acceptLoop()
  l.dispatchLoop()
}

// acceptLoop waits for a new connection, and calls handleConn in a new
// goroutine for each received connection.
func (l *LocalServer) acceptLoop() {
  for {
    conn, err := l.listener.Accept()
    if err != nil {
      log.Printf("Could not accept new connection: ", err)
      continue
    }
    go handleNewConn(conn, l.messages, l.config)
  }
}

// dispatchLoop continuously reads from the message channel and handles one
// message at a time.
func (l *LocalServer) dispatchLoop() {
  for {
    msg := <-l.messages
    log.Printf("Got new message. prefix = '%s', command = '%s', params = %q",
      msg.Prefix, msg.Command, msg.Params)
  }
}
