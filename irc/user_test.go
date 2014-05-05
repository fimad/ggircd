package irc

import (
	"testing"
)

func TestUserPrefix(t *testing.T) {
	tests := []struct {
		in   user
		want string
	}{
		{
			in:   user{nick: "Nick", user: "User", host: "Host"},
			want: "Nick!User@Host",
		},
		{
			in:   user{nick: "nick"},
			want: "nick!@",
		},
		{
			in:   user{},
			want: "!@",
		},
	}

	for i, tt := range tests {
		got := tt.in.prefix()
		if got != tt.want {
			t.Errorf("%d. %+v.prefix() =>\n\tgot %+v\n\twant %+v", i, tt.in, got, tt.want)
		}
	}
}
