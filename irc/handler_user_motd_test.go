package irc

import (
	"testing"
)

func TestUserHandlerMotd(t *testing.T) {
	state := make(chan state, 1)
	handler := func() handler { return newUserHandler(state, "nick") }
	testHandler(t, "userHandler-MOTD", state, handler, []handlerTest{
		{
			desc: "no motd set",
			in:   []message{cmdMotd},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyMOTDStart.
							withPrefix("name").
							withParams("nick").
							withTrailing("- name Message of the day - "),
						replyEndOfMOTD.
							withPrefix("name").
							withParams("nick").
							withTrailing(motdFooter),
					},
				},
			},
			strict: true,
			state:  newMockState().withUser("nick"),
		},
		{
			desc: "with motd set",
			in:   []message{cmdMotd},
			wantNicks: map[string]mockConnection{
				"nick": mockConnection{
					messages: []message{
						replyMOTDStart.
							withPrefix("name").
							withParams("nick").
							withTrailing("- name Message of the day - "),
						replyMOTD.withPrefix("name").withParams("nick").withTrailing("- a"),
						replyMOTD.withPrefix("name").withParams("nick").withTrailing("- b"),
						replyMOTD.withPrefix("name").withParams("nick").withTrailing("- c"),
						replyEndOfMOTD.
							withPrefix("name").
							withParams("nick").
							withTrailing(motdFooter),
					},
				},
			},
			strict: true,
			state:  newMockState().withUser("nick"),
			motd:   []string{"a", "b", "c"},
		},
	})
}
