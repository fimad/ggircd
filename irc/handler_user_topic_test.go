package irc

import (
	"testing"
)

func TestUserHandlerTopic(t *testing.T) {
	state := make(chan State, 1)
	testHandler(t, "UserHandler-TOPIC", state, NewUserHandler(state, "nick"), []handlerTest{
		{
			desc: "successful query empty topic",
			in:   []Message{CmdTopic.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ReplyNoTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick"),
		},
		{
			desc: "successful query non-empty topic",
			in:   []Message{CmdTopic.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ReplyTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "topic").
				withUser("nick"),
		},
		{
			desc: "successful set topic non-op",
			in:   []Message{CmdTopic.WithParams("#channel").WithTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{CmdTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel"),
		},
		{
			desc: "successful set topic op",
			in:   []Message{CmdTopic.WithParams("#channel").WithTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{CmdTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "t", "").
				withUser("nick", "#channel").
				withOps("#channel", "nick"),
		},
		{
			desc: "setting topic broad casts to channel",
			in:   []Message{CmdTopic.WithParams("#channel").WithTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{CmdTopic},
				},
				"foo": mockConnection{
					messages: []Message{CmdTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("foo", "#channel").
				withUser("nick", "#channel"),
		},
		{
			desc: "failure - no channel given",
			in:   []Message{CmdTopic},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick"),
		},
		{
			desc: "failure - no such channel",
			in:   []Message{CmdTopic.WithParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNoSuchChannel},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - set from outside channel",
			in:   []Message{CmdTopic.WithParams("#channel").WithTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorNotOnChannel},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick"),
		},
		{
			desc: "failure - non-op set on +t channel",
			in:   []Message{CmdTopic.WithParams("#channel").WithTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{ErrorChanOPrivsNeeded},
				},
			},
			state: newMockState().
				withChannel("#channel", "t", "").
				withUser("nick", "#channel"),
		},
	})
}
