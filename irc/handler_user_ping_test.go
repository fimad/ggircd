package irc

import (
	"testing"
)

func TestUserHandlerPing(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-PING", state, handler, []handlerTest{
		{
			desc: "successful ping",
			in:   []message{cmdPing},
			want: mockConnection{
				messages: []message{
					cmdPong.withPrefix("name").withParams("name").withTrailing("name"),
				},
			},
			state: newMockState().withUser("nick"),
		},
	})
}
