package irc

// sendNumeric sends a numeric response to the given client.
func (d *Dispatcher) sendNumeric(client *Client, msg Message, extra ...string) {
  params := make([]string, 0, len(extra)+1)
  params = append(params, client.Nick)
  params = append(params, extra...)
  client.Relay.Inbox <- msg.WithPrefix(d.Config.Name).WithParams(params...)
}

// sendNumericTrailing sends a numeric response to the given client.
func (d *Dispatcher) sendNumericTrailing(client *Client, msg Message, trailing string, extra ...string) {
  d.sendNumeric(client, msg.WithTrailing(trailing), extra...)
}
