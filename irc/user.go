package irc

import (
	"fmt"
)

type user struct {
	nick     string
	user     string
	host     string
	server   string
	realName string

	channels map[*channel]bool

	awayMessage string
	mode        mode

	sink sink
}

func (u user) send(msg message) {
	u.sink.send(msg)
}

// forChannels iterates over all of the channels that the user has joined and
// passes a pointer to each to the supplied callback.
func (u user) forChannels(callback func(*channel)) {
	for ch := range u.channels {
		callback(ch)
	}
}

// prefix returns the prefix string that should be used in messages originating
// from this user.
func (u *user) prefix() string {
	return fmt.Sprintf("%s!%s@%s", u.nick, u.user, u.host)
}
