package irc

import (
	"fmt"
)

type User struct {
	Nick     string
	User     string
	Host     string
	Server   string
	RealName string

	Channels map[*Channel]bool

	Mode Mode

	Sink Sink
}

func (u User) Send(msg Message) {
	u.Sink.Send(msg)
}

// ForChannels iterates over all of the channels that the user has joined and
// passes a pointer to each to the supplied callback.
func (u User) ForChannels(callback func(*Channel)) {
	for ch := range u.Channels {
		callback(ch)
	}
}

// Prefix returns the prefix string that should be used in Messages originating
// from this user.
func (u *User) Prefix() string {
	return fmt.Sprintf("%s!%s@%s", u.Nick, u.User, u.Host)
}
