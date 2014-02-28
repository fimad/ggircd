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

  msg.Relay.Handler = d.getHandleStateUser(client)
}

// handleStateUser handles Relay state where the Relay is registering as a
// client and is expected to send the USER command.
func (d *Dispatcher) getHandleStateUser(client *Client) func(Message) {
  return func(msg Message) {
    if len(msg.Params) < 4 {
      msg.Relay.Inbox <- ErrorNeedMoreParams
      return
    }
    user := msg.Params[0]
    host := msg.Params[1]
    server := msg.Params[2]
    realName := msg.Params[3]

    client := d.NewClient(msg.Relay)
    ok, errMsg := d.SetUser(client, user, host, server, realName)
    if !ok {
      errMsg.ShouldKill = true
      msg.Relay.Inbox <- errMsg
      d.KillClient(client)
    }

    msg.Relay.Handler = d.getHandleStateConnected(client, nil)
  }
}
