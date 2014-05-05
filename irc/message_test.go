package irc

import (
	"testing"
)

func TestMessageToString(t *testing.T) {
	tests := []struct {
		msg  message
		want string
		ok   bool
	}{
		{
			msg:  message{},
			want: "",
			ok:   false,
		},
		{
			msg: message{
				command: "cmd",
			},
			want: "cmd\r\n",
			ok:   true,
		},
		{
			msg: message{
				prefix:  "prefix",
				command: "cmd",
			},
			want: ":prefix cmd\r\n",
			ok:   true,
		},
		{
			msg: message{
				prefix:  "prefix",
				command: "cmd",
				params:  []string{"1", "2", "3"},
			},
			want: ":prefix cmd 1 2 3\r\n",
			ok:   true,
		},
		{
			msg: message{
				prefix:   "prefix",
				command:  "cmd",
				params:   []string{"1", "2", "3"},
				trailing: "hello world",
			},
			want: ":prefix cmd 1 2 3 :hello world\r\n",
			ok:   true,
		},
		{
			msg: message{
				prefix:   "prefix",
				command:  "cmd",
				params:   []string{"1", "2", "3"},
				trailing: "hello",
			},
			want: ":prefix cmd 1 2 3 :hello\r\n",
			ok:   true,
		},
		{
			msg: message{
				command: "cmd",
				params:  []string{"hello world"},
			},
			want: "",
			ok:   false,
		},
	}

	for i, tt := range tests {
		got, ok := tt.msg.toString()
		if tt.ok != ok || got != tt.want {
			t.Errorf("%d.\nmessage.toString(%+v) =>\n\tgot %q, %v\n\twant %q, %v",
				i, tt.msg, got, ok, tt.want, tt.ok)
		}
	}
}
