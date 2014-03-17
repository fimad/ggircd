package irc

func (d *Dispatcher) handleCmdTopic(msg Message, client *Client) {
  if client == nil {
    return
  }

  if len(msg.Params) < 1 {
    msg.Relay.Inbox <- ErrorNeedMoreParams.WithParams(msg.Command)
    return
  }

  name := msg.Params[0]
  channel := d.channels[name]
  if channel == nil {
    msg.Relay.Inbox <- ErrorNoSuchChannel.WithParams(name)
    return
  }

  if msg.Trailing == "" {
    d.sendTopic(client, channel)
    return
  }

  if !channel.Clients[client.ID] {
    msg.Relay.Inbox <- ErrorNotOnChannel.WithParams(name)
    return
  }

  if channel.Mode[ChannelModeTopicOp] && !channel.Ops[client.ID] {
    msg.Relay.Inbox <- ErrorChanOPrivsNeeded.WithParams(name)
    return
  }

  msg.Prefix = client.Prefix()
  channel.Topic = msg.Trailing
  d.SendToChannel(channel, msg)
}
