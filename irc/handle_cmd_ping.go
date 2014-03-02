package irc

func (d *Dispatcher) handleCmdPing(msg Message) {
  msg.Relay.Inbox <- Message{Command: CmdPong, Params: []string{d.Config.Name}}
}
