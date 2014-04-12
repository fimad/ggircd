package irc

import (
	"testing"
)

func TestUserHandlerPing(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-PING", state, handler, []handlerTest{
		{
			desc: "successful ping",
			in:   []Message{CmdPing},
			want: mockConnection{
				messages: []Message{
					CmdPong.WithPrefix("name").WithParams("name").WithTrailing("name"),
				},
			},
			state: newMockState().withUser("nick"),
		},
	})
}
