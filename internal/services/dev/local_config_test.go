package devservice_test

import (
	"path/filepath"
	"testing"

	devservice "github.com/guionardo/gs-dev/internal/services/dev"
	"github.com/stretchr/testify/require"
)

func TestNewLocalConfig(t *testing.T) {
	t.Parallel()

	projectFolder, _ := filepath.Abs("../../..")

	lc, err := devservice.NewLocalConfig(projectFolder)
	require.NoError(t, err)

	require.Equal(t, "Guiosoft Development Assistant", lc.Description)
}
