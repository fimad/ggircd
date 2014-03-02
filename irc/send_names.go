package irc

// sendNames sends the messages associated with a NAMES request.
func (d *Dispatcher) sendNames(relay *Relay, channels ...*Channel) {
  for _, channel := range channels {
    params := make([]string, 2, 3)
    if channel.Mode[ChannelModeSecret] {
      params[0] = "@"
    } else if channel.Mode[ChannelModePrivate] {
      params[0] = "*"
    } else {
      params[0] = "="
    }

    params[1] = channel.Name

    for cid := range channel.Clients {
      params = append(params, d.clients[cid].Nick)
    }

    relay.Inbox <- ReplyNamReply.WithParams(params...)
  }
  relay.Inbox <- ReplyEndOfNames.WithParams("End of NAMES list")
}
