package config

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewConfigFile(t *testing.T) {
	t.Parallel()

	t.Run("NewConfigFile_with_existing_file_should_return_config", func(t *testing.T) {
		t.Parallel()

		filename := t.TempDir() + "/config.yaml"
		_ = os.WriteFile(filename, []byte("test: value"), 0600)

		config, err := NewConfigFile(filename)
		require.NoError(t, err)
		require.Equal(t, "value", config.data["test"])
	})

	t.Run("NewConfigFile_with_non_existing_file_should_return_config_with_empty_data", func(t *testing.T) {
		t.Parallel()

		config, err := NewConfigFile("non_existing_file.yaml")
		require.Error(t, err)
		require.ErrorIs(t, err, fs.ErrNotExist)
		require.NotNil(t, config.data)
	})

	t.Run("NewConfigFile_with_invalid_file_should_return_error", func(t *testing.T) {
		t.Parallel()

		filename := t.TempDir() + "/config.yaml"
		_ = os.WriteFile(filename, []byte("invalid- value"), 0600)

		config, err := NewConfigFile(filename)
		require.Error(t, err)
		require.NotNil(t, config.data)
	})
}

func TestGetSetValue(t *testing.T) {
	t.Parallel()

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

	cfg, err := NewConfigFile(t.TempDir() + "/config.yaml")
	require.Error(t, err)

	now := time.Now().UTC().Round(time.Second)

	value1 := ts1{Data1: "value1", Data2: 1}
	value2 := ts2{Data3: true, Data4: 1.0, Data5: []string{"value1", "value2"}, Data6: map[string]string{"key1": "value1", "key2": "value2"}}
	value3 := ts3{Data7: now, Data8: time.Second * 10}

	cfg.SetValue("data1", value1)
	cfg.SetValue("data2", value2)
	cfg.SetValue("data3", value3)

	v1, err := GetValue[ts1](cfg, "data1")
	require.NoError(t, err)
	require.Equal(t, ts1{Data1: "value1", Data2: 1}, v1)

	v2, err := GetValue[ts2](cfg, "data2")
	require.NoError(t, err)
	require.Equal(t, ts2{Data3: true, Data4: 1.0, Data5: []string{"value1", "value2"}, Data6: map[string]string{"key1": "value1", "key2": "value2"}}, v2)

	v3, err := GetValue[ts3](cfg, "data3")
	require.NoError(t, err)
	require.Equal(t, ts3{Data7: now, Data8: time.Second * 10}, v3)

	require.NoError(t, cfg.Save())

	cfg2, err := NewConfigFile(cfg.fileName)
	require.NoError(t, err)

	v1, err = GetValue[ts1](cfg2, "data1")
	require.NoError(t, err)
	require.Equal(t, value1, v1)

	v2, err = GetValue[ts2](cfg2, "data2")
	require.NoError(t, err)
	require.Equal(t, value2, v2)

	v3, err = GetValue[ts3](cfg2, "data3")
	require.NoError(t, err)
	require.Equal(t, value3, v3)
}
