package irc

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMessageParser(t *testing.T) {
	tests := []struct {
		raw  string
		want []message
	}{
		{
			raw:  "",
			want: []message{},
		},
		{
			raw: "foobar\r\n",
			want: []message{
				{
					command: "foobar",
					params:  []string{},
				},
			},
		},
		{
			raw: ":test foobar\r\n",
			want: []message{
				{
					prefix:  "test",
					command: "foobar",
					params:  []string{},
				},
			},
		},
		{
			raw: ":test foobar test\r\n",
			want: []message{
				{
					prefix:  "test",
					command: "foobar",
					params:  []string{"test"},
				},
			},
		},
		{
			raw: ":test foobar 1 2 3 4\r\n",
			want: []message{
				{
					prefix:  "test",
					command: "foobar",
					params:  []string{"1", "2", "3", "4"},
				},
			},
		},
		{
			raw: ":test foobar 1 2 3 4 :hello world\r\n",
			want: []message{
				{
					prefix:   "test",
					command:  "foobar",
					params:   []string{"1", "2", "3", "4"},
					trailing: "hello world",
				},
			},
		},
		{
			raw: "a\r\nb\r\nc\r\n",
			want: []message{
				{command: "a", params: []string{}},
				{command: "b", params: []string{}},
				{command: "c", params: []string{}},
			},
		},
		{
			raw: "a\nb\rc\n\rd\r\n",
			want: []message{
				{command: "a", params: []string{}},
				{command: "b", params: []string{}},
				{command: "c", params: []string{}},
				{command: "d", params: []string{}},
			},
		},
		{
			raw: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\n\r:foo bar 1 2 3\r\n",
			want: []message{
				{prefix: "foo", command: "bar", params: []string{"1", "2", "3"}},
			},
		},
		{
			raw: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\r\n",
			want: []message{
				{command: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", params: []string{}},
			},
		},
	}

	for i, tt := range tests {
		reader := bytes.NewReader([]byte(tt.raw))
		parser := newMessageParser(reader)
		got := make([]message, 0)
		for msg, hasMore := parser(); hasMore; msg, hasMore = parser() {
			got = append(got, msg)
		}

		if (len(tt.want) > 0 || len(got) > 0) && !reflect.DeepEqual(tt.want, got) {
			t.Errorf("%d.\nmessageParser(%q) =>\n\tgot %+v\n\twant %+v",
				i, tt.raw, got, tt.want)
		}
	}
}
