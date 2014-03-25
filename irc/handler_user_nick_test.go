package irc

import (
	"testing"
)

func TestUserHandlerNicks(t *testing.T) {
	state := make(chan State, 1)
	testHandler(t, "UserHandler-NAMES", state, NewUserHandler(state, "nick"), []handlerTest{
		{
			desc: "nick change succesful",
			in:   []Message{CmdNick.WithParams("foo")},
			wantNicks: map[string]mockConnection{
				"foo": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "nick change broadcasts to channels",
			in:   []Message{CmdNick.WithParams("foo")},
			wantNicks: map[string]mockConnection{
				"foo": mockConnection{
					messages: []Message{CmdNick},
				},
				"bar": mockConnection{
					messages: []Message{CmdNick},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("bar", "#channel"),
		},
		{
			desc: "nick fails not given",
			in:   []Message{CmdNick},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNoNicknameGiven},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "nick fails taken",
			in:   []Message{CmdNick.WithParams("foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNicknameInUse},
				},
			},
			state: newMockState().withUser("nick").withUser("foo"),
		},
	})
}
