package irc

import (
	"testing"
)

func TestUserHandlerMode(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-MODE", state, handler, []handlerTest{

		// General mode tests...

		{
			desc: "failure - no target given",
			in:   []message{cmdMode},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().withUser("nick"),
		},

		// Channel mode specific tests...

		{
			desc: "successful query channel mode",
			in:   []message{cmdMode.withParams("#foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{replyChannelModeIs},
				},
			},
			state: newMockState().
				withChannel("#foo", "nmit", "").
				withUser("nick", "#foo"),
		},
		{
			desc: "successful parse mode positive",
			in:   []message{cmdMode.withParams("#foo", "+vno", "nick", "nick")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdMode},
				},
			},
			assert: []assert{
				assertChannelMode("#foo", "n"),
				assertChannelVoice("#foo", "nick"),
				assertChannelOp("#foo", "nick"),
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "successful parse mode negative",
			in:   []message{cmdMode.withParams("#foo", "-stn")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdMode},
				},
			},
			assert: []assert{
				assertChannelMode("#foo", "m"),
			},
			state: newMockState().
				withChannel("#foo", "stnm", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "successful set channel private",
			in:   []message{cmdMode.withParams("#foo", "+p")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdMode},
				},
			},
			assert: []assert{
				assertChannelMode("#foo", "p"),
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "successful set channel secret",
			in:   []message{cmdMode.withParams("#foo", "+s")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdMode},
				},
			},
			assert: []assert{
				assertChannelMode("#foo", "s"),
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "successful swap channel secret/private",
			in: []message{
				cmdMode.withParams("#foo", "+s"),
				cmdMode.withParams("#foo", "+p-s"),
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdMode, cmdMode},
				},
			},
			assert: []assert{
				assertChannelMode("#foo", "p"),
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - channel cannot be both secret and private",
			in: []message{
				cmdMode.withParams("#foo", "+sp"),
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{},
				},
			},
			assert: []assert{
				assertChannelMode("#foo", ""),
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - channel cannot be both secret and then also private",
			in: []message{
				cmdMode.withParams("#foo", "+s"),
				cmdMode.withParams("#foo", "+p"),
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdMode},
				},
			},
			assert: []assert{
				assertChannelMode("#foo", "s"),
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - non-existent channel",
			in:   []message{cmdMode.withParams("#foo", "+n")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNoSuchChannel},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - non-op attempting to change mode",
			in:   []message{cmdMode.withParams("#foo", "+n")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorChanOPrivsNeeded},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick"),
		},
		{
			desc: "failure - attempt to set invalid mode",
			in:   []message{cmdMode.withParams("#foo", "+w")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorUnknownMode},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for op",
			in:   []message{cmdMode.withParams("#foo", "+o")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for voice",
			in:   []message{cmdMode.withParams("#foo", "+v")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for key",
			in:   []message{cmdMode.withParams("#foo", "+k")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for limit",
			in:   []message{cmdMode.withParams("#foo", "+l")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for ban",
			in:   []message{cmdMode.withParams("#foo", "+b")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - invalid number for limit",
			in:   []message{cmdMode.withParams("#foo", "+l", "nan")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorUnknownMode},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},

		// User mode specific tests...

		{
			desc: "success - make invisible",
			in:   []message{cmdMode.withParams("nick", "+i")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{replyUModeIs},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", "i"),
			},
		},
		{
			desc: "success - make invisible",
			in: []message{
				cmdMode.withParams("nick", "+i"),
				cmdMode.withParams("nick", "-i"),
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyUModeIs,
						replyUModeIs,
					},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", ""),
			},
		},
		{
			desc: "failure - need more params",
			in:   []message{cmdMode.withParams("nick")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "failure - unknown mode",
			in:   []message{cmdMode.withParams("nick", "+v")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorUModeUnknownFlag},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", ""),
			},
		},
		{
			desc: "failure - nick mismatch",
			in:   []message{cmdMode.withParams("foo", "+i")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorUsersDontMatch},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", ""),
			},
		},
		{
			desc: "failure - cannot set away with mode",
			in:   []message{cmdMode.withParams("nick", "+a")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorUModeUnknownFlag},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", ""),
			},
		},
	})
}
