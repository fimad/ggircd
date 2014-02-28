package irc

type Client struct {
  ID   int64
  Nick string
  User string
  Host string

  Relay *Relay
}

func (d *Dispatcher) NewClient(relay *Relay) *Client {
  client := &Client{
    ID:    d.nextID,
    Relay: relay,
  }
  d.relayToClient[relay.ID] = append(d.relayToClient[relay.ID], client)
  d.nextID++
  return client
}

func (d *Dispatcher) SetNick(client *Client, nick string) (bool, Message) {
  if d.nicks[nick] != 0 {
    return false, ErrorNickCollision
  }
  client.Nick = nick
  d.nicks[nick] = client.ID
  return true, Message{}
}

func (d *Dispatcher) KillClient(client *Client) {
  if client.Nick != "" {
    delete(d.nicks, client.Nick)
  }

  if client.User != "" {
    delete(d.users, client.User)
  }
}
