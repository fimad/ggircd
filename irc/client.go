package irc

import (
  "fmt"
)

// Client contains metadata about connected clients.
type Client struct {
  ID int64

  Nick     string
  User     string
  Host     string
  Server   string
  RealName string

  Channels map[string]bool

  // The relay that should be used to send messages to this client.
  Relay *Relay
}

// NewClient creates a new client and registers it with the Dispatcher by
// assigning it a unique id.
func (d *Dispatcher) NewClient(relay *Relay) *Client {
  client := &Client{
    ID:       d.nextID,
    Channels: make(map[string]bool),
    Relay:    relay,
  }
  d.clients[client.ID] = client
  d.relayToClient[relay.ID][client.ID] = true
  d.nextID++
  return client
}

// KillClient removes a client from the Dispatchers registry of clients. This
// method does not send any messages, and does not terminate any Relays.
func (d *Dispatcher) KillClient(client *Client) {
  for ch := range client.Channels {
    if channel := d.ChannelForName(ch); channel != nil {
      d.RemoveFromChannel(channel, client, "QUIT")
    }
  }

  delete(d.clients, client.ID)
  delete(d.nicks, client.Nick)
  delete(d.relayToClient[client.Relay.ID], client.ID)
}

// SetNick attempts to set the nick name of a given client. Returns a boolean
// indicating success.
func (d *Dispatcher) SetNick(client *Client, nick string) bool {
  nick = Lowercase(nick)

  if d.nicks[nick] != 0 {
    d.sendNumeric(client, ErrorNicknameInUse, nick)
    return false
  }

  msg := Message{
    Command: CmdNick,
    Prefix:  client.Prefix(),
    Params:  []string{nick},
  }
  for name := range client.Channels {
    d.SendToChannel(d.ChannelForName(name), msg)
  }

  oldNick := client.Nick
  client.Nick = nick
  d.nicks[nick] = client.ID
  delete(d.nicks, oldNick)

  return true
}

// SetUser attempts to set the user info of a given client. Returns a boolean
// indicating success.
func (d *Dispatcher) SetUser(client *Client, user, host, server, realName string) bool {
  client.User = user
  client.Host = host
  client.Server = server
  client.RealName = realName
  return true
}

// Prefix returns the prefix string that should be used in Messages originating
// from this Client.
func (c *Client) Prefix() string {
  return fmt.Sprintf("%s!%s@%s", c.Nick, c.User, c.Host)
}
