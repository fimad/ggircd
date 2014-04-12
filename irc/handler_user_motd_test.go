package irc

import (
	"testing"
)

func TestUserHandlerMotd(t *testing.T) {
	state := make(chan State, 1)
	handler := func() Handler { return NewUserHandler(state, "nick") }
	testHandler(t, "UserHandler-MOTD", state, handler, []handlerTest{
		{
			desc: "no motd set",
			in:   []Message{CmdMotd},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyMOTDStart.
							WithPrefix("name").
							WithParams("nick").
							WithTrailing("- name Message of the day - "),
						ReplyEndOfMOTD.
							WithPrefix("name").
							WithParams("nick").
							WithTrailing(motdFooter),
					},
				},
			},
			strict: true,
			state:  newMockState().withUser("nick"),
		},
		{
			desc: "with motd set",
			in:   []Message{CmdMotd},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []Message{
						ReplyMOTDStart.
							WithPrefix("name").
							WithParams("nick").
							WithTrailing("- name Message of the day - "),
						ReplyMOTD.WithPrefix("name").WithParams("nick").WithTrailing("- a"),
						ReplyMOTD.WithPrefix("name").WithParams("nick").WithTrailing("- b"),
						ReplyMOTD.WithPrefix("name").WithParams("nick").WithTrailing("- c"),
						ReplyEndOfMOTD.
							WithPrefix("name").
							WithParams("nick").
							WithTrailing(motdFooter),
					},
				},
			},
			strict: true,
			state:  newMockState().withUser("nick"),
			motd:   []string{"a", "b", "c"},
		},
	})
}
