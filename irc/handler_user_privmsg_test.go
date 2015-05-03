package irc

import (
	"testing"
)

func TestUserHandlerPrivMsg(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-PRIVMSG", state, handler, []handlerTest{
		{
			desc:   "successful privmsg user",
			in:     []message{cmdPrivMsg.withParams("foo").withTrailing("msg")},
			strict: true,
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
				"foo": mockConnection{
					messages: []message{
						cmdPrivMsg.
							withPrefix("nick!nick@nick").
							withParams("foo").
							withTrailing("msg"),
					},
				},
			},
			state: newMockState().
				withUser("nick").
				withUser("foo"),
		},
		{
			desc:   "successful privmsg user w/ no trailing",
			in:     []message{cmdPrivMsg.withParams("foo", "msg")},
			strict: true,
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
				"foo": mockConnection{
					messages: []message{
						cmdPrivMsg.
							withPrefix("nick!nick@nick").
							withParams("foo").
							withTrailing("msg"),
					},
				},
			},
			state: newMockState().
				withUser("nick").
				withUser("foo"),
		},
		{
			desc:   "successful privmsg channel",
			in:     []message{cmdPrivMsg.withParams("#channel").withTrailing("msg")},
			strict: true,
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
				"foo": mockConnection{
					messages: []message{
						cmdPrivMsg.
							withPrefix("nick!nick@nick").
							withParams("#channel").
							withTrailing("msg"),
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
			in:   []message{cmdPrivMsg.withParams("#channel").withTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
				"foo": mockConnection{
					messages: []message{cmdPrivMsg},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick").
				withUser("foo", "#channel"),
		},
		{
			desc: "successful privmsg moderated channel with voice",
			in:   []message{cmdPrivMsg.withParams("#channel").withTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
				"foo": mockConnection{
					messages: []message{cmdPrivMsg},
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
			in:   []message{cmdPrivMsg.withParams("#channel").withTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
				"foo": mockConnection{
					messages: []message{cmdPrivMsg},
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
			in:   []message{cmdPrivMsg.withTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNoRecipient},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - no message",
			in:   []message{cmdPrivMsg.withParams("foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNoTextToSend},
				},
			},
			state: newMockState().
				withUser("nick").
				withUser("foo"),
		},
		{
			desc: "failure - invalid user",
			in:   []message{cmdPrivMsg.withParams("foo").withTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNoSuchNick},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - non-existent channel",
			in:   []message{cmdPrivMsg.withParams("#channel").withTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorCannotSendToChan},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - message from outside to +n channel",
			in:   []message{cmdPrivMsg.withParams("#channel").withTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorCannotSendToChan},
				},
				"foo": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "n", "").
				withUser("foo", "#channel").
				withUser("nick"),
		},
		{
			desc: "failure - non-voice messaging moderated channel",
			in:   []message{cmdPrivMsg.withParams("#channel").withTrailing("msg")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorCannotSendToChan},
				},
				"foo": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "m", "").
				withUser("foo", "#channel").
				withUser("nick", "#channel"),
		},
	})
}
