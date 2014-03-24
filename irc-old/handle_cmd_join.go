package irc

import (
  "strings"
)

func (d *Dispatcher) handleCmdJoin(msg Message, client *Client, server *Server) {
  if client == nil {
    Todo("handle nil clients")
    return
  }

  if len(msg.Params) == 0 {
    d.sendNumeric(client, ErrorNeedMoreParams)
    return
  }
  channels := strings.Split(msg.Params[0], ",")

  var keys []string
  if len(msg.Params) > 1 {
    keys = strings.Split(msg.Params[1], ",")
    return
  }

  for i := 0; i < len(channels); i++ {
    name := channels[i]
    channel := d.GetChannel(name)

    if channel == nil {
      d.sendNumeric(client, ErrorNoSuchChannel, name)
      continue
    }

    if channel.Mode[ChannelModeInvite] {
      d.sendNumeric(client, ErrorInviteOnlyChan, name)
      continue
    }

    if channel.Mode[ChannelModeKey] && keys[i] != channel.Key {
      d.sendNumeric(client, ErrorBadChannelKey, name)
      continue
    }

    if channel.Mode[ChannelModeUserLimit] &&
      len(channel.Clients) >= channel.Limit {
      d.sendNumeric(client, ErrorChannelIsFull, name)
      continue
    }

    if channel.IsBanned(client) {
      d.sendNumeric(client, ErrorBannedFromChan, name)
      continue
    }

    d.AddToChannel(channel, client)
  }
}
