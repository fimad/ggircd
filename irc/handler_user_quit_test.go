package irc

import (
	"testing"
)

func TestUserHandlerQuit(t *testing.T) {
	state := make(chan State, 1)
	testHandler(t, "UserHandler-QUIT", state, NewUserHandler(state, "nick"), []handlerTest{
		{
			desc:  "succesful quit",
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
