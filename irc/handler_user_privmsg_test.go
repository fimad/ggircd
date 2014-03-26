package irc

import (
	"testing"
)

func TestUserHandlerPrivMsg(t *testing.T) {
	state := make(chan State, 1)
	testHandler(t, "UserHandler-PRIVMSG", state, NewUserHandler(state, "nick"), []handlerTest{
		{
			desc:   "successful privmsg user",
			in:     []Message{CmdPrivMsg.WithParams("foo").WithTrailing("msg")},
			strict: true,
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdPrivMsg.
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
			desc:   "successful privmsg channel",
			in:     []Message{CmdPrivMsg.WithParams("#channel").WithTrailing("msg")},
			strict: true,
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						CmdPrivMsg.
							WithPrefix("nick!nick@nick").
							WithParams("#channel").
							WithTrailing("msg"),
					},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdPrivMsg.
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
			desc: "successful privmsg -n channel from outside",
			in:   []Message{CmdPrivMsg.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{CmdPrivMsg},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick").
				withUser("foo", "#channel"),
		},
		{
			desc: "successful privmsg moderated channel with voice",
			in:   []Message{CmdPrivMsg.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{CmdPrivMsg},
				},
				"foo": mockConnection{
					messages: []Message{CmdPrivMsg},
				},
			},
			state: newMockState().
				withChannel("#channel", "m", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel").
				withVoices("#channel", "nick"),
		},
		{
			desc: "successful privmsg moderated channel with op",
			in:   []Message{CmdPrivMsg.WithParams("#channel").WithTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{CmdPrivMsg},
				},
				"foo": mockConnection{
					messages: []Message{CmdPrivMsg},
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
			in:   []Message{CmdPrivMsg.WithTrailing("msg")},
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
			in:   []Message{CmdPrivMsg.WithParams("foo")},
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
			in:   []Message{CmdPrivMsg.WithParams("foo").WithTrailing("msg")},
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
			in:   []Message{CmdPrivMsg.WithParams("#channel").WithTrailing("msg")},
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
			in:   []Message{CmdPrivMsg.WithParams("#channel").WithTrailing("msg")},
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
			in:   []Message{CmdPrivMsg.WithParams("#channel").WithTrailing("msg")},
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
