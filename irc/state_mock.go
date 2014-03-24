package irc

type mockState struct {
  stateImpl
}

func newMockState() *mockState {
  return &mockState{
    stateImpl{
      users:    make(map[string]*User),
      channels: make(map[string]*Channel),
      config: Config{
        Name:    "name",
        Network: "network",
        Port:    6667,
        MOTD:    "",
      },
    },
  }
}

func (s *mockState) withChannel(name string) *mockState {
  s.channels[name] = &Channel{
    Name:    name,
    Topic:   "",
    Limit:   0,
    Key:     "",
    BanNick: "",
    BanUser: "",
    BanHost: "",
    Users:   make(map[*User]bool),
    Ops:     make(map[*User]bool),
    Voices:  make(map[*User]bool),
  }
  return s
}

func (s *mockState) withUser(nick string, channels ...string) *mockState {
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
  }
  return s
}
