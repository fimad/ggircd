package main

import (
  "flag"
  "log"
  "runtime"

  "github.com/fimad/ggircd/irc"
)

var configFile = flag.String("config", "/etc/ggircd/ggircd.conf",
  "Path to a file containing the irc daemon's configuration.")

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  flag.Parse()
  log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

  cfg := irc.ConfigFromJSONFile(*configFile)
  server := irc.NewDispatcher(cfg)
  server.Loop()
}
