package irc

import (
  "reflect"
  "testing"
)

type handlerTest struct {
  desc  string
  in    []Message
  state *mockState
  want  mockConnection

  // If true, perform a fuzzy comparison for messages where only the commands
  // are compared.
  fuzzy bool

  // If true, call Handler.Close after the sequence of messages has been
  // processed.
  hangup bool
}

// testHandler takes a seque
func testHandler(t *testing.T, name string, state chan State, handler Handler, tests []handlerTest) {
  for i, tt := range tests {
    state <- tt.state
    got := runHandler(tt, handler)
    _ = <-state
    if !compareMessages(tt.fuzzy, got, tt.want) {
      t.Errorf("%d. %s: %s\n%+v =>\n\tgot %+v\n\twant %+v",
        i, name, tt.desc, tt.in, got, tt.want)
    }
  }
}

func runHandler(tt handlerTest, handler Handler) mockConnection {
  conn := mockConnection{}
  for _, message := range tt.in {
    handler = handler.Handle(&conn, message)
  }
  if tt.hangup {
    handler.Closed(&conn)
  }
  return conn
}

func compareMessages(fuzzy bool, got, want mockConnection) bool {
  if !fuzzy {
    return reflect.DeepEqual(got, want)
  }

  if got.killed != want.killed {
    return false
  }

  if len(got.messages) != len(want.messages) {
    return false
  }

  for i := 0; i < len(got.messages); i++ {
    if got.messages[i].Command != want.messages[i].Command {
      return false
    }
  }

  return true
}
