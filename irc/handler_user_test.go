package irc

import (
	"testing"
)

func TestUserHandler(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewFreshHandler(state) }
	testHandler(t, "UserHandler-EndToEnd", state, handler, []handlerTest{
		{
			desc: "connect to server, join a channel and set topic",
			in: []Message{
				CmdNick.WithParams("nick"),
				CmdUser.WithParams("user", "host", "server").WithTrailing("real name"),
				CmdJoin.WithParams("#channel"),
				CmdTopic.WithParams("#channel").WithTrailing("blah blah blah"),
			},
			want: mockConnection{
				messages: []Message{
					ReplyWelcome,
					ReplyYourHost,
					ReplyMOTDStart,
					ReplyMOTD,
					ReplyEndOfMOTD,
					ReplyNoTopic,
					ReplyNamReply,
					ReplyEndOfNames,
					CmdTopic,
				},
			},
			assert: []assert{
				assertChannelOp("#channel", "nick"),
			},
			state: newMockState(),
			motd:  []string{"foobar"},
		},
		{
			desc: "connect to server, connect to channel, and change nick",
			in: []Message{
				CmdNick.WithParams("nick"),
				CmdUser.WithParams("user", "host", "server").WithTrailing("real name"),
				CmdJoin.WithParams("#channel"),
				CmdNick.WithParams("foo"),
				CmdNick.WithParams("bar"),
				CmdNick.WithParams("baz"),
			},
			want: mockConnection{
				messages: []Message{
					ReplyWelcome,
					ReplyYourHost,
					ReplyMOTDStart,
					ReplyEndOfMOTD,
					ReplyNoTopic,
					ReplyNamReply,
					ReplyEndOfNames,
					CmdNick,
					CmdNick,
					CmdNick,
				},
			},
			assert: []assert{
				assertChannelOp("#channel", "baz"),
			},
			state: newMockState(),
		},
	})
}
