package irc

func sendChannelMode(state State, user *User, channel *Channel) {
  var mode string
  for flag, set := range channel.Mode {
    if set {
      mode += string(flag)
    }
  }
  sendNumeric(state, user, ReplyChannelModeIs, channel.Name, mode)
}
