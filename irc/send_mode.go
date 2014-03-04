package irc

func (d *Dispatcher) sendChannelMode(client *Client, channel *Channel) {
  var mode string
  for flag := range channel.Mode {
    mode += string(flag)
  }
  client.Relay.Inbox <- ReplyChannelModeIs.
    WithParams(client.Nick, channel.Name, mode)
}
