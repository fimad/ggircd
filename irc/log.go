package irc

import (
	"log"
)

const (
	fatal = 0
	warn  = 2
	info  = 3
	debug = 4
)

var logLevel = info

var levelToSymbol = map[int]string{
	fatal: "F",
	warn:  "W",
	info:  "I",
	debug: "D",
}

// SetLogLevel sets the verbosity of logging.
func SetLogLevel(level int) {
	logLevel = level
}

func logf(level int, format string, args ...interface{}) {
	if level <= logLevel {
		symbol := levelToSymbol[level]
		if symbol == "" {
			symbol = "?"
		}
		format = symbol + ": " + format
		if level > 0 {
			log.Printf(format, args...)
		} else {
			log.Panicf(format, args...)
		}
	}
}
