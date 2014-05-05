package irc

import (
	"testing"
)

func TestUserHandlerMode(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-MODE", state, handler, []handlerTest{

		// General mode tests...

		{
			desc: "failure - no target given",
			in:   []Message{CmdMode},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNeedMoreParams},
				},
			},
			state: newMockState().withUser("nick"),
		},

		// Channel mode specific tests...

		{
			desc: "successful query channel mode",
			in:   []Message{CmdMode.WithParams("#foo")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ReplyChannelModeIs},
				},
			},
			state: newMockState().
				withChannel("#foo", "nmit", "").
				withUser("nick", "#foo"),
		},
		{
			desc: "successful parse mode positive",
			in:   []Message{CmdMode.WithParams("#foo", "+vno", "nick", "nick")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{CmdMode},
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
			in:   []Message{CmdMode.WithParams("#foo", "-stn")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{CmdMode},
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
			desc: "failure - non-existent channel",
			in:   []Message{CmdMode.WithParams("#foo", "+n")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNoSuchChannel},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - non-op attempting to change mode",
			in:   []Message{CmdMode.WithParams("#foo", "+n")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorChanOPrivsNeeded},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick"),
		},
		{
			desc: "failure - attempt to set invalid mode",
			in:   []Message{CmdMode.WithParams("#foo", "+w")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorUnknownMode},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for op",
			in:   []Message{CmdMode.WithParams("#foo", "+o")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for voice",
			in:   []Message{CmdMode.WithParams("#foo", "+v")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for key",
			in:   []Message{CmdMode.WithParams("#foo", "+k")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for limit",
			in:   []Message{CmdMode.WithParams("#foo", "+l")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - not enough parameters for ban",
			in:   []Message{CmdMode.WithParams("#foo", "+b")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withUser("nick", "#foo").
				withOps("#foo", "nick"),
		},
		{
			desc: "failure - invalid number for limit",
			in:   []Message{CmdMode.WithParams("#foo", "+l", "nan")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorUnknownMode},
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
			in:   []Message{CmdMode.WithParams("nick", "+i")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ReplyUModeIs},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", "i"),
			},
		},
		{
			desc: "success - make invisible",
			in: []Message{
				CmdMode.WithParams("nick", "+i"),
				CmdMode.WithParams("nick", "-i"),
			},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyUModeIs,
						ReplyUModeIs,
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
			in:   []Message{CmdMode.WithParams("nick")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNeedMoreParams},
				},
			},
			state: newMockState().withUser("nick"),
		},
		{
			desc: "failure - unknown mode",
			in:   []Message{CmdMode.WithParams("nick", "+v")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorUModeUnknownFlag},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", ""),
			},
		},
		{
			desc: "failure - nick mismatch",
			in:   []Message{CmdMode.WithParams("foo", "+i")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorUsersDontMatch},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", ""),
			},
		},
		{
			desc: "failure - cannot set away with mode",
			in:   []Message{CmdMode.WithParams("nick", "+a")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorUModeUnknownFlag},
				},
			},
			state: newMockState().withUser("nick"),
			assert: []assert{
				assertUserMode("nick", ""),
			},
		},
	})
}
