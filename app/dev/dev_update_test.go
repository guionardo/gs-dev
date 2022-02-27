package dev

import (
	"testing"

	"github.com/guionardo/gs-dev/configs"
)

func TestUpdatePaths(t *testing.T) {
	paths := []configs.DevPathConfig{
		{FullPath: "/tmp/unexistent1"},
		{FullPath: "/tmp/unexistent2"},
	}
	devConfig := &configs.DevConfig{
		MaxSubLevels: 2,
		Paths:        paths,
	}
	cfg := &configs.RootConfig{DevConfig: *devConfig}

	UpdatePaths(cfg)

	type args struct {
		cfg *configs.RootConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := UpdatePaths(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePaths() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
