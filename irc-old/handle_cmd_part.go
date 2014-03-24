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
    d.sendNumeric(client, ErrorNeedMoreParams)
    return
  }

  reason := "PARTing"
  if len(msg.Params) > 1 {
    reason = msg.Params[1]
  }

  channels := strings.Split(msg.Params[0], ",")
  for i := 0; i < len(channels); i++ {
    name := channels[i]
    channel := d.ChannelForName(name)

    if channel == nil {
      d.sendNumeric(client, ErrorNoSuchChannel, name)
      continue
    }

    if !channel.Clients[client.ID] {
      d.sendNumeric(client, ErrorNotOnChannel, name)
      continue
    }

    d.RemoveFromChannel(channel, client, reason)
  }
}
