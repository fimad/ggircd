package irc

// sendTopic sends the topic of a channel to a given user.
func sendTopic(state State, user *User, channel *Channel) {
  if channel.Topic != "" {
    sendNumericTrailing(state, user, ReplyTopic, channel.Topic, channel.Name)
  }
  sendNumericTrailing(state, user, ReplyNoTopic, "No topic set", channel.Name)
}
