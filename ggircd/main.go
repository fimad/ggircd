package main

import (
  "flag"
  "log"

  "github.com/fimad/ggircd/irc/config"
  "github.com/fimad/ggircd/irc/server"
)

var configFile = flag.String("config", "/etc/ggircd.conf",
  "Path to a file containing the irc daemon's configuration.")

func main() {
  flag.Parse()
  log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

  cfg := config.FromJSONFile(*configFile)
  server := server.NewLocal(cfg)
  server.Loop()
}
