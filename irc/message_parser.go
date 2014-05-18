package irc

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

// messageParser is a function that returns a Message and a boolean indicating
// if the end of the stream has been reached. If the boolean is false, then the
// returned Message should be ignored and the end of the input has been reached.
type messageParser func() (message, bool)

// newMessageParser will create a new Parser function that can be called
// repeatedly to parse Messages from the given io.Reader.
func newMessageParser(reader io.Reader) messageParser {
	scanner := bufio.NewScanner(reader)
	scanner.Split(splitFunc)
	return func() (message, bool) {
		for scanner.Scan() {
			msg, ok := parseMessage(scanner.Text())
			if ok {
				return msg, true
			}
		}
		return message{}, false
	}
}

// splitFunc is a split function used by a scanner. It splits at CR-LF, LF-CR,
// CR and LF. A result of this is that the returned line will not contain the
// CR-LF token.
func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := range data {
		// Handle the case where a '\r' or '\n' is encountered at the very end of
		// the buffer.
		if (data[i] == '\x0d' || data[i] == '\x0a') && i == len(data)-1 {
			return i + 1, data[0:i], nil
		}

		// Handle the case where either a CR-LF or LF-CR is encountered.
		if (data[i] == '\x0d' && data[i+1] == '\x0a') ||
			(data[i] == '\x0a' && data[i+1] == '\x0d') {
			return i + 2, data[0:i], nil
		}

		// Handle the case where either a lone CR or LF is encountered.
		if data[i] == '\x0d' || data[i] == '\x0a' {
			return i + 1, data[0:i], nil
		}
	}

	if atEOF {
		return 0, nil, errors.New("no end of line token in input")
	}
	return 0, nil, nil
}

// See RFC 952 for the definition of a host name.
//var hostRegex = `[a-zA-Z](?:[a-zA-Z0-9.-]{,22}[a-zA-Z0-9])?`
//var nickRegex = `[a-zA-Z][a-zA-Z0-9\-\[\]\\` + "`" + `^{}]*`
var prefix = `(?::([^ ]+) )?`
var command = `([a-zA-Z]+|[0-9]{3})`
var params = `( .*)?`
var messageRegex = regexp.MustCompile(`^` + prefix + command + params + `$`)

// parseMessage takes a raw line from the IRC protocol (minus the trailing CRLF)
// and returns a parsed Message and a bool indicating success.
//
// TODO(will): This does not currently validate the prefix portion of the
// message.
func parseMessage(line string) (message, bool) {
	var msg message

	// Parse the message into prefix, command and parameter strings.
	parts := messageRegex.FindStringSubmatch(line)
	if parts == nil {
		return msg, false
	}

	msg.prefix = parts[1]
	msg.command = parts[2]
	msg.params = make([]string, 0)

	// Split the parameters into separate strings.
	i := 0
	for i < len(parts[3])-1 {
		// Each parameter must begin with a space.
		if parts[3][i] != ' ' {
			return msg, false
		}
		// Skip extra white space.
		for parts[3][i] == ' ' {
			i++
		}

		// If the next character is a ':', then the parameter is the value of the
		// remainder of the string.
		if parts[3][i] == ':' {
			msg.trailing = parts[3][i+1:]
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
		msg.params = append(msg.params, parts[3][start:i])
	}

	return msg, true
}
