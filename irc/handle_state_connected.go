package irc

func (d *Dispatcher) getHandleStateConnected(client *Client, server *Server) func(Message) {
  return func(msg Message) {
    switch msg.Command {
    case CmdJoin:
      d.handleCmdJoin(msg, client, server)
    case CmdPart:
      d.handleCmdPart(msg, client, server)
    case CmdPing:
      d.handleCmdPing(msg)
    case CmdQuit:
      d.handleCmdQuit(msg, client, server)
    }
  }
}
