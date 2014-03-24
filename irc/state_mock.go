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

func (s *mockState) withChannel(name string, mode string) *mockState {
	s.channels[name] = &Channel{
		Name:    name,
		Mode:    ParseMode(ChannelModes, mode),
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

func (s *mockState) withChannelKey(name string, key string) *mockState {
	s.channels[name].Key = key
	return s
}

func (s *mockState) withChannelLimit(name string, limit int) *mockState {
	s.channels[name].Limit = limit
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
		Sink:     &mockConnection{},
	}

	for ch := range chanMap {
		ch.Users[s.users[nick]] = true
	}

	return s
}
