package irc

// sendTopic sends the topic of a channel to a given client.
func (d *Dispatcher) sendTopic(client *Client, channel *Channel) {
  msg := ReplyTopic.WithParams(client.Nick, channel.Name, channel.Topic)
  if channel.Topic == "" {
    msg = ReplyNoTopic.WithParams(channel.Name, "Not topic set")
  }
  client.Relay.Inbox <- msg
}
