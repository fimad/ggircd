package irc

import (
	"testing"
)

func TestUserHandlerInvite(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-INVITE", state, handler, []handlerTest{
		{
			desc: "invite to current channel",
			in:   []Message{CmdInvite.WithParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyInviting,
					},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdInvite,
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
			in:   []Message{CmdInvite.WithParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorNotOnChannel,
					},
				},
				"foo": mockConnection{
					messages: []Message{},
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
			in:   []Message{CmdInvite.WithParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyInviting,
					},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdInvite,
					},
				},
			},
			state: newMockState().
				withUser("nick").
				withUser("foo"),
		},
		{
			desc: "op invite to current invite-only channel",
			in:   []Message{CmdInvite.WithParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyInviting,
					},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdInvite,
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
			in:   []Message{CmdInvite.WithParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorChanOPrivsNeeded,
					},
				},
				"foo": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "i", "").
				withUser("nick", "#channel").
				withUser("foo"),
		},
		{
			desc: "invite to not-joined invite-only channel",
			in:   []Message{CmdInvite.WithParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorNotOnChannel,
					},
				},
				"foo": mockConnection{
					messages: []Message{},
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
			in:   []Message{CmdInvite.WithParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorNoSuchNick,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel"),
		},
		{
			desc: "invite already on channel",
			in:   []Message{CmdInvite.WithParams("foo", "#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorUserOnChannel,
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
			in:   []Message{CmdInvite.WithParams("foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorNeedMoreParams,
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
