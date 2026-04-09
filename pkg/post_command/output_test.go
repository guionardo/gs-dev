package postcommand

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutput(t *testing.T) {
	t.Parallel()

	t.Run("AddOutputLine_with_lines_should_add_lines_to_output_content", func(t *testing.T) { //nolint:paralleltest
		postCommandOutputFile = t.TempDir() + "/output.txt"

		AddOutputLine("line 1")
		AddOutputLine("line 2")

		err := WriteOutput()
		require.NoError(t, err)

		content, err := os.ReadFile(postCommandOutputFile)

		require.NoError(t, err)
		assert.Equal(t, "line 1\nline 2\n", string(content))
	})

	t.Run("WriteOutput_without_output_file_should_return_nil", func(t *testing.T) { //nolint:paralleltest
		postCommandOutputFile = ""

		err := WriteOutput()

		assert.NoError(t, err)
	})
}
