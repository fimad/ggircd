package irc

import (
	"bufio"
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
	bufReader := bufio.NewReader(reader)
	buffer := make([]byte, 1)
	parseBuffer := make([]byte, 0, 512)

	var newConsumeNewLineState func(func() (message, bool)) func() (message, bool)
	var newParseState func() func() (message, bool)
	var throwAwayState func() (message, bool)
	var currentState func() (message, bool)

	hasMore := func(err error) bool {
		netErr, castOk := err.(*net.OpError)
		return castOk && (netErr.Timeout() || netErr.Temporary())
	}

	// Consumes input while the current character is a newline character. Then it
	// transitions to the supplied next state.
	newConsumeNewLineState = func(nextState func() (message, bool)) func() (message, bool) {
		var fn func() (message, bool)
		fn = func() (message, bool) {
			buf, err := bufReader.Peek(1)
			if err != nil {
				return message{}, hasMore(err)
			}

			// If the character is a new line character, keep reading.
			if buf[0] == '\r' || buf[0] == '\n' {
				bufReader.Read(buffer)
				return fn()
			}

			currentState = nextState
			return currentState()
		}

		return fn
	}

	// A state that corresponds to consuming input until a newline character is
	// encountered. Then the state is switched to parsing.
	throwAwayState = func() (message, bool) {
		_, err := bufReader.Read(buffer)
		if err != nil {
			return message{}, hasMore(err)
		}

		// If the character is a new line character, switch to consume newline state
		// followed by a parsing state.
		if buffer[0] == '\r' || buffer[0] == '\n' {
			currentState = newConsumeNewLineState(newParseState())
			return currentState()
		}

		return throwAwayState()
	}

	// Returns a new state that will read a line up to 512 characters long
	// including a terminating newline character and parse it into a message.
	newParseState = func() func() (message, bool) {
		// Reset the parse buffer.
		parseBuffer = parseBuffer[0:0]

		var fn func() (message, bool)
		fn = func() (message, bool) {
			// If the parse buffer has filled up, then throwaway input until a new
			// line character is encountered.
			if len(parseBuffer) == 512 {
				currentState = throwAwayState
				return currentState()
			}

			_, err := bufReader.Read(buffer)
			if err != nil {
				return message{}, hasMore(err)
			}

			// If the character is a new line character, parse the message, and switch
			// to consume newline state followed by a parsing state.
			if buffer[0] == '\r' || buffer[0] == '\n' {
				// The ordering matters here since the same parseBuffer is used for
				// every parse state.
				message, _ := parseMessage(string(parseBuffer))
				currentState = newConsumeNewLineState(newParseState())
				return message, true
			}

			// Save the current character to the parse buffer and keep reading.
			parseBuffer = append(parseBuffer, buffer[0])
			return fn()
		}

		return fn
	}

	// Return a wrapper function that evaluates the function corresponding to the
	// current state.
	currentState = newParseState()
	return func() (message, bool) {
		return currentState()
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
