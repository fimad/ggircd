package irc

import (
	"testing"
)

func TestUserHandlerWho(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-WHO", state, handler, []handlerTest{
		{
			desc: "who with no params",
			in:   []message{cmdWho},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick").
				withUser("foo").
				withUser("bar").
				withOps("#channel", "nick", "foo"),
		},
		{
			desc: "who w/ nick",
			in:   []message{cmdWho.withParams("foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyWhoReply.
							withPrefix("name").
							withParams("*", "foo", "foo", "name", "foo", "H@").
							withTrailing("0 foo"),
						replyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel"),
		},
		{
			desc: "who w/ channel (verify op)",
			in:   []message{cmdWho.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyWhoReply.
							withPrefix("name").
							withParams("#channel", "nick", "nick", "name", "nick", "H@").
							withTrailing("0 nick"),
						replyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withOps("#channel", "nick"),
		},
		{
			desc: "who w/ channel (verify voice)",
			in:   []message{cmdWho.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyWhoReply.
							withPrefix("name").
							withParams("#channel", "nick", "nick", "name", "nick", "H+").
							withTrailing("0 nick"),
						replyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withVoices("#channel", "nick"),
		},
		{
			desc: "who w/ channel (verify no op/voice)",
			in:   []message{cmdWho.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyWhoReply.
							withPrefix("name").
							withParams("#channel", "nick", "nick", "name", "nick", "H").
							withTrailing("0 nick"),
						replyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel"),
		},
		{
			desc: "who w/ channel",
			in:   []message{cmdWho.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyWhoReply,
						replyWhoReply,
						replyWhoReply,
						replyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("bar", "#channel").
				withUserMode("bar", "i").
				withUser("foo", "#channel"),
		},
		{
			desc: "who w/ channel and op",
			in:   []message{cmdWho.withParams("#channel", "o")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyWhoReply,
						replyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel").
				withOps("#channel", "nick"),
		},
	})
}
