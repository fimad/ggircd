package irc

// handleStateNew handles Relays that have just connected and have not been
// classified as either a server or a client.
func (d *Dispatcher) handleStateNew(msg Message) {
  switch msg.Command {
  case CmdPass:
    Todo("Implement the PASS case")
  case CmdNick:
    d.handleStateNewCmdNick(msg)
  case CmdServer:
    Todo("Implement the SERVER case")
  case CmdQuit:
    d.handleCmdQuit(msg, nil, nil)
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
  if !d.SetNick(client, nick) {
    return
  }

  msg.Relay.Handler = d.getHandleStateUser(client)
}

// getHandleStateUser returns a handler for a Relay state where the Relay is
// registering as a client and is expected to send the USER command.
func (d *Dispatcher) getHandleStateUser(client *Client) func(Message) {
  return func(msg Message) {
    switch msg.Command {
    case CmdUser:
      d.handleStateUserCmdUser(msg, client)
    case CmdQuit:
      d.handleCmdQuit(msg, client, nil)
    }
  }
}

// getHandleStateUser returns a handler for a Relay state where the Relay is
// registering as a client and is expected to send the USER command.
func (d *Dispatcher) handleStateUserCmdUser(msg Message, client *Client) {
  if len(msg.Params) < 3 || msg.Trailing == "" {
    msg.Relay.Inbox <- ErrorNeedMoreParams
    return
  }

  user := msg.Params[0]
  host := msg.Params[1]
  server := msg.Params[2]
  realName := msg.Trailing
  if !d.SetUser(client, user, host, server, realName) {
    return
  }

  d.sendIntro(msg.Relay, client)

  msg.Relay.Handler = d.getHandleStateConnected(client, nil)
}
