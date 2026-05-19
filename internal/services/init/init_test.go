package initservice_test

import (
	"testing"

	initservice "github.com/guionardo/gs-dev/internal/services/init"
)

func TestInitService_CheckAvailableVersion(t *testing.T) {
	t.Parallel()

	initService := initservice.NewInitService()

	release, err := initService.GetAvailableVersion()
	if err != nil {
		t.Errorf("CheckAvailableVersion() failed: %v", err)
	}

	if release == nil {
		t.Errorf("GetAvailableVersion() returned nil")
	}

	t.Logf("Version: %+v", release)
}
