package irc

import (
	"testing"
)

func TestUserHandlerMode(t *testing.T) {
	state := make(chan State, 1)
	testHandler(t, "UserHandler-MODE", state, NewUserHandler(state, "nick"), []handlerTest{

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

	})
}
