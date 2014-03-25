package irc

import (
	"testing"
)

func TestUserHandlerNames(t *testing.T) {
	state := make(chan State, 1)
	testHandler(t, "UserHandler-NAMES", state, NewUserHandler(state, "nick"), []handlerTest{
		{
			desc: "names successful",
			in:   []Message{CmdNames.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyNamReply,
						ReplyNamReply,
						ReplyNamReply,
						ReplyEndOfNames,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel").
				withUser("bar", "#channel"),
		},
		{
			desc: "names all",
			in:   []Message{CmdNames.WithParams()},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyNamReply,
						ReplyEndOfNames,
						ReplyNamReply,
						ReplyEndOfNames,
					},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withChannel("#bar", "", "").
				withUser("nick").
				withUser("foo", "#foo").
				withUser("bar", "#bar"),
		},
		{
			desc: "names all secret",
			in:   []Message{CmdNames.WithParams()},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#foo", "s", "").
				withUser("nick").
				withUser("foo", "#foo"),
		},
		{
			desc: "names successful private",
			in:   []Message{CmdNames.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyNamReply,
						ReplyNamReply,
						ReplyNamReply,
						ReplyEndOfNames,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "p", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel").
				withUser("bar", "#channel"),
		},
		{
			desc: "names fails private",
			in:   []Message{CmdNames.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "p", "").
				withUser("nick").
				withUser("foo", "#channel").
				withUser("bar", "#channel"),
		},
	})
}
