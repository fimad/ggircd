package irc

import (
  "fmt"
  "log"
  "net"
)

type Dispatcher struct {
  listener net.Listener

  Config   Config
  Inbox    chan Message

  nicks map[string]int64
  users map[string]int64

  //  servers map[int64]ServerInfo
  clients map[int64]*Client

  relayToClient map[int64][]*Client
  //relayToServer map[int64][]*Server

  nextID int64
}

// NewDispatcher creates a new dispatcher with the given configuration. This
// method is responsible for initializing the socket.
func NewDispatcher(cfg Config) Dispatcher {
  ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
  if err != nil {
    log.Fatalf("Could not create server: %v", err)
  }

  inbox := make(chan Message)

  return Dispatcher{
    listener: ln,

    Config:   cfg,
    Inbox:    inbox,

    nicks: make(map[string]int64),
    users: make(map[string]int64),

    nextID: 1,
  }
}

// Loop is the entry point for the local server. This method does not return.
func (d *Dispatcher) Loop() {
  go d.acceptLoop()
  d.handleLoop()
}

// acceptLoop waits for a new connection, and spins off a new Relay in a
// separate go routine for each new connection.
func (d *Dispatcher) acceptLoop() {
  for {
    conn, err := d.listener.Accept()
    if err != nil {
      log.Printf("Could not accept new connection: ", err)
      continue
    }

    // Kick off a new Relay.
    relay := d.NewRelay(conn)
    go relay.Loop()
  }
}

// dispatchLoop continuously reads from the message channel and handles one
// message at a time.
func (d *Dispatcher) handleLoop() {
  for {
    msg := <-d.Inbox
    log.Printf("Got new message. prefix = '%s', command = '%s', params = %q",
      msg.Prefix, msg.Command, msg.Params)
    msg.Relay.Handler(msg)
  }
}
