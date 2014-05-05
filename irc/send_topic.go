package irc

// sendTopic sends the topic of a channel to a given user.
func sendTopic(state state, user *user, channel *channel) {
	if channel.topic != "" {
		sendNumericTrailing(state, user, replyTopic, channel.topic, channel.name)
	} else {
		sendNumericTrailing(state, user, replyNoTopic, "No topic set", channel.name)
	}
}
