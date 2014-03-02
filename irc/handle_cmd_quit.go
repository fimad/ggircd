package irc

func (d *Dispatcher) handleCmdQuit(msg Message, client *Client, server *Server) {
  if client != nil {
    d.KillClient(client)
  }

  if server != nil {
    Todo("server QUIT")
  }

  // Only kill the relay if there is exactly one connection.
  oneEntity := 1 >=
    len(d.relayToServer[msg.Relay.ID])+
      len(d.relayToClient[msg.Relay.ID])
  if oneEntity {
    d.sendKillingMessage(msg.Relay, Message{})
    //d.KillRelay(msg.Relay)
    //msg.Relay.Kill()
  }
}
