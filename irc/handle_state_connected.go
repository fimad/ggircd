package irc

func (d *Dispatcher) getHandleStateConnected(client *Client, server *Server) func(Message) {
  return func(msg Message) {
    switch msg.Command {
    case CmdQuit:
      d.handleCmdQuit(msg, client, server)
    case CmdJoin:
      d.handleCmdJoin(msg, client, server)
    case CmdPing:
      d.handleCmdPing(msg)
    }
  }
}
