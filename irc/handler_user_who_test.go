package irc

import (
	"testing"
)

func TestUserHandlerWho(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-WHO", state, handler, []handlerTest{
		{
			desc: "who with no params",
			in:   []Message{CmdWho},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyEndOfWho,
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
			in:   []Message{CmdWho.WithParams("foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyWhoReply.
							WithPrefix("name").
							WithParams("*", "foo", "foo", "name", "foo", "H@").
							WithTrailing("0 foo"),
						ReplyEndOfWho,
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
			in:   []Message{CmdWho.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyWhoReply.
							WithPrefix("name").
							WithParams("#channel", "nick", "nick", "name", "nick", "H@").
							WithTrailing("0 nick"),
						ReplyEndOfWho,
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
			in:   []Message{CmdWho.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyWhoReply.
							WithPrefix("name").
							WithParams("#channel", "nick", "nick", "name", "nick", "H+").
							WithTrailing("0 nick"),
						ReplyEndOfWho,
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
			in:   []Message{CmdWho.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyWhoReply.
							WithPrefix("name").
							WithParams("#channel", "nick", "nick", "name", "nick", "H").
							WithTrailing("0 nick"),
						ReplyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel"),
		},
		{
			desc: "who w/ channel",
			in:   []Message{CmdWho.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyWhoReply,
						ReplyWhoReply,
						ReplyEndOfWho,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel"),
		},
		{
			desc: "who w/ channel and op",
			in:   []Message{CmdWho.WithParams("#channel", "o")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyWhoReply,
						ReplyEndOfWho,
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
