package irc

import (
  "strings"
)

func (d *Dispatcher) handleCmdNames(msg Message, client *Client) {
  if client == nil {
    return
  }

  if len(msg.Params) == 0 {
    d.handleCmdNamesAll(msg, client)
    return
  }
  d.handleCmdNamesChannel(msg, client)
}

func (d *Dispatcher) handleCmdNamesAll(msg Message, client *Client) {
  for _, channel := range d.channels {
    if channel.Mode[ChannelModePrivate] && !channel.Clients[client.ID] {
      break
    }

    if channel.Mode[ChannelModeSecret] && !channel.Clients[client.ID] {
      break
    }

    d.sendNames(client, channel)
  }
}

func (d *Dispatcher) handleCmdNamesChannel(msg Message, client *Client) {
  names := strings.Split(msg.Params[0], ",")
  for _, name := range names {
    channel := d.channels[name]
    if channel == nil {
      break
    }

    if channel.Mode[ChannelModePrivate] && !channel.Clients[client.ID] {
      break
    }

    d.sendNames(client, channel)
  }
}
