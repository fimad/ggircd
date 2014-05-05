package irc

import (
	"strings"
)

type message struct {
	prefix   string
	command  string
	params   []string
	trailing string
}

// withParams creates a new copy of a message with the given parameters.
func (m message) withParams(params ...string) message {
	m.params = params
	return m
}

// withParams creates a new copy of a message with the given parameters.
func (m message) withTrailing(trailing string) message {
	m.trailing = trailing
	return m
}

// withPrefix creates a new copy of a message with the given prefix.
func (m message) withPrefix(prefix string) message {
	m.prefix = prefix
	return m
}

// toString serializes a Message to an IRC protocol compatible string.
func (m message) toString() (string, bool) {
	if m.command == "" {
		return "", false
	}

	var msg string
	if len(m.prefix) > 0 {
		msg = ":" + m.prefix + " "
	}

	msg += m.command

	for i := 0; i < len(m.params); i++ {
		param := m.params[i]
		if strings.Index(param, " ") != -1 {
			return "", false
		}
		msg += " " + param
	}

	if m.trailing != "" {
		msg += " :" + m.trailing
	}

	msg += "\r\n"

	return msg, true
}
