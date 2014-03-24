package irc

import (
	"reflect"
	"testing"
)

func TestParseMode(t *testing.T) {
	tests := []struct {
		line  string
		valid Mode
		want  Mode
	}{
		{
			line:  "abcd",
			valid: Mode{},
			want:  Mode{},
		},
		{
			line: "ab",
			valid: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			want: Mode{
				"a": true,
				"b": true,
			},
		},
		{
			line: "abcd",
			valid: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			want: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
		},
	}

	for i, tt := range tests {
		got := ParseMode(tt.valid, tt.line)
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("%d.\nParseMode(%+v, %q) =>\n\tgot %+v\b\twant %+v",
				i, tt.valid, tt.line, got, tt.want)
		}
	}
}
