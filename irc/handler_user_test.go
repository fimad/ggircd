package irc

import (
	"testing"
)

func TestUserHandler(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newFreshHandler(state) }
	testHandler(t, "userHandler-EndToEnd", state, handler, []handlerTest{
		{
			desc: "connect to server, join a channel and set topic",
			in: []message{
				cmdNick.withParams("nick"),
				cmdUser.withParams("user", "host", "server").withTrailing("real name"),
				cmdJoin.withParams("#channel"),
				cmdTopic.withParams("#channel").withTrailing("blah blah blah"),
			},
			want: mockConnection{
				messages: []message{
					replyWelcome,
					replyYourHost,
					replyMOTDStart,
					replyMOTD,
					replyEndOfMOTD,
					replyNoTopic,
					replyNamReply,
					replyEndOfNames,
					cmdTopic,
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
			in: []message{
				cmdNick.withParams("nick"),
				cmdUser.withParams("user", "host", "server").withTrailing("real name"),
				cmdJoin.withParams("#channel"),
				cmdNick.withParams("foo"),
				cmdNick.withParams("bar"),
				cmdNick.withParams("baz"),
			},
			want: mockConnection{
				messages: []message{
					replyWelcome,
					replyYourHost,
					replyMOTDStart,
					replyEndOfMOTD,
					replyNoTopic,
					replyNamReply,
					replyEndOfNames,
					cmdNick,
					cmdNick,
					cmdNick,
				},
			},
			assert: []assert{
				assertChannelOp("#channel", "baz"),
			},
			state: newMockState(),
		},
	})
}
