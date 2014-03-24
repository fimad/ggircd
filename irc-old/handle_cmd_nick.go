package irc

func (d *Dispatcher) handleCmdNick(msg Message, client *Client) {
  if client == nil {
    Todo("handle nil clients")
    return
  }

  if len(msg.Params) != 1 {
    d.sendNumeric(client, ErrorNoNicknameGiven)
    return
  }

  d.SetNick(client, msg.Params[0])
}
