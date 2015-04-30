package irc

import (
	"testing"
)

func TestUserHandlerTopic(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-TOPIC", state, handler, []handlerTest{
		{
			desc: "successful query empty topic",
			in:   []message{cmdTopic.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{replyNoTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick"),
		},
		{
			desc: "successful query non-empty topic",
			in:   []message{cmdTopic.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{replyTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "topic").
				withUser("nick"),
		},
		{
			desc: "successful set topic non-op",
			in:   []message{cmdTopic.withParams("#channel").withTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel"),
		},
		{
			desc: "successful set topic w/ no trailing",
			in:   []message{cmdTopic.withParams("#channel", "topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick", "#channel"),
		},
		{
			desc: "successful set topic op",
			in:   []message{cmdTopic.withParams("#channel").withTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "t", "").
				withUser("nick", "#channel").
				withOps("#channel", "nick"),
		},
		{
			desc: "setting topic broad casts to channel",
			in:   []message{cmdTopic.withParams("#channel").withTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{cmdTopic},
				},
				"foo": mockConnection{
					messages: []message{cmdTopic},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("foo", "#channel").
				withUser("nick", "#channel"),
		},
		{
			desc: "failure - no channel given",
			in:   []message{cmdTopic},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNeedMoreParams},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick"),
		},
		{
			desc: "failure - no such channel",
			in:   []message{cmdTopic.withParams("#channel")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNoSuchChannel},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "failure - set from outside channel",
			in:   []message{cmdTopic.withParams("#channel").withTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorNotOnChannel},
				},
			},
			state: newMockState().
				withChannel("#channel", "", "").
				withUser("nick"),
		},
		{
			desc: "failure - non-op set on +t channel",
			in:   []message{cmdTopic.withParams("#channel").withTrailing("topic")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{errorChanOPrivsNeeded},
				},
			},
			state: newMockState().
				withChannel("#channel", "t", "").
				withUser("nick", "#channel"),
		},
	})
}
