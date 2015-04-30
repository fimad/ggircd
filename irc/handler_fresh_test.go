package irc

import (
	"testing"
)

func TestFreshHandlerHandle(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newFreshHandler(state) }
	testHandler(t, "FreshHandler", state, handler, []handlerTest{
		{
			desc:  "empty",
			in:    []message{},
			want:  mockConnection{},
			state: newMockState(),
		},
		{
			desc:   "kill closed connection",
			in:     []message{},
			want:   mockConnection{killed: true},
			state:  newMockState(),
			hangup: true,
		},
		{
			desc:  "nick without parameters",
			in:    []message{cmdNick},
			want:  mockConnection{messages: []message{errorNoNicknameGiven}},
			state: newMockState(),
		},
		{
			desc:  "nick, no user message",
			in:    []message{cmdNick.withParams("foo")},
			want:  mockConnection{},
			state: newMockState(),
		},
		{
			desc:  "nick using in-use nickname",
			in:    []message{cmdNick.withParams("foo")},
			want:  mockConnection{messages: []message{errorNicknameInUse}},
			state: newMockState().withUser("foo"),
		},
		{
			desc: "user missing parameters",
			in: []message{
				cmdNick.withParams("foo"),
				cmdUser.withTrailing("real name"),
			},
			want:  mockConnection{messages: []message{errorNeedMoreParams}},
			state: newMockState(),
		},
		{
			desc: "successful registration",
			in: []message{
				cmdNick.withParams("foo"),
				cmdUser.withParams("user", "host", "server").withTrailing("real name"),
			},
			want: mockConnection{
				messages: []message{
					replyWelcome,
					replyYourHost,
					replyMOTDStart,
					replyEndOfMOTD,
				},
			},
			state: newMockState(),
		},
		{
			desc: "successful registration w/ no trailing",
			in: []message{
				cmdNick.withParams("foo"),
				cmdUser.withParams("user", "host", "server", "realname"),
			},
			want: mockConnection{
				messages: []message{
					replyWelcome,
					replyYourHost,
					replyMOTDStart,
					replyEndOfMOTD,
				},
			},
			state: newMockState(),
		},
	})
}
