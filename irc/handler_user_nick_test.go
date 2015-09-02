package irc

import (
	"testing"
)

func TestUserHandlerNicks(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-NAMES", state, handler, []handlerTest{
		{
			desc: "nick change successful",
			in:   []message{cmdNick.withParams("foo")},
			wantNicks: map[string]mockConnection{
				"foo": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "nick change broadcasts to channels",
			in:   []message{cmdNick.withParams("foo")},
			wantNicks: map[string]mockConnection{
				"foo": mockConnection{
					messages: []message{cmdNick},
				},
				"bar": mockConnection{
					messages: []message{cmdNick},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("bar", "#channel"),
		},
		{
			desc: "nick fails not given",
			in:   []message{cmdNick},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNoNicknameGiven},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "nick fails taken",
			in:   []message{cmdNick.withParams("foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNicknameInUse},
				},
			},
			state: newMockState().withUser("nick").withUser("foo"),
		},
		{
			desc: "nick fails too long",
			in:   []message{cmdNick.withParams("0123456789")},  // Max is 9 chars.
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNicknameInUse},
				},
			},
			state: newMockState().withUser("nick"),
		},
	})
}
