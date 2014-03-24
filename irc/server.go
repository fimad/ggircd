package irc

import (
  "fmt"
  "log"
  "net"
)

func RunServer(cfg Config) {
  ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
  if err != nil {
    log.Fatalf("Could not create server: %v", err)
  }

  state := make(chan State, 1)
  state <- NewState(cfg)
  acceptLoop(ln, state)
}

func acceptLoop(listener net.Listener, state chan State) {
  for {
    conn, err := listener.Accept()
    if err != nil {
      log.Printf("Could not accept new connection: ", err)
      continue
    }

    c := NewConnection(conn, NewFreshHandler())
    go c.Loop()
  }
}
