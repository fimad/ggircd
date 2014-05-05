package irc

import (
	"testing"
)

func TestUserHandlerQuit(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-QUIT", state, handler, []handlerTest{
		{
			desc:  "successful quit",
			in:    []message{cmdQuit},
			want:  mockConnection{killed: true},
			state: newMockState().withUser("nick"),
		},
		{
			desc:  "handles no messages after quiting",
			in:    []message{cmdQuit, cmdPing},
			want:  mockConnection{killed: true},
			state: newMockState().withUser("nick"),
		},
	})
}
