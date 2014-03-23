package irc

func (d *Dispatcher) handleCmdNick(msg Message, client *Client) {
  if client == nil {
    Todo("handle nil clients")
    return
  }

  if len(msg.Params) != 1 {
    client.Relay.Inbox <- ErrorNoNicknameGiven.
      WithTrailing("No Nickname given")
    return
  }

  nick := msg.Params[0]
  if _, found := d.nicks[nick]; found {
    client.Relay.Inbox <- ErrorNicknameInUse.
      WithParams(client.Nick, nick).
      WithTrailing("Nick name in use")
    return
  }

  msg.Prefix = client.Prefix()
  for name := range client.Channels {
    d.SendToChannel(d.ChannelForName(name), msg)
  }

  oldNick := client.Nick
  client.Nick = nick
  d.nicks[nick] = client.ID
  delete(d.nicks, oldNick)
}
