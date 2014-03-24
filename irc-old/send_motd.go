package irc

import (
  "bufio"
  "fmt"
  "log"
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
func (d *Dispatcher) sendMOTD(relay *Relay, client *Client) {
  motdOnce.Do(d.loadMOTD)

  relay.Inbox <- ReplyMOTDStart.
    WithParams(client.Nick).
    WithTrailing(fmt.Sprintf(motdHeader, d.Config.Name))

  for _, line := range motd {
    relay.Inbox <- ReplyMOTD.
      WithParams(client.Nick).
      WithTrailing("- " + line)
  }

  relay.Inbox <- ReplyEndOfMOTD.
    WithParams(client.Nick).
    WithTrailing(motdFooter)
}

func (d *Dispatcher) loadMOTD() {
  if d.Config.MOTD == "" {
    return
  }

  file, err := os.Open(d.Config.MOTD)
  if err != nil {
    log.Printf("Could not open MOTD: %v", err)
  }

  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  motd = make([]string, 0)
  for scanner.Scan() {
    motd = append(motd, scanner.Text())
  }

  file.Close()
}
