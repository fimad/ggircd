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

// NewRelay creates a new Relay and registers it with the Dispatcher by
// assigning it a unique ID.
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
  d.relayToClient[relay.ID] = make(map[int64]bool)
  d.relayToServer[relay.ID] = make(map[int64]bool)
  d.nextID++
  return relay
}

// KillRelay removes a relay from a Dispatcher. It does not send any messages or
// remove servers or clients.
func (d *Dispatcher) KillRelay(relay *Relay) {
  if d.relayToClient[relay.ID] != nil {
    delete(d.relayToClient, relay.ID)
  }

  if d.relayToServer[relay.ID] != nil {
    delete(d.relayToServer, relay.ID)
  }
}

// Kill shuts down a Relay. It does not handle unregistering it from the
// Dispatcher.
func (r *Relay) Kill() {
  r.killInbox <- killToken{}
  r.killOutbox <- killToken{}
}

// Loop is the entry point for the local server. This method does not return.
func (r *Relay) Loop() {
  go r.inboxLoop()
  r.outboxLoop()
}

// outboxLoop reads messages from the connected client and continuously pushes
// Messages to the Dispatcher via the Outbox channel.
func (r *Relay) outboxLoop() {
  parser := NewMessageParser(r.conn)

  var msg Message
  hasMore := true
  alive := true
  didQuit := false
  for alive && hasMore {
    select {
    case _ = <-r.killOutbox:
      alive = false
    default:
      msg, hasMore = parser()

      // Don't double send QUIT messages. QUIT messages will cause the
      // Dispatcher to send a message that will result in the Relay's loops
      // being torn down. If multiple quit messages are sent then the dispatcher
      // will hang on the second because there will be nothing at the other end
      // to read from the channel.
      didQuit = didQuit || msg.Command == "QUIT"
      if didQuit {
        break
      }

      if !hasMore {
        r.Outbox <- Message{Command: "QUIT", Relay: r}
        break
      }

      msg.Relay = r
      r.Outbox <- msg
    }
  }

  if alive {
    _ = <-r.killOutbox
  }

  r.conn.Close()
}

// inboxLoop continuously pulls messages from the Inbox channel and sends the
// message to the connected client.
func (r *Relay) inboxLoop() {
  alive := true
  shouldKill := false
  for alive {
    select {
    case _ = <-r.killInbox:
      alive = false
    case msg := <-r.Inbox:
      shouldKill = msg.ShouldKill

      line, ok := msg.ToString()
      if !ok {
        break
      }

      _, err := io.WriteString(r.conn, line)
      if err != nil {
        log.Printf("Error encountered sending message to client: %v", err)
        break
      }
    }

    if alive && shouldKill {
      go r.Kill()
    }
  }
}
