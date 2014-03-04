package irc

// sendNames sends the messages associated with a NAMES request.
func (d *Dispatcher) sendNames(client *Client, channels ...*Channel) {
  for _, channel := range channels {
    params := make([]string, 3, 4)

    params[0] = client.Nick

    if channel.Mode[ChannelModeSecret] {
      params[1] = "@"
    } else if channel.Mode[ChannelModePrivate] {
      params[1] = "*"
    } else {
      params[1] = "="
    }

    params[2] = channel.Name

    for cid := range channel.Clients {
      nick := d.clients[cid].Nick
      if channel.Ops[cid] {
        nick = "@" + nick
      } else if channel.Voice[cid] {
        nick = "+" + nick
      }
      params = append(params, nick)
    }

    client.Relay.Inbox <- ReplyNamReply.
      WithPrefix(d.Config.Name).
      WithParams(params...)
    client.Relay.Inbox <- ReplyEndOfNames.
      WithPrefix(d.Config.Name).
      WithParams(client.Nick, channel.Name, "End of NAMES list")
  }
}
