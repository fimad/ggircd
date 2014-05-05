package irc

import (
	"testing"
)

func TestUserHandlerInvite(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-INVITE", state, handler, []handlerTest{
		{
			desc: "invite to current channel",
			in:   []message{cmdInvite.withParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyInviting,
					},
				},
				"foo": mockConnection{
					messages: []message{
						cmdInvite,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo"),
		},
		{
			desc: "invite to not-joined channel",
			in:   []message{cmdInvite.withParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorNotOnChannel,
					},
				},
				"foo": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick").
				withUser("bar", "#channel").
				withUser("foo"),
		},
		{
			desc: "invite to non-existent channel",
			in:   []message{cmdInvite.withParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyInviting,
					},
				},
				"foo": mockConnection{
					messages: []message{
						cmdInvite,
					},
				},
			},
			state: newMockState().
				withUser("nick").
				withUser("foo"),
		},
		{
			desc: "op invite to current invite-only channel",
			in:   []message{cmdInvite.withParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyInviting,
					},
				},
				"foo": mockConnection{
					messages: []message{
						cmdInvite,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "i", "").
				withUser("nick", "#channel").
				withOps("#channel", "nick").
				withUser("foo"),
		},
		{
			desc: "non-op invite to current invite-only channel",
			in:   []message{cmdInvite.withParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorChanOPrivsNeeded,
					},
				},
				"foo": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "i", "").
				withUser("nick", "#channel").
				withUser("foo"),
		},
		{
			desc: "invite to not-joined invite-only channel",
			in:   []message{cmdInvite.withParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorNotOnChannel,
					},
				},
				"foo": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "i", "").
				withUser("nick").
				withUser("bar", "#channel").
				withUser("foo"),
		},
		{
			desc: "invite bad nick",
			in:   []message{cmdInvite.withParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorNoSuchNick,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel"),
		},
		{
			desc: "invite already on channel",
			in:   []message{cmdInvite.withParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorUserOnChannel,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel"),
		},
		{
			desc: "invite with no channel",
			in:   []message{cmdInvite.withParams("foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorNeedMoreParams,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", ""),
		},
		// TODO(will): Add a case for away once that is supported.
	})
}
