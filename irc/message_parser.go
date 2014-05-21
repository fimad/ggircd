package irc

import (
	"io"
	"net"
	"regexp"
)

// messageParser is a function that returns a Message and a boolean indicating
// if the end of the stream has been reached. If the boolean is false, then the
// returned Message should be ignored and the end of the input has been reached.
// It is also possible that there was no available message but the end of the
// stream has not been reached. This case is denoted by returning a message with
// an empty command field.
type messageParser func() (message, bool)

func newMessageParser(reader io.Reader) messageParser {
	buffer := make([]byte, 512)
	hasMore := true
	startOfBuffer := 0
	throwAway := false
	return func() (message, bool) {
		var message message

		// Check if there is already a newline character. If there is, then don't
		// bother waiting for another io op.
		hasNewLine := false
		for j := 0; j < startOfBuffer; j++ {
			if buffer[j] == '\n' || buffer[j] == '\r' {
				hasNewLine = true
				break
			}
		}

		// Only deal with network io if we don't have another message ready to go in
		// the buffer.
		var n int
		var err error
		if !hasNewLine {
			n, err = reader.Read(buffer[startOfBuffer:])

			// Handle the error. If it is a timeout or temporary then the connection
			// still has more data. Otherwise signal that there isn't any more data to
			// be had.
			if startOfBuffer == 0 && err != nil {
				netErr, castOk := err.(*net.OpError)
				if !castOk || (!netErr.Timeout() && !netErr.Temporary()) {
					hasMore = false
				}
				return message, hasMore
			}
		}

		// Scan over the buffer looking for some combination of CR-LF.
		var i int
		for i = 0; i < n+startOfBuffer; i++ {
			if buffer[i] != '\r' && buffer[i] != '\n' {
				continue
			}

			// We found a new line. If we are throwing away this message, don't throw
			// away the next one. Otherwise parse the message and store it in the
			// variable that we will eventually return.
			if throwAway {
				throwAway = false
			} else {
				// Parse the message up to but not including the new line character.
				message, _ = parseMessage(string(buffer[0:i]))
			}

			// Copy the remainder of the buffer to the head of the buffer. Do not
			// clobber i since it is used to determine if we ran out of space in the
			// buffer.
			k := i
			if i < len(buffer) && (buffer[i+1] == '\n' || buffer[i+1] == '\r') {
				k++
			}
			k++
			j := 0
			for k < n+startOfBuffer {
				buffer[j] = buffer[k]
				k++
				j++
			}
			startOfBuffer = j
			break
		}

		// If we scan the entire buffer and can't find a line terminating character,
		// then all data we read is thrown away until we read a new line.
		if i == len(buffer) {
			throwAway = true
			startOfBuffer = 0
		}

		return message, hasMore
	}
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
