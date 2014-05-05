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
func sendIntro(state state, sink sink) {
	sendNumericTrailing(state, sink, replyWelcome,
		fmt.Sprintf(welcomeMessage, state.getConfig().Network))

	sendNumericTrailing(state, sink, replyYourHost,
		fmt.Sprintf(yourHostMessage, state.getConfig().Name, version))

	sendMOTD(state, sink)
}
