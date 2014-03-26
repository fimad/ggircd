package irc

import (
	"testing"
)

func TestUserHandlerQuit(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-QUIT", state, handler, []handlerTest{
		{
			desc:  "successful quit",
			in:    []Message{CmdQuit},
			want:  mockConnection{killed: true},
			state: newMockState().withUser("nick"),
		},
		{
			desc:  "handles no messages after quiting",
			in:    []Message{CmdQuit, CmdPing},
			want:  mockConnection{killed: true},
			state: newMockState().withUser("nick"),
		},
	})
}
