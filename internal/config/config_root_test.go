package config

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type (
	ts1 struct {
		Data1 string `yaml:"data1"`
		Data2 int    `yaml:"data2"`
	}
	ts2 struct {
		Data3 bool              `yaml:"data3"`
		Data4 float64           `yaml:"data4"`
		Data5 []string          `yaml:"data5"`
		Data6 map[string]string `yaml:"data6"`
	}
	ts3 struct {
		Data7 time.Time     `yaml:"data7"`
		Data8 time.Duration `yaml:"data8"`
	}
)

func (ts1) Key() string {
	return "ts1"
}

func (ts1) Validate() error {
	return nil
}

func (ts2) Key() string {
	return "ts2"
}

func (ts2) Validate() error {
	return nil
}
func (ts3) Key() string {
	return "ts3"
}

func (ts3) Validate() error {
	return nil
}
func TestNewConfigFile(t *testing.T) {
	t.Parallel()

	t.Run("NewConfigFile_with_existing_file_should_return_config", func(t *testing.T) {
		t.Parallel()

		configDir := t.TempDir()
		_ = os.WriteFile(path.Join(configDir, "ts1.yml"), []byte("data1: value"), 0600)

		config, err := NewConfigRoot(configDir)
		require.NoError(t, err)
		value, err := GetValue[ts1](config)
		require.NoError(t, err)
		require.Equal(t, "value", value.Data1)
	})

	t.Run("NewConfigFile_with_non_existing_file_should_return_config_with_empty_data", func(t *testing.T) {
		t.Parallel()

		_, err := NewConfigRoot("non_existing_file")
		require.Error(t, err)
	})

	t.Run("NewConfigFile_with_invalid_file_should_return_error", func(t *testing.T) {
		t.Parallel()

		configDir := t.TempDir()
		_ = os.WriteFile(path.Join(configDir, "ts1.yml"), []byte("invalid- value"), 0600)

		config, err := NewConfigRoot(configDir)
		require.NoError(t, err)

		_, err = GetValue[ts1](config)
		require.Error(t, err)
	})
}

func TestGetSetValue(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfigRoot(t.TempDir())
	require.NoError(t, err)

	now := time.Now().UTC().Round(time.Second)

	value1 := ts1{Data1: "value1", Data2: 1}
	value2 := ts2{Data3: true, Data4: 1.0, Data5: []string{"value1", "value2"}, Data6: map[string]string{"key1": "value1", "key2": "value2"}}
	value3 := ts3{Data7: now, Data8: time.Second * 10}

	SetValue(cfg, value1)
	SetValue(cfg, value2)
	SetValue(cfg, value3)

	v1, err := GetValue[ts1](cfg)
	require.NoError(t, err)
	require.Equal(t, ts1{Data1: "value1", Data2: 1}, v1)

	v2, err := GetValue[ts2](cfg)
	require.NoError(t, err)
	require.Equal(t, ts2{Data3: true, Data4: 1.0, Data5: []string{"value1", "value2"}, Data6: map[string]string{"key1": "value1", "key2": "value2"}}, v2)

	v3, err := GetValue[ts3](cfg)
	require.NoError(t, err)
	require.Equal(t, ts3{Data7: now, Data8: time.Second * 10}, v3)

	require.NoError(t, cfg.Save())

	cfg2, err := NewConfigRoot(cfg.configDir)
	require.NoError(t, err)

	v1, err = GetValue[ts1](cfg2)
	require.NoError(t, err)
	require.Equal(t, value1, v1)

	v2, err = GetValue[ts2](cfg2)
	require.NoError(t, err)
	require.Equal(t, value2, v2)

	v3, err = GetValue[ts3](cfg2)
	require.NoError(t, err)
	require.Equal(t, value3, v3)
}
