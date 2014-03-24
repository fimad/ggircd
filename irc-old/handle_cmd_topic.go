package irc

func (d *Dispatcher) handleCmdTopic(msg Message, client *Client) {
  if client == nil {
    return
  }

  if len(msg.Params) < 1 {
    d.sendNumeric(client, ErrorNeedMoreParams)
    return
  }

  name := msg.Params[0]
  channel := d.ChannelForName(name)
  if channel != nil {
    d.sendNumeric(client, ErrorNoSuchChannel, name)
    return
  }

  if msg.Trailing == "" {
    d.sendTopic(client, channel)
    return
  }

  if !channel.Clients[client.ID] {
    d.sendNumeric(client, ErrorNotOnChannel, name)
    return
  }

  if channel.Mode[ChannelModeTopicOp] && !channel.Ops[client.ID] {
    d.sendNumeric(client, ErrorChanOPrivsNeeded, name)
    return
  }

  msg.Prefix = client.Prefix()
  channel.Topic = msg.Trailing
  d.SendToChannel(channel, msg)
}
