package irc

import (
	"reflect"
	"testing"
)

func TestParseMode(t *testing.T) {
	tests := []struct {
		line  string
		valid mode
		want  mode
	}{
		{
			line:  "abcd",
			valid: mode{},
			want:  mode{},
		},
		{
			line: "ab",
			valid: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			want: mode{
				"a": true,
				"b": true,
			},
		},
		{
			line: "abcd",
			valid: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			want: mode{
				"a": true,
				"b": true,
				"c": true,
			},
		},
	}

	for i, tt := range tests {
		got := parseMode(tt.valid, tt.line)
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("%d.\nparseMode(%+v, %q) =>\n\tgot %+v\n\twant %+v",
				i, tt.valid, tt.line, got, tt.want)
		}
	}
}

func TestParseModeDiff(t *testing.T) {
	tests := []struct {
		line       []string
		valid      mode
		errMessage message
		posParams  mode
		negParams  mode
		wantPos    map[string][]string
		wantNeg    map[string][]string
		wantErrors []message
	}{
		{
			valid: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			errMessage: errorUnknownMode,
			posParams:  mode{},
			negParams:  mode{},
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
			valid: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			errMessage: errorUnknownMode,
			posParams:  mode{},
			negParams:  mode{},
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
			valid: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			errMessage: errorUnknownMode,
			posParams:  mode{},
			negParams:  mode{},
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
			valid: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			errMessage: errorUnknownMode,
			posParams:  mode{},
			negParams:  mode{},
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
			valid: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			errMessage: errorUnknownMode,
			posParams: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			negParams: mode{
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
			valid: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			errMessage: errorUnknownMode,
			posParams: mode{
				"a": true,
				"b": true,
				"c": true,
			},
			negParams: mode{
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
			valid: mode{
				"a": true,
			},
			errMessage: errorUnknownMode,
			posParams: mode{
				"a": true,
			},
			negParams: mode{},
			line: []string{
				"ab",
			},
			wantPos:    nil,
			wantNeg:    nil,
			wantErrors: []message{errorNeedMoreParams},
		},
		{
			valid:      mode{},
			errMessage: errorUnknownMode,
			posParams:  mode{},
			negParams:  mode{},
			line: []string{
				"ab",
			},
			wantPos: nil,
			wantNeg: nil,
			wantErrors: []message{
				errorUnknownMode.withParams("a"),
				errorUnknownMode.withParams("b"),
			},
		},
	}

	for i, tt := range tests {
		gotPos, gotNeg, gotErrors := parseModeDiff(
			tt.valid, tt.posParams, tt.negParams, tt.errMessage, tt.line)

		if !reflect.DeepEqual(tt.wantPos, gotPos) ||
			!reflect.DeepEqual(tt.wantNeg, gotNeg) ||
			!reflect.DeepEqual(tt.wantErrors, gotErrors) {
			t.Errorf("%d.\nparseModeDiff(%+v, %+v, %+v, %q) =>\n\tgot (%+v, %+v, %+v)\n\twant (%+v, %+v, %+v)",
				i, tt.valid, tt.posParams, tt.negParams, tt.line,
				gotPos, gotNeg, gotErrors,
				tt.wantPos, tt.wantNeg, tt.wantErrors)
		}
	}
}
