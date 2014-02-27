package irc

import (
  "io"
  "log"
  "net"
)

type remoteClient struct {
  conn   net.Conn
  inbox  chan Message
  outbox chan<- Message
}

func (c *remoteClient) GetInbox() chan<- Message {
  return c.inbox
}

func (c *remoteClient) Loop() {
  go c.inboxLoop()
  c.outboxLoop()
}

// outboxLoop reads messages from the connected client and continuously pushes
// Messages to the LocalServer via the send channel.
func (c *remoteClient) outboxLoop() {
  parser := NewMessageParser(c.conn)
  for {
    msg, ok := parser()
    if !ok {
      continue
    }
    c.outbox <- msg
  }
}

// inboxLoop continuously pulls messages from the recv channel and sends the
// message to the connected client.
func (c *remoteClient) inboxLoop() {
  for {
    msg := <-c.inbox
    line, ok := msg.ToString()
    if !ok {
      log.Printf("Malformed message: %v", msg)
      continue
    }

    _, err := io.WriteString(c.conn, line)
    if err != nil {
      log.Printf("Error encountered sending message to client: %v", err)
      continue
    }
  }
}
