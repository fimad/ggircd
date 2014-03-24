package irc

type MockState stateImpl

func NewMockState() MockState {
  return MockState(stateImpl{
    users:    make(map[string]*User),
    channels: make(map[string]*Channel),
    config: Config{
      Name:    "name",
      Network: "network",
      Port:    6667,
      MOTD:    "",
    },
  })
}

func (s MockState) WithChannel(name string) MockState {
  s.channels[name] = &Channel{
    Name:    name,
    Topic:   "",
    Limit:   0,
    Key:     "",
    BanNick: "",
    BanUser: "",
    BanHost: "",
    Clients: make(map[*User]bool),
    Ops:     make(map[*User]bool),
    Voice:   make(map[*User]bool),
    Sink:    &SliceSink{},
  }
  return s
}

func (s MockState) WithUser(nick string, channels ...string) MockState {
  chanMap := make(map[*Channel]bool)
  for _, name := range channels {
    ch := s.channels[name]
    if ch == nil {
      continue
    }
    chanMap[ch] = true
  }

  s.users[nick] = &User{
    Nick:     nick,
    User:     nick,
    Host:     nick,
    Server:   nick,
    RealName: nick,
    Channels: chanMap,
    Sink:     &SliceSink{},
  }
  return s
}
