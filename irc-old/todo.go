package irc

import (
  "log"
  "runtime"
)

func Todo(msg string) {
  _, file, line, _ := runtime.Caller(1)
  log.Printf("TODO(%s:%d): %s", file, line, msg)
}
