package irc

import (
  "io"
  "log"
  "net"
)

type State int

type killToken struct{}

type Relay struct {
  conn net.Conn

  ID      int64
  Inbox   chan Message
  Outbox  chan<- Message
  Handler func(Message)

  killInbox  chan killToken
  killOutbox chan killToken
}

func (d *Dispatcher) NewRelay(conn net.Conn) *Relay {
  relay := &Relay{
    ID:         d.nextID,
    conn:       conn,
    Inbox:      make(chan Message),
    Outbox:     d.Inbox,
    Handler:    d.handleStateNew,
    killInbox:  make(chan killToken),
    killOutbox: make(chan killToken),
  }
  d.nextID++
  return relay
}

// Kill shuts down a Relay. It does not handle unregistering it from the
// Dispatcher.
func (r *Relay) Kill() {
  r.killInbox <- killToken{}
  r.killOutbox <- killToken{}
  r.conn.Close()
}

// Loop is the entry point for the local server. This method does not return.
func (r *Relay) Loop() {
  go r.inboxLoop()
  r.outboxLoop()
}

// outboxLoop reads messages from the connected client and continuously pushes
// Messages to the LocalServer via the send channel.
func (r *Relay) outboxLoop() {
  parser := NewMessageParser(r.conn)
  for {
    select {
    case _ = <-r.killOutbox:
      break
    default:
      msg, ok := parser()
      if !ok {
        continue
      }
      msg.Relay = r
      r.Outbox <- msg
    }
  }
}

// inboxLoop continuously pulls messages from the recv channel and sends the
// message to the connected client.
func (r *Relay) inboxLoop() {
  for {
    select {
    case _ = <-r.killInbox:
      break
    case msg := <-r.Inbox:
      line, ok := msg.ToString()
      if !ok {
        log.Printf("Malformed message: %v", msg)
        continue
      }

      _, err := io.WriteString(r.conn, line)
      if err != nil {
        log.Printf("Error encountered sending message to client: %v", err)
        continue
      }
    }
  }
}
