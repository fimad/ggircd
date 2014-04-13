package irc

import (
	"testing"
)

func TestUserHandlerKick(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-KICK", state, handler, []handlerTest{
		{
			desc: "op kick one user, one channel",
			in:   []Message{CmdKick.WithParams("#channel", "foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						CmdKick,
					},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdKick,
					},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withOps("#channel", "nick").
				withUser("foo", "#channel"),
			assert: []assert{
				assertUserNotOnChannel("foo", "#channel"),
			},
		},
		{
			desc: "op kick many user, one channels",
			in:   []Message{CmdKick.WithParams("#a", "foo,bar,baz")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						CmdKick,
						CmdKick,
						CmdKick,
					},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdKick,
					},
				},
				"bar": mockConnection{
					messages: []Message{
						CmdKick,
						CmdKick,
					},
				},
				"baz": mockConnection{
					messages: []Message{
						CmdKick,
						CmdKick,
						CmdKick,
					},
				},
			},
			state: newMockState().
				withChannel("#a", "", "").
				withUser("nick", "#a").
				withOps("#a", "nick").
				withUser("foo", "#a").
				withUser("bar", "#a").
				withUser("baz", "#a"),
			assert: []assert{
				assertUserNotOnChannel("foo", "#a"),
				assertUserNotOnChannel("bar", "#a"),
				assertUserNotOnChannel("baz", "#a"),
			},
		},
		{
			desc: "op kick many users, many channels",
			in:   []Message{CmdKick.WithParams("#a,#b", "foo,bar")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						CmdKick,
						CmdKick,
					},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdKick,
						CmdKick,
					},
				},
				"bar": mockConnection{
					messages: []Message{
						CmdKick,
						CmdKick,
					},
				},
			},
			state: newMockState().
				withChannel("#a", "", "").
				withChannel("#b", "", "").
				withUser("nick", "#a", "#b").
				withOps("#a", "nick").
				withOps("#b", "nick").
				withUser("foo", "#a", "#b").
				withUser("bar", "#a", "#b"),
			assert: []assert{
				assertUserNotOnChannel("foo", "#a"),
				assertUserNotOnChannel("bar", "#b"),
				assertUserOnChannel("foo", "#b"),
				assertUserOnChannel("bar", "#a"),
			},
		},
		{
			desc: "non-op kick",
			in:   []Message{CmdKick.WithParams("#a", "foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorChanOPrivsNeeded},
				},
				"foo": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#a", "", "").
				withUser("nick", "#a").
				withUser("foo", "#a"),
			assert: []assert{
				assertUserOnChannel("foo", "#a"),
			},
		},
		{
			desc: "mixed op kick",
			in:   []Message{CmdKick.WithParams("#a,#b,#c", "foo,foo,foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						CmdKick,
						ErrorChanOPrivsNeeded,
						CmdKick,
					},
				},
				"foo": mockConnection{
					messages: []Message{
						CmdKick,
						CmdKick,
					},
				},
			},
			state: newMockState().
				withChannel("#a", "", "").
				withChannel("#b", "", "").
				withChannel("#c", "", "").
				withUser("nick", "#a", "#b", "#c").
				withOps("#a", "nick").
				withOps("#c", "nick").
				withUser("foo", "#a", "#b", "#c"),
			assert: []assert{
				assertUserNotOnChannel("foo", "#a"),
				assertUserNotOnChannel("foo", "#c"),
				assertUserOnChannel("foo", "#b"),
			},
		},
		{
			desc: "kick with unequal users and channels",
			in:   []Message{CmdKick.WithParams("#a,#b,#c", "foo,bar")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorNeedMoreParams,
					},
				},
				"foo": mockConnection{
					messages: []Message{},
				},
				"bar": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#a", "", "").
				withChannel("#b", "", "").
				withChannel("#c", "", "").
				withUser("nick", "#a", "#b", "#c").
				withOps("#a", "nick").
				withOps("#b", "nick").
				withOps("#c", "nick").
				withUser("foo", "#a", "#b", "#c").
				withUser("bar", "#a", "#b", "#c"),
			assert: []assert{
				assertUserOnChannel("foo", "#a"),
				assertUserOnChannel("foo", "#b"),
				assertUserOnChannel("foo", "#c"),
				assertUserOnChannel("bar", "#a"),
				assertUserOnChannel("bar", "#b"),
				assertUserOnChannel("bar", "#c"),
			},
		},
		{
			desc: "kick user not on channel",
			in:   []Message{CmdKick.WithParams("#a", "foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorUserNotInChannel,
					},
				},
				"foo": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#a", "", "").
				withUser("nick", "#a").
				withOps("#a", "nick").
				withUser("foo"),
		},
		{
			desc: "kick bad channel",
			in:   []Message{CmdKick.WithParams("#b", "foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ErrorNoSuchChannel,
					},
				},
				"foo": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#a", "", "").
				withUser("nick", "#a").
				withOps("#a", "nick").
				withUser("foo", "#a"),
			assert: []assert{
				assertUserOnChannel("foo", "#a"),
			},
		},
		{
			desc: "non-op kick",
			in:   []Message{CmdKick.WithParams("#a", "foo")},
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
				withChannel("#a", "", "").
				withUser("nick", "#a").
				withUser("foo", "#a"),
			assert: []assert{
				assertUserOnChannel("foo", "#a"),
			},
		},
		{
			desc: "kick while not on channel",
			in:   []Message{CmdKick.WithParams("#a", "foo")},
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
				withChannel("#a", "", "").
				withUser("nick").
				withUser("foo", "#a"),
			assert: []assert{
				assertUserOnChannel("foo", "#a"),
			},
		},
	})
}
