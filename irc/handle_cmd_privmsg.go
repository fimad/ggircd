package irc

func (d *Dispatcher) handleCmdPrivMsg(msg Message, client *Client) {
  if client == nil {
    Todo("handle nil clients")
    return
  }

  if len(msg.Params) < 1 {
    client.Relay.Inbox <- ErrorNoRecipient
    return
  }

  if msg.Trailing == "" {
    client.Relay.Inbox <- ErrorNoTextToSend
    return
  }

  target := msg.Params[0]
  msg.Prefix = client.Prefix()

  if target[0] != '#' && target[0] != '&' {
    cid, ok := d.nicks[target]
    if !ok {
      client.Relay.Inbox <- ErrorNoSuchNick.
        WithParams(target).
        WithTrailing("No such nick")
      return
    }
    d.clients[cid].Relay.Inbox <- msg
    return
  }

  channel := d.channels[target]

  if channel == nil || !channel.CanPrivMsg(client) {
    client.Relay.Inbox <- ErrorCannotSendToChan.
      WithParams(target).
      WithTrailing("Not allowed")
    return
  }

  msg.Prefix = client.Prefix()
  for cid := range channel.Clients {
    // Don't send the message to the client that sent it.
    if cid == client.ID {
      continue
    }
    d.clients[cid].Relay.Inbox <- msg
  }
}
