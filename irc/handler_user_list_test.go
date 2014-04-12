package irc

import (
	"testing"
)

func TestUserHandlerList(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-LIST", state, handler, []handlerTest{
		{
			desc: "list all w/ no channels",
			in:   []Message{CmdList},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyListEnd,
					},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
		{
			desc: "list all w/ channels",
			in:   []Message{CmdList},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyList,
						ReplyList,
						ReplyListEnd,
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
			in:   []Message{CmdList.WithParams("#foo,#baz")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyList.
							WithPrefix("name").
							WithParams("nick", "#foo", "2").
							WithTrailing("topic"),
						ReplyList.
							WithPrefix("name").
							WithParams("nick", "#baz", "1"),
						ReplyListEnd.
							WithPrefix("name").
							WithParams("nick").
							WithTrailing(endOfListMessage),
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
			in:   []Message{CmdList.WithParams("#foo,#baz")},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyListEnd,
					},
				},
			},
			state: newMockState().
				withUser("nick"),
		},
	})
}
