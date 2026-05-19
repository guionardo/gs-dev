package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfigDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		appName string
		want    string
	}{
		{"default", "gs-dev-test", ".config/gs-dev-test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.True(t, strings.HasSuffix(GetConfigDir(tt.appName), tt.want), "GetConfigDir() should return the correct config dir")
		})
	}
}
