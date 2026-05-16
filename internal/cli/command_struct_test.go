package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	SampleCommandStruct struct {
		Name      string        `flag:"name,n" description:"Name of the user" validate:"required"`
		Age       int           `flag:"age,a" description:"Age of the user" validate:"required" default:"12"`
		Enabled   bool          `flag:"enabled" description:"User is enabled" default:"false"`
		TTL       time.Duration `flag:"ttl"`
		JustAtime time.Time     `flag:"time" description:"Time" default:"2026-01-01"`
		Extra     []string      `flag:"extra" default:"1,2,3"`

		RunCommand  SampleSubCommandStruct  `subcommand:"sub,sub command in the sample command"`
		RunCommand2 SampleSubcommandStruct2 `subcommand:"sub2,another sub command in the sample command"`
	}
	SampleSubCommandStruct struct {
		Args []string `args:"1,description for args"` // Count of nargs, description
	}
	SampleSubcommandStruct2 struct {
		Args []string `args:"-1,description for args"` // -1 means at least 1 argument
	}
	SampleConfig struct {
		text string //nolint: unused
	}
)

func TestGenerateCobraCommand(t *testing.T) { //nolint: funlen
	t.Parallel()

	t.Run("show_help", func(t *testing.T) {
		t.Parallel()

		var sample SampleCommandStruct

		cmd := GenerateCobraCommand(&sample, "sample_command", "Sample command for test purposes", "This is a sample command for test purposes", false)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"--help"})
		err := cmd.ExecuteContext(t.Context())

		require.NoError(t, err)

		cmd.SetArgs([]string{"sub", "--help"})
		require.NoError(t, cmd.ExecuteContext(t.Context()))

		t.Logf("Output:\n%s", buf.String())
	})

	t.Run("run_command", func(t *testing.T) {
		t.Parallel()

		var sample SampleCommandStruct

		cmd := GenerateCobraCommand(&sample, "sample_command", "Sample command for test purposes", "This is a sample command for test purposes", false)
		require.NotNil(t, cmd)

		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"--name", "Guionardo", "--enabled", "--ttl", "1h", "--age", "25", "--extra", "A", "--extra", "B"})
		require.NoError(t, cmd.ExecuteContext(t.Context()))

		assert.Equal(t, "Guionardo", sample.Name)
		assert.True(t, sample.Enabled)
		assert.Equal(t, time.Hour, sample.TTL)
		assert.Equal(t, []string{"A", "B"}, sample.Extra)

		require.Contains(t, buf.String(), "Name:Guionardo")
		require.Contains(t, buf.String(), "Enabled:true")
		require.Contains(t, buf.String(), "TTL:1h0m0s")
		require.Contains(t, buf.String(), "Extra:[A B]")
	})

	t.Run("run_subcommand", func(t *testing.T) {
		t.Parallel()

		var sample SampleCommandStruct

		cmd := GenerateCobraCommand(&sample, "sample_command", "Sample command for test purposes", "This is a sample command for test purposes", false)
		require.NotNil(t, cmd)

		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"sub", "Another arguments"})

		require.NoError(t, cmd.ExecuteContext(t.Context()))
		require.Contains(t, buf.String(), "{Args:[Another arguments]}")
	})
}

func (c SampleCommandStruct) Run(ctx context.Context, output io.Writer) error {
	_, _ = fmt.Fprintf(output, "%+v\n", c)
	return nil
}

func (c SampleCommandStruct) Setup(ctx context.Context) error {
	_, _ = fmt.Println("Running setup for SampleCommandStruct")
	return nil
}

func (c SampleSubCommandStruct) Run(ctx context.Context, output io.Writer) error {
	_, _ = fmt.Fprintf(output, "%+v\n", c)
	return nil
}

func (c SampleSubCommandStruct) Setup(ctx context.Context) error {
	_, _ = fmt.Println("Running setup for SampleSubCommandStruct")
	return nil
}
