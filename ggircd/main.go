package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/fimad/ggircd/irc"
)

var configFile = flag.String("config", "/etc/ggircd/ggircd.conf",
	"Path to a file containing the irc daemon's configuration.")

var logLevel = flag.Int("log", 3,
	"The log level, the higher the move verbose. (default: 3).")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime)

	irc.SetLogLevel(*logLevel)

	cfg := irc.ConfigFromJSONFile(*configFile)
	irc.RunServer(cfg)
}
