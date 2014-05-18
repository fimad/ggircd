package irc

import (
	"fmt"
	"net"
)

// RunServer starts the GGircd IRC server. This method will not return.
func RunServer(cfg Config) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logf(fatal, "Could not create server: %v", err)
	}

	state := make(chan state, 1)
	state <- newState(cfg)
	acceptLoop(cfg, ln, state)
}

func acceptLoop(cfg Config, listener net.Listener, state chan state) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logf(warn, "Could not accept new connection: ", err)
			continue
		}

		c := newConnection(cfg, conn, newFreshHandler(state))
		go c.loop()
	}
}
