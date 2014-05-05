package irc

func sendChannelMode(state state, user *user, channel *channel) {
	var mode string
	for flag, set := range channel.mode {
		if set {
			mode += string(flag)
		}
	}
	sendNumeric(state, user, replyChannelModeIs, channel.name, mode)
}
