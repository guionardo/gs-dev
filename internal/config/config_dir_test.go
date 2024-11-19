package config

import (
	"strings"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		want    string
	}{
		{"default", "gs-dev-test", ".config/gs-dev-test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigDir(tt.appName); !strings.HasSuffix(got, tt.want) {
				t.Errorf("GetConfigDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
