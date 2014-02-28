package irc

// Client contains metadata about connected clients.
type Client struct {
  ID int64

  Nick     string
  User     string
  Host     string
  Server   string
  RealName string

  // The relay that should be used to send messages to this client.
  Relay *Relay
}

// NewClient creates a new client and registers it with the Dispatcher by
// assigning it a unique id.
func (d *Dispatcher) NewClient(relay *Relay) *Client {
  client := &Client{
    ID:    d.nextID,
    Relay: relay,
  }
  d.relayToClient[relay.ID][client.ID] = true
  d.nextID++
  return client
}

// KillClient removes a client from the Dispatchers registry of clients. This
// method does not send any messages, and does not terminate any Relays.
func (d *Dispatcher) KillClient(client *Client) {
  if client.Nick != "" {
    delete(d.nicks, client.Nick)
  }

  if client.User != "" {
    delete(d.users, client.User)
  }

  if d.relayToClient[client.Relay.ID][client.ID] {
    delete(d.relayToClient[client.Relay.ID], client.ID)
  }
}

// SetNick attempts to set the nick name of a given client. Returns a boolean
// indicating success and an error message in the case of failure.
func (d *Dispatcher) SetNick(client *Client, nick string) (bool, Message) {
  if d.nicks[nick] != 0 {
    return false, ErrorNickCollision
  }

  client.Nick = nick
  d.nicks[nick] = client.ID
  return true, Message{}
}

// SetUser attempts to set the user info of a given client. Returns a boolean
// indicating success and an error message in the case of failure.
func (d *Dispatcher) SetUser(client *Client, user, host, server, realName string) (bool, Message) {
  if d.users[user] != 0 {
    return false, ErrorAlreadyRegistred
  }

  client.User = user
  client.Host = host
  client.Server = server
  client.RealName = realName
  d.users[user] = client.ID
  return true, Message{}
}
