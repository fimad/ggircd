package irc

import (
  "fmt"
  "io"
)

const (
  ErrorNeedMoreParams = 461
)

func SendError(conn io.Writer, err int) {
  io.WriteString(conn, fmt.Sprintf("%03d\x0d\x0a", err))
}
