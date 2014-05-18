package irc

import (
	"io"
	"net"
)

// connection corresponds to some end-point that has connected to the IRC
// server.
type connection interface {
	sink

	// loop reads messages from the connection and passes them to the handler.
	loop()

	// kill stops the execution of the go routine running Loop.
	kill()
}

type connectionImpl struct {
	conn      net.Conn
	handler   handler
	inbox     chan message
	killRead  chan struct{}
	killWrite chan struct{}
}

// newConnection creates a new connection with the given network connection and
// handler.
func newConnection(conn net.Conn, handler handler) connection {
	return &connectionImpl{
		conn:      conn,
		handler:   handler,
		inbox:     make(chan message),
		killRead:  make(chan struct{}),
		killWrite: make(chan struct{}),
	}
}

func (c *connectionImpl) send(msg message) {
	c.inbox <- msg
}

func (c *connectionImpl) loop() {
	go c.writeLoop()
	c.readLoop()
}

func (c *connectionImpl) kill() {
	go func() {
		c.killRead <- struct{}{}
		c.killWrite <- struct{}{}
	}()
}

func (c *connectionImpl) readLoop() {
	var msg message
	parser := newMessageParser(c.conn)

	didQuit := false
	alive, hasMore := true, true
	for hasMore && alive {
		select {
		case _ = <-c.killRead:
			alive = false
		default:
			msg, hasMore = parser()
			logf(debug, "< %+v", msg)
			didQuit = didQuit || msg.command == cmdQuit.command
			c.handler = c.handler.handle(c, msg)

			if !hasMore {
				break
			}
		}
	}

	c.conn.Close()

	// If there was never a QUIT message then this is a premature termination and
	// a quit message should be faked.
	if !didQuit {
		c.handler = c.handler.handle(c, cmdQuit.withTrailing("QUITing"))
	}

	if alive {
		_ = <-c.killRead
	}
	close(c.killRead)
}

func (c *connectionImpl) writeLoop() {
	alive := true
	for alive {
		select {
		case _ = <-c.killWrite:
			alive = false
		case msg := <-c.inbox:
			logf(debug, "> %+v", msg)

			line, ok := msg.toString()
			if !ok {
				break
			}

			_, err := io.WriteString(c.conn, line)
			if err != nil {
				logf(warn, "Error encountered sending message to client: %v", err)
				break
			}
		}
	}

	c.conn.Close()

	if alive {
		_ = <-c.killWrite
	}
	close(c.inbox)
	close(c.killWrite)
}
