package irc

import (
	"reflect"
	"testing"
)

type handlerTest struct {
	// A string describing the behaviour that this particular test is testing.
	desc string

	// A sequence of messages that will be sent from the nick supplied to
	// testHandler.
	in []message

	// The initial IRC state.
	state *mockState

	// The desired mock connection state for the connection that is passed to the
	// hnadler.
	want mockConnection

	// A map from nick names to desired mock connection states.
	wantNicks map[string]mockConnection

	// A slice of assertions that will be applied to the resulting state.
	assert []assert

	// If false, perform a fuzzy comparison for messages where only the commands
	// are compared.
	strict bool

	// If true, call Handler.Close after the sequence of messages has been
	// processed.
	hangup bool

	// Override the message of the day so that it is not loaded from disk.
	motd []string
}

// testHandler is a helper method for use in testing handlers.
func testHandler(t *testing.T, name string, state chan state, handler func() handler, tests []handlerTest) {
	for i, tt := range tests {
		if tt.motd == nil {
			motd = []string{}
		} else {
			motd = tt.motd
		}

		state <- tt.state
		got := runHandler(tt, handler())
		_ = <-state

		gotNicks := make(map[string]mockConnection)
		if !compareMessages(tt.strict, got, tt.want) {
			t.Errorf("%d. %s: %s\n%+v =>\n\tgot %+v\n\twant %+v",
				i, name, tt.desc, tt.in, got, tt.want)
		}

		for nick, want := range tt.wantNicks {
			user := tt.state.getUser(nick)
			gotNicks[nick] = *user.sink.(*mockConnection)
			if !compareMessages(tt.strict, gotNicks[nick], want) {
				t.Errorf("%d. %s: %s\n%+v => nick = %q\n\tgot %+v\n\twant %+v",
					i, name, tt.desc, tt.in, nick, gotNicks[nick], tt.wantNicks[nick])
			}
		}

		for j, assert := range tt.assert {
			err := assert(tt.state)
			if err != nil {
				t.Errorf("%d. %s: %s\n%+v =>\n\tfailed assert (%d) %s",
					i, name, tt.desc, tt.in, j, err)
			}
		}
	}
}

func runHandler(tt handlerTest, handler handler) mockConnection {
	conn := mockConnection{}
	for _, message := range tt.in {
		handler = handler.handle(&conn, message)
	}
	if tt.hangup {
		handler.closed(&conn)
	}
	return conn
}

func compareMessages(strict bool, got, want mockConnection) bool {
	if got.killed != want.killed {
		return false
	}

	// Special case to handle 0 length slice. It is possible to have an implicit
	// empty in want, but an explicit 0 length slice in want. These should
	// evaluate to equivalent.
	if len(got.messages) == 0 && len(want.messages) == 0 {
		return true
	}

	if strict {
		return reflect.DeepEqual(got.messages, want.messages)
	}

	if len(got.messages) != len(want.messages) {
		return false
	}

	for i := 0; i < len(got.messages); i++ {
		if got.messages[i].command != want.messages[i].command {
			return false
		}
	}

	return true
}
