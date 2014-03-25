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
	in []Message

	// The initial IRC state.
	state *mockState

	// The desired mock connection state for the connection that is passed to the
	// hnadler.
	want mockConnection

	// A map from nick names to desired mock connection states.
	wantNicks map[string]mockConnection

	// If false, perform a fuzzy comparison for messages where only the commands
	// are compared.
	strict bool

	// If true, call Handler.Close after the sequence of messages has been
	// processed.
	hangup bool
}

// testHandler is a helper method for use in testing handlers.
func testHandler(t *testing.T, name string, state chan State, handler Handler, tests []handlerTest) {
	for i, tt := range tests {
		state <- tt.state
		got := runHandler(tt, handler)
		_ = <-state

		gotNicks := make(map[string]mockConnection)
		ok := compareMessages(tt.strict, got, tt.want)
		for nick, want := range tt.wantNicks {
			user := tt.state.GetUser(nick)
			gotNicks[nick] = *user.Sink.(*mockConnection)
			ok = ok && compareMessages(tt.strict, gotNicks[nick], want)
		}

		if !ok {
			t.Errorf("%d. %s: %s\n%+v =>\n\tgot %+v\n\twant %+v\n\tgot nicks %+v\n\twant nicks %+v",
				i, name, tt.desc, tt.in, got, tt.want, gotNicks, tt.wantNicks)
		}
	}
}

func runHandler(tt handlerTest, handler Handler) mockConnection {
	conn := mockConnection{}
	for _, message := range tt.in {
		handler = handler.Handle(&conn, message)
	}
	if tt.hangup {
		handler.Closed(&conn)
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
		if got.messages[i].Command != want.messages[i].Command {
			return false
		}
	}

	return true
}
