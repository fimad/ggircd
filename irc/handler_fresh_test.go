package irc

import (
  "testing"
)

func TestFreshHandlerHandle(t *testing.T) {
  testHandler(t, "FreshHandler", NewFreshHandler(), []handlerTest{
    {
      in:    []Message{},
      want:  MockConnection{Killed: true},
      state: NewMockState(),
    },
  })
}
