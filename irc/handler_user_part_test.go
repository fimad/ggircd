package irc

import (
	"testing"
)

func TestUserHandlerPart(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-PART", state, handler, []handlerTest{
		{
			desc: "successful part channel",
			in:   []message{cmdPart.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdPart},
				},
				"foo": mockConnection{
					messages: []message{cmdPart},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel").
				withUser("foo", "#channel"),
		},
		{
			desc: "successful part multiple channels",
			in:   []message{cmdPart.withParams("#foo,#bar")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdPart, cmdPart},
				},
				"foo": mockConnection{
					messages: []message{cmdPart},
				},
				"bar": mockConnection{
					messages: []message{cmdPart},
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
			in:   []message{cmdPart.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNotOnChannel},
				},
				"foo": mockConnection{
					messages: []message{},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick").
				withUser("foo", "#channel"),
		},
		{
			desc: "failure - part from bad channel",
			in:   []message{cmdPart.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNoSuchChannel},
				},
			},
			state: newMockState().withUser("nick"),
		},
	})
}
