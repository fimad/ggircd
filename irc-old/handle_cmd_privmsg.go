package irc

func (d *Dispatcher) handleCmdPrivMsg(msg Message, client *Client) {
	if client == nil {
		Todo("handle nil clients")
		return
	}

	if len(msg.Params) < 1 {
		d.sendNumeric(client, ErrorNoRecipient)
		return
	}

	if msg.Trailing == "" {
		d.sendNumeric(client, ErrorNoTextToSend)
		return
	}

	target := msg.Params[0]
	msg.Prefix = client.Prefix()

	if target[0] != '#' && target[0] != '&' {
		nickClient, ok := d.ClientForNick(target)
		if !ok {
			d.sendNumeric(client, ErrorNoSuchNick, target)
			return
		}
		nickClient.Relay.Inbox <- msg
		return
	}

	channel := d.ChannelForName(target)

	if channel == nil || !channel.CanPrivMsg(client) {
		d.sendNumeric(client, ErrorCannotSendToChan, target)
		return
	}

	msg.Prefix = client.Prefix()
	for cid := range channel.Clients {
		// Don't send the message to the client that sent it.
		if cid == client.ID {
			continue
		}
		d.clients[cid].Relay.Inbox <- msg
	}
}
