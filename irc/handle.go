package irc

// handleStateNew handles Relays that have just connected and have not been
// classified as either a server or a client.
func (d *Dispatcher) handleStateNew(msg Message) {
  switch msg.Command {
  case "PASS":
    Todo("Implement the PASS case")
  case "NICK":
    d.handleStateNewCmdNick(msg)
  case "SERVER":
    Todo("Implement the SERVER case")
  }
}

// handleStateNewCmdNick handles the NICK command from Relays that have yet to
// be identified as either a client or a server.
func (d *Dispatcher) handleStateNewCmdNick(msg Message) {
  if len(msg.Params) < 1 {
    msg.Relay.Inbox <- ErrorNoNicknameGiven
    return
  }
  nick := msg.Params[0]

  client := d.NewClient(msg.Relay)
  ok, errMsg := d.SetNick(client, nick)
  if !ok {
    errMsg.ShouldKill = true
    msg.Relay.Inbox <- errMsg
    d.KillClient(client)
  }

  msg.Relay.Handler = d.handleStateUser
}

// handleStateUser handles Relay state where the Relay is registering as a
// client and is expected to send the USER command.
func (d *Dispatcher) handleStateUser(msg Message) {
}
