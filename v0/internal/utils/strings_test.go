package utils

import "testing"

func TestPad(t *testing.T) {
	type args struct {
		text   string
		length int
		pad    PadType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Left", args{"Test", 10, Left}, "Test      "},
		{"Right", args{"Test", 10, Right}, "      Test"},
		{"Center", args{"Test", 10, Center}, "   Test   "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pad(tt.args.text, tt.args.length, tt.args.pad); got != tt.want {
				t.Errorf("Pad() = %v, want %v", got, tt.want)
			}
		})
	}
}
