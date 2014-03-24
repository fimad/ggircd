package irc

type mockConnection struct {
  messages []Message
  killed   bool
}

func (c *mockConnection) Send(msg Message) {
  c.messages = append(c.messages, msg)
}

func (c *mockConnection) Kill() {
  c.killed = true
}

func (c *mockConnection) Loop() {}
