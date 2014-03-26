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
			t.Errorf("%d.\nParseMode(%+v, %q) =>\n\tgot %+v\n\twant %+v",
				i, tt.valid, tt.line, got, tt.want)
		}
	}
}

func TestParseModeDiff(t *testing.T) {
	tests := []struct {
		line       []string
		valid      Mode
		posParams  Mode
		negParams  Mode
		wantPos    map[string][]string
		wantNeg    map[string][]string
		wantErrors []Message
	}{
		{
			valid: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			posParams: Mode{},
			negParams: Mode{},
			line: []string{
				"abc",
			},
			wantPos: map[string][]string{
				"a": []string{""},
				"b": []string{""},
				"c": []string{""},
			},
			wantNeg:    map[string][]string{},
			wantErrors: nil,
		},
		{
			valid: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			posParams: Mode{},
			negParams: Mode{},
			line: []string{
				"-abc",
			},
			wantPos: map[string][]string{},
			wantNeg: map[string][]string{
				"a": []string{""},
				"b": []string{""},
				"c": []string{""},
			},
			wantErrors: nil,
		},
		{
			valid: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			posParams: Mode{},
			negParams: Mode{},
			line: []string{
				"-a+b-c",
			},
			wantPos: map[string][]string{
				"b": []string{""},
			},
			wantNeg: map[string][]string{
				"a": []string{""},
				"c": []string{""},
			},
			wantErrors: nil,
		},
		{
			valid: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			posParams: Mode{},
			negParams: Mode{},
			line: []string{
				"-a+b-c+c-c",
			},
			wantPos: map[string][]string{
				"b": []string{""},
			},
			wantNeg: map[string][]string{
				"a": []string{""},
				"c": []string{""},
			},
			wantErrors: nil,
		},
		{
			valid: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			posParams: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			negParams: Mode{
				"a": true,
				"c": true,
			},
			line: []string{
				"-a+b-c+b-b",
				"1",
				"2",
				"3",
				"4",
				"5",
			},
			wantPos: map[string][]string{},
			wantNeg: map[string][]string{
				"a": []string{"1"},
				"b": []string{""},
				"c": []string{"3"},
			},
			wantErrors: nil,
		},
		{
			valid: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			posParams: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			negParams: Mode{
				"a": true,
				"b": true,
				"c": true,
			},
			line: []string{
				"-a+b-c+b-b",
				"1",
				"2",
				"3",
				"4",
				"5",
			},
			wantPos: map[string][]string{
				"b": []string{"2", "4"},
			},
			wantNeg: map[string][]string{
				"a": []string{"1"},
				"b": []string{"5"},
				"c": []string{"3"},
			},
			wantErrors: nil,
		},
		{
			valid: Mode{
				"a": true,
			},
			posParams: Mode{
				"a": true,
			},
			negParams: Mode{},
			line: []string{
				"ab",
			},
			wantPos:    nil,
			wantNeg:    nil,
			wantErrors: []Message{ErrorNeedMoreParams},
		},
		{
			valid:     Mode{},
			posParams: Mode{},
			negParams: Mode{},
			line: []string{
				"ab",
			},
			wantPos: nil,
			wantNeg: nil,
			wantErrors: []Message{
				ErrorUnknownMode.WithParams("a"),
				ErrorUnknownMode.WithParams("b"),
			},
		},
	}

	for i, tt := range tests {
		gotPos, gotNeg, gotErrors := ParseModeDiff(
			tt.valid, tt.posParams, tt.negParams, tt.line)

		if !reflect.DeepEqual(tt.wantPos, gotPos) ||
			!reflect.DeepEqual(tt.wantNeg, gotNeg) ||
			!reflect.DeepEqual(tt.wantErrors, gotErrors) {
			t.Errorf("%d.\nParseModeDiff(%+v, %+v, %+v, %q) =>\n\tgot (%+v, %+v, %+v)\n\twant (%+v, %+v, %+v)",
				i, tt.valid, tt.posParams, tt.negParams, tt.line,
				gotPos, gotNeg, gotErrors,
				tt.wantPos, tt.wantNeg, tt.wantErrors)
		}
	}
}
