package irc

import (
	"testing"
)

func TestUserHandlerJoin(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-JOIN", state, handler, []handlerTest{
		{
			desc: "successful join",
			in:   []message{cmdJoin.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdJoin,
						replyNoTopic,
						replyNamReply,
						replyEndOfNames,
					},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "successful join, and then rejoin",
			in: []message{
				cmdJoin.withParams("#channel"),
				cmdJoin.withParams("#channel"),
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdJoin,
						replyNoTopic,
						replyNamReply,
						replyEndOfNames,
					},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "successful join with key",
			in:   []message{cmdJoin.withParams("#channel", "key")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdJoin,
						replyNoTopic,
						replyNamReply,
						replyEndOfNames,
					},
				},
			},
			state: newMockState().
				withUser("nick").
				withChannel("#channel", "k", "").
				withChannelKey("#channel", "key"),
		},
		{
			desc: "successful join limited channel",
			in:   []message{cmdJoin.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdJoin,
						replyNoTopic,
						replyNamReply,
						replyEndOfNames,
					},
				},
			},
			state: newMockState().
				withUser("nick").
				withChannel("#channel", "l", "").
				withChannelLimit("#channel", 1),
		},
		{
			desc: "successful join multiple",
			in:   []message{cmdJoin.withParams("#foo,#bar")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdJoin,
						replyNoTopic,
						replyNamReply,
						replyEndOfNames,
						cmdJoin,
						replyNoTopic,
						replyNamReply,
						replyEndOfNames,
					},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "no channel given",
			in:   []message{cmdJoin},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "bad channel name",
			in:   []message{cmdJoin.withParams("channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNoSuchChannel},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "invite only",
			in:   []message{cmdJoin.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorInviteOnlyChan},
				},
			},
			state: newMockState().withUser("nick").withChannel("#channel", "i", ""),
		},
		{
			desc: "wrong channel key",
			in:   []message{cmdJoin.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorBadChannelKey},
				},
			},
			state: newMockState().
				withUser("nick").
				withChannel("#channel", "k", "").
				withChannelKey("#channel", "key"),
		},
		{
			desc: "channel full",
			in:   []message{cmdJoin.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorChannelIsFull},
				},
			},
			state: newMockState().
				withUser("nick").
				withChannel("#channel", "l", "").
				withChannelLimit("#channel", 0),
		},
	})
}
