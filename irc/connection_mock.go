package irc

type MockConnection struct {
  Messages []Message
  Killed   bool
}

func (c *MockConnection) Send(msg Message) {
  c.Messages = append(c.Messages, msg)
}

func (c *MockConnection) Kill() {
  c.Killed = true
}

func (c *MockConnection) Loop() {}
