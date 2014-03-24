package irc

import (
  "testing"
)

func TestUserPrefix(t *testing.T) {
  tests := []struct {
    in   User
    want string
  }{
    {
      in:   User{Nick: "Nick", User: "User", Host: "Host"},
      want: "Nick!User@Host",
    },
    {
      in:   User{Nick: "nick"},
      want: "nick!@",
    },
    {
      in:   User{},
      want: "!@",
    },
  }

  for i, tt := range tests {
    got := tt.in.Prefix()
    if got != tt.want {
      t.Errorf("%d. %+v.Prefix() => got %+v, want %+v", i, tt.in, got, tt.want)
    }
  }
}
