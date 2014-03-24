package irc

type Channel struct {
  Name string

  Mode Mode

  Topic string

  Limit int
  Key   string

  BanNick string
  BanUser string
  BanHost string

  Users  map[*User]bool
  Ops    map[*User]bool
  Voices map[*User]bool
}

func (ch Channel) Send(msg Message) {
  for user := range ch.Users {
    user.Sink.Send(msg)
  }
}

// ForUsers iterates over all of the users in the channel and passes a pointer
// to each to the supplied callback.
func (ch Channel) ForUsers(callback func(*User)) {
  for u := range ch.Users {
    callback(u)
  }
}

// IsBanned takes a user and returns a boolean indicating if that user is banned
// in this channel.
func (ch Channel) IsBanned(user *User) bool {
  // TODO(will): Actually implement banned users.
  return false
}
