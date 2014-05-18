package irc

type mockConnection struct {
	messages []message
	killed   bool
}

func (c *mockConnection) send(msg message) {
	c.messages = append(c.messages, msg)
}

func (c *mockConnection) close() {
}

func (c *mockConnection) kill() {
	c.killed = true
}

func (c *mockConnection) loop() {}
