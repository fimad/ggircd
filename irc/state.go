package irc

// State represents the state of this IRC server. State is not safe for
// concurrent access.
type State interface {
  GetConfig() Config

  ForChannels(func(*Channel))
  ForUsers(func(*User))

  GetChannel(string) *Channel
  GetUser(string) *User

  NewChannel(string) *Channel
  NewUser(string) *User

  SetNick(*User, string) bool
  JoinChannel(*Channel, *User)
}

// stateImpl is a concrete implementation of the State interface.
type stateImpl struct {
  config   Config
  channels map[string]*Channel
  users    map[string]*User
}

func NewState(config Config) State {
  return &stateImpl{
    config:   config,
    channels: make(map[string]*Channel),
    users:    make(map[string]*User),
  }
}

// GetConfig returns the IRC server's configuration.
func (s stateImpl) GetConfig() Config {
  return s.config
}

// ForChannels iterates over all of the registered channels and passes a pointer
// to each to the supplied callback.
func (s stateImpl) ForChannels(callback func(*Channel)) {
  for _, ch := range s.channels {
    callback(ch)
  }
}

// ForUsers iterates over all of the registered users and passes a pointer to
// each to the supplied callback.
func (s stateImpl) ForUsers(callback func(*User)) {
  for _, u := range s.users {
    callback(u)
  }
}

// GetChannel returns a pointer to the channel struct corresponding to the given
// channel name.
func (s stateImpl) GetChannel(name string) *Channel {
  return s.channels[Lowercase(name)]
}

// GetUser returns a pointer to the user struct corresponding to the given
// nickname.
func (s stateImpl) GetUser(nick string) *User {
  return s.users[Lowercase(nick)]
}

// NewChannel creates a new channel with the given name and the appropriate
// default values.
func (s *stateImpl) NewChannel(name string) *Channel {
  name = Lowercase(name)
  if s.channels[name] != nil {
    return nil
  }

  ch := &Channel{Name: name}
  s.channels[name] = ch
  return ch
}

// NewUser creates a new user with the given nickname and the appropriate
// default values.
func (s *stateImpl) NewUser(nick string) *User {
  nick = Lowercase(nick)
  if s.users[nick] != nil {
    return nil
  }

  u := &User{Nick: nick}
  s.users[nick] = u
  return u
}

// SetNick updates the nickname for the given user. If there is a user with the
// given nickname then this method does nothing and returns false.
func (s *stateImpl) SetNick(user *User, nick string) bool {
  nick = Lowercase(nick)
  if s.users[nick] != nil {
    return false
  }

  delete(s.users, user.Nick)
  s.users[nick] = user
  user.Nick = nick
  return true
}

// JoinChannel adds a user to a channel. It does not perform any permissions
// checking, it only updates pointers.
func (s *stateImpl) JoinChannel(channel *Channel, user *User) {
  channel.Clients[user] = true
  user.Channels[channel] = true
}
