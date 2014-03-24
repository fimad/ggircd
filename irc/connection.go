package irc

import (
	"io"
	"log"
	"net"
)

// Connection corresponds to some end-point that has connected to the IRC
// server.
type Connection interface {
	Sink

	// Loop reads messages from the connection and passes them to the handler.
	Loop()

	// Kill stops the execution of the go routine running Loop.
	Kill()
}

type connectionImpl struct {
	conn      net.Conn
	handler   Handler
	inbox     chan Message
	killRead  chan struct{}
	killWrite chan struct{}
}

// NewConnection creates a new connection with the given network connection and
// handler.
func NewConnection(conn net.Conn, handler Handler) Connection {
	return &connectionImpl{
		conn:      conn,
		handler:   handler,
		inbox:     make(chan Message),
		killRead:  make(chan struct{}),
		killWrite: make(chan struct{}),
	}
}

func (c *connectionImpl) Send(msg Message) {
	c.inbox <- msg
}

func (c *connectionImpl) Loop() {
	go c.writeLoop()
	c.readLoop()
}

func (c *connectionImpl) Kill() {
	go func() {
		c.killRead <- struct{}{}
		c.killWrite <- struct{}{}
	}()
}

func (c *connectionImpl) readLoop() {
	var msg Message
	parser := NewMessageParser(c.conn)

	alive, hasMore := true, true
	for hasMore && alive {
		select {
		case _ = <-c.killRead:
			alive = false
		default:
			msg, hasMore = parser()
			c.handler = c.handler.Handle(c, msg)
		}
	}

	if alive {
		_ = <-c.killRead
	}
	c.conn.Close()
}

func (c *connectionImpl) writeLoop() {
	alive := true
	for alive {
		select {
		case _ = <-c.killWrite:
			alive = false
		case msg := <-c.inbox:
			line, ok := msg.ToString()
			if !ok {
				break
			}

			_, err := io.WriteString(c.conn, line)
			if err != nil {
				log.Printf("Error encountered sending message to client: %v", err)
				break
			}
		}
	}

	if alive {
		_ = <-c.killWrite
	}
}
