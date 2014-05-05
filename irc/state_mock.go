package irc

type mockState struct {
	stateImpl
}

func newMockState() *mockState {
	return &mockState{
		stateImpl{
			users:    make(map[string]*user),
			channels: make(map[string]*channel),
			config: Config{
				Name:    "name",
				Network: "network",
				Port:    6667,
				MOTD:    "",
			},
		},
	}
}

func (s *mockState) withChannel(name, mode, topic string) *mockState {
	s.channels[name] = &channel{
		name:    name,
		mode:    parseMode(channelModes, mode),
		topic:   topic,
		limit:   0,
		key:     "",
		banNick: "",
		banUser: "",
		banHost: "",
		users:   make(map[*user]bool),
		ops:     make(map[*user]bool),
		voices:  make(map[*user]bool),
		invited: make(map[*user]bool),
	}
	return s
}

func (s *mockState) withChannelKey(name, key string) *mockState {
	s.channels[name].key = key
	return s
}

func (s *mockState) withChannelLimit(name string, limit int) *mockState {
	s.channels[name].limit = limit
	return s
}

func (s *mockState) withUser(nick string, channels ...string) *mockState {
	chanMap := make(map[*channel]bool)
	for _, name := range channels {
		ch := s.channels[name]
		if ch == nil {
			continue
		}
		chanMap[ch] = true
	}

	s.users[nick] = &user{
		nick:     nick,
		user:     nick,
		host:     nick,
		server:   nick,
		realName: nick,
		channels: chanMap,
		mode:     make(mode),
		sink:     &mockConnection{},
	}

	for ch := range chanMap {
		ch.users[s.users[nick]] = true
	}

	return s
}

func (s *mockState) withUserMode(nick, modeLine string) *mockState {
	s.users[nick].mode = parseMode(channelModes, modeLine)
	return s
}

func (s *mockState) withUserAway(nick, awayMessage string) *mockState {
	s.users[nick].mode[userModeAway] = awayMessage != ""
	s.users[nick].awayMessage = awayMessage
	return s
}

func (s *mockState) withOps(channel string, nicks ...string) *mockState {
	ch := s.channels[channel]
	for _, nick := range nicks {
		ch.ops[s.users[nick]] = true
	}
	return s
}

func (s *mockState) withVoices(channel string, nicks ...string) *mockState {
	ch := s.channels[channel]
	for _, nick := range nicks {
		ch.voices[s.users[nick]] = true
	}
	return s
}
