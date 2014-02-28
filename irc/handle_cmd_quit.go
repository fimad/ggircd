package irc

func (d *Dispatcher) handleCmdQuit(msg Message, client *Client, server *Server) {
  if client == nil {
    Todo("non-client QUIT")
  }

  nonServer := len(d.relayToServer[msg.Relay.ID]) == 0
  d.KillClient(client)
  if nonServer {
    d.KillRelay(msg.Relay)
    msg.Relay.Kill()
  }
}
