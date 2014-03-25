package irc

import (
	"testing"
)

func TestUserHandlerPart(t *testing.T) {
	state := make(chan State, 1)
	testHandler(t, "UserHandler-PART", state, NewUserHandler(state, "nick"), []handlerTest{
		{
			desc: "succesful part channel",
			in:   []Message{CmdPart.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{CmdPart},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel"),
		},
		{
			desc: "successful part multiple channels",
			in:   []Message{CmdPart.WithParams("#foo,#bar")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{},
				},
				"foo": mockConnection{
					messages: []Message{CmdPart},
				},
				"bar": mockConnection{
					messages: []Message{CmdPart},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "").
				withChannel("#bar", "", "").
				withUser("nick", "#foo", "#bar").
				withUser("foo", "#foo").
				withUser("bar", "#bar"),
		},
		{
			desc: "failure - not on channel",
			in:   []Message{CmdPart.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNotOnChannel},
				},
				"foo": mockConnection{
					messages: []Message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick").
				withUser("foo", "#channel"),
		},
		{
			desc: "failure - part from bad channel",
			in:   []Message{CmdPart.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNoSuchChannel},
				},
			},
			state: newMockState().withUser("nick"),
		},
	})
}
