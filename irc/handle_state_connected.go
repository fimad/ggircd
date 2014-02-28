package irc

func (d *Dispatcher) getHandleStateConnected(client *Client, server *Server) func(Message) {
  return func(msg Message) {
    switch msg.Command {
    }
  }
}
