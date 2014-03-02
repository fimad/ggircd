package irc

import (
  "log"
  "strings"
)

func (d *Dispatcher) handleCmdJoin(msg Message, client *Client, server *Server) {
  if client == nil {
    Todo("handle nil clients")
    return
  }

  if len(msg.Params) == 0 {
    msg.Relay.Inbox <- ErrorNeedMoreParams
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

    log.Printf("trying to join channel: %s", name)

    if channel == nil {
      msg.Relay.Inbox <- ErrorNoSuchChannel.WithParams(name, "No such channel")
      continue
    }

    if channel.Mode[ChannelModeInvite] {
      msg.Relay.Inbox <- ErrorInviteOnlyChan.WithParams(name, "Invite only")
      continue
    }

    if channel.Mode[ChannelModeKey] && keys[i] != channel.Key {
      msg.Relay.Inbox <- ErrorBadChannelKEY.WithParams(name, "Incorrect key")
      continue
    }

    if channel.Mode[ChannelModeUserLimit] &&
      len(channel.Clients) >= channel.Limit {
      msg.Relay.Inbox <- ErrorChannelIsFull.WithParams(name, "Channel full")
      continue
    }

    if channel.IsBanned(client) {
      msg.Relay.Inbox <- ErrorBannedFromChan.WithParams(name, "Banned")
      continue
    }

    d.AddToChannel(channel, client)
  }
}
