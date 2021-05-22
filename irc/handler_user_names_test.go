package irc

import (
	"testing"
)

func TestUserHandlerNames(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-NAMES", state, handler, []handlerTest{
		{
			desc: "names successful",
			in:   []message{cmdNames.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyNamReply,
						replyNamReply,
						replyNamReply,
						replyNamReply,
						replyEndOfNames,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel").
				withUser("bar", "#channel").
				// User baz should be listed even though they are invisible because
				// they share a channel with the user that is requesting names.
				withUser("baz", "#channel").
				withUserMode("baz", "i"),
		},
		{
			desc: "names all",
			in:   []message{cmdNames.withParams()},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyNamReply,
						replyEndOfNames,
						replyNamReply,
						replyEndOfNames,
					},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withChannel("#bar", "", "").
				withUser("nick").
				withUser("foo", "#foo").
				withUser("bar", "#bar").
				// User baz should not be listed because their are invisible.
				withUser("baz", "#bar").
				withUserMode("baz", "i"),
		},
		{
			desc: "names all private",
			in:   []message{cmdNames.withParams()},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#foo", "p", "").
				withUser("nick").
				withUser("foo", "#foo"),
		},
		{
			desc: "names successful private",
			in:   []message{cmdNames.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyNamReply,
						replyNamReply,
						replyNamReply,
						replyEndOfNames,
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
			in:   []message{cmdNames.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "p", "").
				withUser("nick").
				withUser("foo", "#channel").
				withUser("bar", "#channel"),
		},
		{
			desc: "names all secret",
			in:   []message{cmdNames.withParams()},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#foo", "s", "").
				withUser("nick").
				withUser("foo", "#foo"),
		},
		{
			desc: "names successful secret",
			in:   []message{cmdNames.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyNamReply,
						replyNamReply,
						replyNamReply,
						replyEndOfNames,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "s", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel").
				withUser("bar", "#channel"),
		},
		{
			desc: "names fails secret",
			in:   []message{cmdNames.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "s", "").
				withUser("nick").
				withUser("foo", "#channel").
				withUser("bar", "#channel"),
		},
	})
}
