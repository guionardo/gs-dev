package environment

import (
	"os"
	"path"
	"testing"
)

func Test_getProfileFile(t *testing.T) {
	tests := []struct {
		name            string
		wantProfileFile string
	}{
		{
			name:            "bash",
			wantProfileFile: path.Join(os.Getenv("HOME"), ".bashrc"),
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotProfileFile := getProfileFile(); gotProfileFile != tt.wantProfileFile {
				t.Errorf("getProfileFile() = %v, want %v", gotProfileFile, tt.wantProfileFile)
			}
		})
	}
}
