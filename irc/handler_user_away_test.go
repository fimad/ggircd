package irc

import (
	"testing"
)

func TestUserHandlerAway(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-AWAY", state, handler, []handlerTest{
		{
			desc: "set away",
			in:   []message{cmdAway.withTrailing("away right now")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{replyNowAway},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", "a"),
			},
		},
		{
			desc: "set away w/ no trailing",
			in:   []message{cmdAway.withParams("away")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{replyNowAway},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", "a"),
			},
		},
		{
			desc: "set un-away",
			in: []message{
				cmdAway.withTrailing("away right now"),
				cmdAway,
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyNowAway,
						replyUnaway,
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
			in: []message{
				cmdPrivMsg.withParams("foo").withTrailing("yo"),
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdPrivMsg,
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
