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
func sendIntro(state State, sink Sink) {
  sendNumericTrailing(state, sink, ReplyWelcome,
    fmt.Sprintf(welcomeMessage, state.GetConfig().Network))
  sendNumericTrailing(state, sink, ReplyYourHost,
    fmt.Sprintf(yourHostMessage, state.GetConfig().Name, Version))
  // d.sendMOTD(relay, client)
}
