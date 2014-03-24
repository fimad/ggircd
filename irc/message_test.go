package irc

import (
	"testing"
)

func TestMessageToString(t *testing.T) {
	tests := []struct {
		msg  Message
		want string
		ok   bool
	}{
		{
			msg:  Message{},
			want: "",
			ok:   false,
		},
		{
			msg: Message{
				Command: "cmd",
			},
			want: "cmd\r\n",
			ok:   true,
		},
		{
			msg: Message{
				Prefix:  "prefix",
				Command: "cmd",
			},
			want: ":prefix cmd\r\n",
			ok:   true,
		},
		{
			msg: Message{
				Prefix:  "prefix",
				Command: "cmd",
				Params:  []string{"1", "2", "3"},
			},
			want: ":prefix cmd 1 2 3\r\n",
			ok:   true,
		},
		{
			msg: Message{
				Prefix:   "prefix",
				Command:  "cmd",
				Params:   []string{"1", "2", "3"},
				Trailing: "hello world",
			},
			want: ":prefix cmd 1 2 3 :hello world\r\n",
			ok:   true,
		},
		{
			msg: Message{
				Prefix:   "prefix",
				Command:  "cmd",
				Params:   []string{"1", "2", "3"},
				Trailing: "hello",
			},
			want: ":prefix cmd 1 2 3 :hello\r\n",
			ok:   true,
		},
		{
			msg: Message{
				Command: "cmd",
				Params:  []string{"hello world"},
			},
			want: "",
			ok:   false,
		},
	}

	for i, tt := range tests {
		got, ok := tt.msg.ToString()
		if tt.ok != ok || got != tt.want {
			t.Errorf("%d.\nMessage.ToString(%+v) =>\n\tgot %q, %v\n\twant %q, %v",
				i, tt.msg, got, ok, tt.want, tt.ok)
		}
	}
}
