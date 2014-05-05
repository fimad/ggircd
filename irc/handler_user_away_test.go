package irc

import (
	"testing"
)

func TestUserHandlerAway(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-AWAY", state, handler, []handlerTest{
		{
			desc: "set away",
			in:   []Message{CmdAway.WithTrailing("away right now")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ReplyNowAway},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", "a"),
			},
		},
		{
			desc: "set un-away",
			in: []Message{
				CmdAway.WithTrailing("away right now"),
				CmdAway,
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyNowAway,
						ReplyUnaway,
					},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", ""),
			},
		},
		{
			desc: "automatic priv replies",
			in: []Message{
				CmdPrivMsg.WithParams("foo").WithTrailing("yo"),
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						CmdPrivMsg,
					},
				},
			},
			state: newMockState().
				withUser("nick").
				withUser("foo").
				withUserAway("foo", "not here right now"),
		},
	})
}
