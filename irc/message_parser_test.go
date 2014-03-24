package irc

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMessageParser(t *testing.T) {
	tests := []struct {
		raw  string
		want []Message
	}{
		{
			raw:  "",
			want: []Message{},
		},
		{
			raw: "foobar\r\n",
			want: []Message{
				{
					Command: "foobar",
					Params:  []string{},
				},
			},
		},
		{
			raw: ":test foobar\r\n",
			want: []Message{
				{
					Prefix:  "test",
					Command: "foobar",
					Params:  []string{},
				},
			},
		},
		{
			raw: ":test foobar test\r\n",
			want: []Message{
				{
					Prefix:  "test",
					Command: "foobar",
					Params:  []string{"test"},
				},
			},
		},
		{
			raw: ":test foobar 1 2 3 4\r\n",
			want: []Message{
				{
					Prefix:  "test",
					Command: "foobar",
					Params:  []string{"1", "2", "3", "4"},
				},
			},
		},
		{
			raw: ":test foobar 1 2 3 4 :hello world\r\n",
			want: []Message{
				{
					Prefix:   "test",
					Command:  "foobar",
					Params:   []string{"1", "2", "3", "4"},
					Trailing: "hello world",
				},
			},
		},
		{
			raw: "a\r\nb\r\nc\r\n",
			want: []Message{
				{Command: "a", Params: []string{}},
				{Command: "b", Params: []string{}},
				{Command: "c", Params: []string{}},
			},
		},
	}

	for i, tt := range tests {
		reader := bytes.NewReader([]byte(tt.raw))
		parser := NewMessageParser(reader)
		got := make([]Message, 0)
		for msg, hasMore := parser(); hasMore; msg, hasMore = parser() {
			got = append(got, msg)
		}

		if (len(tt.want) > 0 || len(got) > 0) && !reflect.DeepEqual(tt.want, got) {
			t.Errorf("%d.\nMessageParser(%q) =>\n\tgot %+v\n\twant %+v",
				i, tt.raw, got, tt.want)
		}
	}
}
