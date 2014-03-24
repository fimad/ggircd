package irc

// sendTopic sends the topic of a channel to a given client.
func (d *Dispatcher) sendTopic(client *Client, channel *Channel) {
  msg := ReplyTopic.
    WithPrefix(d.Config.Name).
    WithParams(client.Nick, channel.Name).
    WithTrailing(channel.Topic)

  if channel.Topic == "" {
    msg = ReplyNoTopic.
      WithPrefix(d.Config.Name).
      WithParams(channel.Name).
      WithTrailing("Not topic set")
  }

  client.Relay.Inbox <- msg
}
