package colors

import "testing"

func TestColors(t *testing.T) {
	tests := []struct {
		name string
		fnc  func(string, ...interface{})
	}{
		{"Normal", Normal},
		{"Red", Red},
		{"Green", Green},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fnc("%s: %s - %d", tt.name, "test", 1)
		})
	}
}
