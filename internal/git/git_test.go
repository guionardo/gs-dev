package git

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRemoteHttpURL(t *testing.T) {
	t.Parallel()

	thisFolderWithRepo, _ := os.Getwd()
	thisFolderWithRepo = path.Join(thisFolderWithRepo, "..", "..")
	t.Logf("thisFolderWithRepo: %s", thisFolderWithRepo)
	anotherFolderWithoutRepo := t.TempDir()
	t.Logf("anotherFolderWithoutRepo: %s", anotherFolderWithoutRepo)

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "GITHUB") {
			t.Logf("Skipping in github action job")
			return
		}
	}

	tests := []struct {
		name    string
		folder  string
		want    string
		wantErr bool
	}{
		{"Repo GIT", thisFolderWithRepo, "https://github.com/guionardo/gs-dev", false},
		{"Repo not GIT", anotherFolderWithoutRepo, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetRemoteHttpURL(tt.folder)
			if tt.wantErr {
				require.Error(t, err, "GetRemoteHttpURL() should return error")
			}

			require.Equal(t, tt.want, got, "GetRemoteHttpURL() should return the correct URL")
		})
	}
}
