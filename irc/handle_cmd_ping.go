package irc

func (d *Dispatcher) handleCmdPing(msg Message) {
  msg.Relay.Inbox <- Message{Command: "PONG", Params: []string{d.Config.Name}}
}
