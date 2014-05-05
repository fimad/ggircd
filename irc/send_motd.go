package irc

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

const (
	motdHeader = "- %s Message of the day - "
	motdFooter = "- End of /MOTD command"
)

var (
	motdOnce sync.Once
	motd     []string
)

// sendMOTD will send the message of the day to a relay.
func sendMOTD(state State, sink Sink) {
	motdOnce.Do(func() { loadMOTD(state) })

	sendNumericTrailing(state, sink, ReplyMOTDStart,
		fmt.Sprintf(motdHeader, state.GetConfig().Name))

	for _, line := range motd {
		sendNumericTrailing(state, sink, ReplyMOTD, "- "+line)
	}

	sendNumericTrailing(state, sink, ReplyEndOfMOTD, motdFooter)
}

func loadMOTD(state State) {
	motdFile := state.GetConfig().MOTD
	if motdFile == "" || motd != nil {
		return
	}

	file, err := os.Open(motdFile)
	if err != nil {
		logf(Error, "Could not open MOTD: %v", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	motd = make([]string, 0)
	for scanner.Scan() {
		motd = append(motd, scanner.Text())
	}

	file.Close()
}
