package irc

func (d *Dispatcher) getHandleStateConnected(client *Client, server *Server) func(Message) {
  return func(msg Message) {
    switch msg.Command {
    case CmdJoin:
      d.handleCmdJoin(msg, client, server)
    case CmdMode:
      d.handleCmdMode(msg, client)
    case CmdNames:
      d.handleCmdNames(msg, client)
    case CmdPart:
      d.handleCmdPart(msg, client, server)
    case CmdPing:
      d.handleCmdPing(msg)
    case CmdPrivMsg:
      d.handleCmdPrivMsg(msg, client)
    case CmdQuit:
      d.handleCmdQuit(msg, client, server)
    case CmdTopic:
      d.handleCmdTopic(msg, client)
    }
  }
}
