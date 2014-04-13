package irc

import (
	"testing"
)

func TestUserHandlerNotice(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-NOTICE", state, handler, []handlerTest{
		{
			desc:   "successful notice user",
			in:     []Message{CmdNotice.WithParams("foo").WithTrailing("msg")},
			strict: true,
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdNotice.
							WithPrefix("nick!nick@nick").
							WithParams("foo").
							WithTrailing("msg"),
					},
				},
			},
			state: newMockState().
				withUser("nick").
				withUser("foo"),
		},
		{
			desc:   "successful notice channel",
			in:     []Message{CmdNotice.WithParams("#channel").WithTrailing("msg")},
			strict: true,
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdNotice.
							WithPrefix("nick!nick@nick").
							WithParams("#channel").
							WithTrailing("msg"),
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel"),
		},
		{
			desc: "successful notice -n channel from outside",
			in:   []Message{CmdNotice.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{CmdNotice},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick").
				withUser("foo", "#channel"),
		},
		{
			desc: "successful notice moderated channel with voice",
			in:   []Message{CmdNotice.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{CmdNotice},
				},
			},
			state: newMockState().
				withChannel("#channel", "m", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel").
				withVoices("#channel", "nick"),
		},
		{
			desc: "successful notice moderated channel with op",
			in:   []Message{CmdNotice.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{CmdNotice},
				},
			},
			state: newMockState().
				withChannel("#channel", "m", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel").
				withOps("#channel", "nick"),
		},
		{
			desc: "failure - no channel or user",
			in:   []Message{CmdNotice.WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNoRecipient},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - no message",
			in:   []Message{CmdNotice.WithParams("foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNoTextToSend},
				},
			},
			state: newMockState().
				withUser("nick").
				withUser("foo"),
		},
		{
			desc: "failure - invalid user",
			in:   []Message{CmdNotice.WithParams("foo").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNoSuchNick},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - non-existent channel",
			in:   []Message{CmdNotice.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorCannotSendToChan},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - message from outside to +n channel",
			in:   []Message{CmdNotice.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorCannotSendToChan},
				},
				"foo": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "n", "").
				withUser("foo", "#channel").
				withUser("nick"),
		},
		{
			desc: "failure - non-voice messaging moderated channel",
			in:   []Message{CmdNotice.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorCannotSendToChan},
				},
				"foo": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "m", "").
				withUser("foo", "#channel").
				withUser("nick", "#channel"),
		},
	})
}
