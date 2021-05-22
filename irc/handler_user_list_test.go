package irc

import (
	"testing"
)

func TestUserHandlerList(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-LIST", state, handler, []handlerTest{
		{
			desc: "list all w/ no channels",
			in:   []message{cmdList},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyListEnd,
					},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "list all w/ channels",
			in:   []message{cmdList},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyList,
						replyList,
						replyListEnd,
					},
				},
			},
			state: newMockState().
				withChannel("#foo", "", "topic").
				withChannel("#bar", "", "").
				withUser("nick").
				withUser("nick1", "#foo").
				withUser("nick2", "#foo").
				withUser("nick3", "#bar"),
		},
		{
			desc: "list subset of channels",
			in:   []message{cmdList.withParams("#foo,#baz")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyList.
							withPrefix("name").
							withParams("nick", "#foo", "2").
							withTrailing("topic"),
						replyList.
							withPrefix("name").
							withParams("nick", "#baz", "1"),
						replyListEnd.
							withPrefix("name").
							withParams("nick").
							withTrailing(endOfListMessage),
					},
				},
			},
			strict: true,
			state: newMockState().
				withChannel("#foo", "", "topic").
				withChannel("#bar", "", "").
				withChannel("#baz", "", "").
				withUser("nick").
				withUser("nick1", "#foo").
				withUser("nick2", "#foo").
				withUser("nick3", "#bar").
				withUser("nick4", "#baz"),
		},
		{
			desc: "list of invalid of channels",
			in:   []message{cmdList.withParams("#foo,#baz")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyListEnd,
					},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "secret/private channels are hidden from outside",
			in:   []message{cmdList},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyListEnd,
					},
				},
			},
			state: newMockState().
				withChannel("#foo", "p", "topic").
				withChannel("#bar", "s", "").
				withUser("nick").
				withUser("nick1", "#foo").
				withUser("nick2", "#foo").
				withUser("nick3", "#bar"),
		},
		{
			desc: "secret/private channels are hidden from outside by name",
			in:   []message{cmdList.withParams("#foo,#baz")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyListEnd,
					},
				},
			},
			state: newMockState().
				withChannel("#foo", "p", "topic").
				withChannel("#bar", "s", "").
				withUser("nick").
				withUser("nick1", "#foo").
				withUser("nick2", "#foo").
				withUser("nick3", "#bar"),
		},
		{
			desc: "secret channels are visible from inside",
			in:   []message{cmdList},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyList.
							withPrefix("name").
							withParams("nick", "#foo", "3").
							withTrailing("topic"),
						replyListEnd,
					},
				},
			},
			state: newMockState().
				withChannel("#foo", "s", "topic").
				withChannel("#bar", "s", "").
				withUser("nick", "#foo").
				withUser("nick1", "#foo").
				withUser("nick2", "#foo").
				withUser("nick3", "#bar"),
		},
		{
			desc: "private channels are visible from inside",
			in:   []message{cmdList},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyList.
							withPrefix("name").
							withParams("nick", "#foo", "3").
							withTrailing("topic"),
						replyListEnd,
					},
				},
			},
			state: newMockState().
				withChannel("#foo", "p", "topic").
				withChannel("#bar", "p", "").
				withUser("nick", "#foo").
				withUser("nick1", "#foo").
				withUser("nick2", "#foo").
				withUser("nick3", "#bar"),
		},
	})
}
