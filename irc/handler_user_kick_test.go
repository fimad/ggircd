package irc

import (
	"testing"
)

func TestUserHandlerKick(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-KICK", state, handler, []handlerTest{
		{
			desc: "op kick one user, one channel",
			in:   []message{cmdKick.withParams("#channel", "foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdKick,
					},
				},
				"foo": mockConnection{
					messages: []message{
						cmdKick,
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
			in:   []message{cmdKick.withParams("#a", "foo,bar,baz")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdKick,
						cmdKick,
						cmdKick,
					},
				},
				"foo": mockConnection{
					messages: []message{
						cmdKick,
					},
				},
				"bar": mockConnection{
					messages: []message{
						cmdKick,
						cmdKick,
					},
				},
				"baz": mockConnection{
					messages: []message{
						cmdKick,
						cmdKick,
						cmdKick,
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
			in:   []message{cmdKick.withParams("#a,#b", "foo,bar")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdKick,
						cmdKick,
					},
				},
				"foo": mockConnection{
					messages: []message{
						cmdKick,
						cmdKick,
					},
				},
				"bar": mockConnection{
					messages: []message{
						cmdKick,
						cmdKick,
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
			in:   []message{cmdKick.withParams("#a", "foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorChanOPrivsNeeded},
				},
				"foo": mockConnection{
					messages: []message{},
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
			in:   []message{cmdKick.withParams("#a,#b,#c", "foo,foo,foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						cmdKick,
						errorChanOPrivsNeeded,
						cmdKick,
					},
				},
				"foo": mockConnection{
					messages: []message{
						cmdKick,
						cmdKick,
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
			in:   []message{cmdKick.withParams("#a,#b,#c", "foo,bar")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorNeedMoreParams,
					},
				},
				"foo": mockConnection{
					messages: []message{},
				},
				"bar": mockConnection{
					messages: []message{},
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
			in:   []message{cmdKick.withParams("#a", "foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorUserNotInChannel,
					},
				},
				"foo": mockConnection{
					messages: []message{},
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
			in:   []message{cmdKick.withParams("#b", "foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						errorNoSuchChannel,
					},
				},
				"foo": mockConnection{
					messages: []message{},
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
			in:   []message{cmdKick.withParams("#a", "foo")},
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
				withChannel("#a", "", "").
				withUser("nick", "#a").
				withUser("foo", "#a"),
			assert: []assert{
				assertUserOnChannel("foo", "#a"),
			},
		},
		{
			desc: "kick while not on channel",
			in:   []message{cmdKick.withParams("#a", "foo")},
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
				withChannel("#a", "", "").
				withUser("nick").
				withUser("foo", "#a"),
			assert: []assert{
				assertUserOnChannel("foo", "#a"),
			},
		},
	})
}
