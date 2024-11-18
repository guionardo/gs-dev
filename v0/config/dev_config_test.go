package config

// import (
// 	"path"
// 	"reflect"
// 	"testing"

// 	"github.com/guionardo/gs-dev/config"
// )

// func TestNewDevConfig(t *testing.T) {

// 	t.Run("Default", func(t *testing.T) {
// 		tmp := t.TempDir()
// 		err := ValidateRepositoryFolder(path.Join(tmp, "gs_dev"))
// 		if err != nil {
// 			t.Errorf("GetConfigRepositoryFolder error - %v", err)
// 			return
// 		}
// 		configFolder := GetConfigRepositoryFolder()
// 		devCfg0 := NewDevConfig(configFolder)
// 		devCfg0.Folders = map[string]*config.DevFolder{}

// 		if err := devCfg0.Save(); err != nil {
// 			t.Errorf("NewDevConfig save error - %v", err)
// 			return
// 		}

// 		devCfg1 := NewDevConfig(configFolder)
// 		if err := devCfg1.Load(); err != nil {
// 			t.Errorf("NewDevConfig load error - %v", err)
// 			return
// 		}

// 		if !reflect.DeepEqual(devCfg0, devCfg1) {
// 			t.Errorf("Expected %v - got %v", devCfg0, devCfg1)
// 		}

// 	})

// }
