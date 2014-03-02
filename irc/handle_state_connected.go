package irc

func (d *Dispatcher) getHandleStateConnected(client *Client, server *Server) func(Message) {
  return func(msg Message) {
    switch msg.Command {
    case "QUIT":
      d.handleCmdQuit(msg, client, server)
    case "JOIN":
      d.handleCmdJoin(msg, client, server)
    case "PING":
      d.handleCmdPing(msg)
    }
  }
}
