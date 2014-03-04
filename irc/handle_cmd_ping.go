package irc

func (d *Dispatcher) handleCmdPing(msg Message) {
  msg.Relay.Inbox <- Message{
    Prefix:     d.Config.Name,
    Command:    CmdPong,
    Params:     []string{d.Config.Name, d.Config.Name},
    ForceColon: true,
  }
}
