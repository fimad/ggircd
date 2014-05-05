package irc

import "log"

const (
	Fatal = 0
	Error = 1
	Warn  = 2
	Info  = 3
	Debug = 4
)

var logLevel = Info

var levelToSymbol = map[int]string{
	Fatal: "F",
	Error: "E",
	Warn:  "W",
	Info:  "I",
	Debug: "D",
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
