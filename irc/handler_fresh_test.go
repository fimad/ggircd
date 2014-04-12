package irc

import (
	"testing"
)

func TestFreshHandlerHandle(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewFreshHandler(state) }
	testHandler(t, "FreshHandler", state, handler, []handlerTest{
		{
			desc:  "empty",
			in:    []Message{},
			want:  mockConnection{},
			state: newMockState(),
		},
		{
			desc:   "kill closed connection",
			in:     []Message{},
			want:   mockConnection{killed: true},
			state:  newMockState(),
			hangup: true,
		},
		{
			desc:  "nick without parameters",
			in:    []Message{CmdNick},
			want:  mockConnection{messages: []Message{ErrorNoNicknameGiven}},
			state: newMockState(),
		},
		{
			desc:  "nick, no user message",
			in:    []Message{CmdNick.WithParams("foo")},
			want:  mockConnection{},
			state: newMockState(),
		},
		{
			desc:  "nick using in-use nickname",
			in:    []Message{CmdNick.WithParams("foo")},
			want:  mockConnection{messages: []Message{ErrorNicknameInUse}},
			state: newMockState().withUser("foo"),
		},
		{
			desc: "user missing parameters",
			in: []Message{
				CmdNick.WithParams("foo"),
				CmdUser.WithTrailing("real name"),
			},
			want:  mockConnection{messages: []Message{ErrorNeedMoreParams}},
			state: newMockState(),
		},
		{
			desc: "successful registration",
			in: []Message{
				CmdNick.WithParams("foo"),
				CmdUser.WithParams("user", "host", "server").WithTrailing("real name"),
			},
			want: mockConnection{
				messages: []Message{
					ReplyWelcome,
					ReplyYourHost,
					ReplyMOTDStart,
					ReplyEndOfMOTD,
				},
			},
			state: newMockState(),
		},
	})
}
