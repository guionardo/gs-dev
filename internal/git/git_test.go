package git

import (
	"os"
	"path"
	"testing"
)

func TestGetRemoteHttpURL(t *testing.T) {
	thisFolderWithRepo, _ := os.Getwd()
	thisFolderWithRepo = path.Join(thisFolderWithRepo, "..", "..")
	t.Logf("thisFolderWithRepo: %s", thisFolderWithRepo)
	anotherFolderWithoutRepo := t.TempDir()
	t.Logf("anotherFolderWithoutRepo: %s", anotherFolderWithoutRepo)

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
			got, err := GetRemoteHttpURL(tt.folder)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemoteHttpURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRemoteHttpURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
