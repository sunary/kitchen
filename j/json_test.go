package j

import (
	"testing"
)

func TestMarshalJson(t *testing.T) {
	type args struct {
		body map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				body: map[string]string{
					"a": "1",
					"b": "\"2\"",
					"c": "3",
				},
			},
			want: "{\"a\":\"1\",\"b\":\"\\\"2\\\"\",\"c\":\"3\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarshalJson(tt.args.body); got != tt.want {
				t.Errorf("MarshalJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalJson(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			args: args{
				message: "{\"a\":\"1\",\"b\":\"\\\"2\\\"\",\"c\":\"3\"}",
			},
			want: map[string]string{
				"a": "1",
				"b": "\\\"2\\\"",
				"c": "3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnmarshalJson(tt.args.message); !compareJson(got, tt.want) {
				t.Errorf("UnmarshalJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readChar(t *testing.T) {
	type args struct {
		message string
		char    rune
		start   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				message: ", \"",
				char:    '"',
				start:   0,
			},
			want: 3,
		},
		{
			args: args{
				message: "\\:",
				char:    ':',
				start:   0,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readChar(tt.args.message, tt.args.char, tt.args.start); got != tt.want {
				t.Errorf("readChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readString(t *testing.T) {
	type args struct {
		message string
		start   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				message: "\"sunary\"s",
				start:   1,
			},
			want: "sunary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readString(tt.args.message, tt.args.start); got != tt.want {
				t.Errorf("readString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareJson(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k := range m1 {
		if m1[k] != m2[k] {
			return false
		}
	}

	return true
}
