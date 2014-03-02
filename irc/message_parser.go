package irc

import (
  "bufio"
  "bytes"
  "errors"
  "io"
  "log"
  "regexp"
)

// MessageParser is a function that returns a Message and a boolean indicating
// if the end of the stream has been reached. If the boolean is false, then the
// returned Message should be ignored and the end of the input has been reached.
type MessageParser func() (Message, bool)

// NewMessageParser will create a new Parser function that can be called
// repeatedly to parse Messages from the given io.Reader.
func NewMessageParser(reader io.Reader) MessageParser {
  scanner := bufio.NewScanner(reader)
  scanner.Split(splitFunc)
  return func() (Message, bool) {
    for scanner.Scan() {
      msg, ok := parseMessage(scanner.Text())
      if ok {
        return msg, true
      }
    }
    return Message{}, false
  }
}

// splitFunc is a split function used by a scanner. It splits at CR-LF. A result
// of this is that the returned line will not contain the CR-LF token.
func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
  lines := bytes.SplitN(data, []byte("\x0d\x0a"), 2)

  if len(lines) < 2 && !atEOF {
    return 0, nil, nil
  } else if len(lines) < 2 {
    return 0, nil, errors.New("no end of line token in input")
  } else {
    return len(lines[0]) + 2, lines[0], nil
  }
}

// See RFC 952 for the definition of a host name.
//var hostRegex = `[a-zA-Z](?:[a-zA-Z0-9.-]{,22}[a-zA-Z0-9])?`
//var nickRegex = `[a-zA-Z][a-zA-Z0-9\-\[\]\\` + "`" + `^{}]*`
var prefix = `(?::([^ ]+) )?`
var command = `([a-zA-Z]+|[0-9]{3})`
var params = `( .*)?`
var message = regexp.MustCompile(`^` + prefix + command + params + `$`)

// parseMessage takes a raw line from the IRC protocol (minus the trailing CRLF)
// and returns a parsed Message and a bool indicating success.
//
// TODO(will): This does not currently validate the prefix portion of the
// message.
func parseMessage(line string) (Message, bool) {
  log.Printf("raw line: %s", line)
  var msg Message

  // Parse the message into prefix, command and parameter strings.
  parts := message.FindStringSubmatch(line)
  if parts == nil {
    log.Printf("Does not match regexp")
    return msg, false
  }

  msg.Prefix = parts[1]
  msg.Command = parts[2]
  msg.Params = make([]string, 0)

  // Split the parameters into separate strings.
  i := 0
  for i < len(parts[3])-1 {
    // Each parameter must begin with a space.
    if parts[3][i] != ' ' {
      log.Printf("Does not start with a space")
      return msg, false
    }
    // Skip extra white space.
    for parts[3][i] == ' ' {
      i++
    }

    // If the next character is a ':', then the parameter is the value of the
    // remainder of the string.
    if parts[3][i] == ':' {
      msg.Params = append(msg.Params, parts[3][i+1:])
      break
    }

    // Otherwise the current parameter is the substring from i to either the
    // first space or the end of the string.
    start := i
    for i < len(parts[3]) {
      if parts[3][i] == ' ' {
        break
      }
      i++
    }
    msg.Params = append(msg.Params, parts[3][start:i])
  }

  return msg, true
}
