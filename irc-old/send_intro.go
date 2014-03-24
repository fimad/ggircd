package irc

import (
  "fmt"
)

const (
  welcomeMessage  = "Wecome to the %s Internet Relay Network"
  yourHostMessage = "Your host is %s, running version %s"
  myInfoMessage   = "%s %s %s %s"
)

// sendIntro sends all of the welcome messages that clients expect to receive
// after joining the server.
func (d *Dispatcher) sendIntro(relay *Relay, client *Client) {
  relay.Inbox <- ReplyWelcome.WithParams(client.Nick,
    fmt.Sprintf(welcomeMessage, d.Config.Network))
  relay.Inbox <- ReplyYourHost.WithParams(client.Nick,
    fmt.Sprintf(yourHostMessage, d.Config.Name, Version))
  d.sendMOTD(relay, client)
}
