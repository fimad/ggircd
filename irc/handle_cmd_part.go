package irc

import (
  "strings"
)

func (d *Dispatcher) handleCmdPart(msg Message, client *Client, server *Server) {
  if client == nil {
    Todo("handle nil clients")
    return
  }

  if len(msg.Params) == 0 {
    msg.Relay.Inbox <- ErrorNeedMoreParams
    return
  }

  reason := "PARTing"
  if len(msg.Params) > 1 {
    reason = msg.Params[1]
  }

  channels := strings.Split(msg.Params[0], ",")
  for i := 0; i < len(channels); i++ {
    name := channels[i]
    channel := d.GetChannel(name)

    if channel == nil {
      client.Relay.Inbox <- ErrorNoSuchChannel.WithParams(
        channel.Name, "No such channel")
      continue
    }

    if !channel.Clients[client.ID] {
      client.Relay.Inbox <- ErrorNotOnChannel.WithParams(
        channel.Name, "Not on channel")
      continue
    }

    d.RemoveFromChannel(channel, client, reason)
  }
}
