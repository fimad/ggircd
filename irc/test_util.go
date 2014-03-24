package irc

import (
  "reflect"
  "testing"
)

type handlerTest struct {
  in      []Message
  state   MockState
  want    MockConnection
}

// testHandler takes a seque
func testHandler(t *testing.T, name string, handler Handler, tests []handlerTest) {
  for i, tt := range tests {
    got := runHandler(tt.in, tt.state, handler)
    if !reflect.DeepEqual(got, tt.want) {
      t.Errorf("%d. %s(%+v) => got %+v, want %+v", i, name, tt.in, got, tt.want)
    }
  }
}

func runHandler(in []Message, state MockState, handler Handler) MockConnection {
  conn := MockConnection{}
  for _, message := range in {
    handler = handler.Handle(&conn, message)
  }
  return conn
}
