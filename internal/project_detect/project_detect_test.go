package projectdetect

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProjectType(t *testing.T) {
	t.Parallel()

	t.Run("getProjectType_should_prioritize_go_detector_order", func(t *testing.T) {
		t.Parallel()

		folder := path.Join(t.TempDir(), "mixed-project")
		require.NoError(t, os.MkdirAll(folder, 0700))
		require.NoError(t, os.WriteFile(path.Join(folder, "go.mod"), []byte("module github.com/acme/priority\n"), 0600))
		require.NoError(t, os.WriteFile(path.Join(folder, "pyproject.toml"), []byte("[project]\nname = \"ignored-py\"\n"), 0600))

		projectType, projectName := getProjectType(folder)
		require.Equal(t, "go", projectType)
		require.Equal(t, "github.com/acme/priority", projectName)
	})
}

func Test_getProjectType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		folder              string
		filename            string
		fileContent         string
		expectedProjectType string
		expectedProjectName string
	}{
		{
			name:                "go project",
			folder:              "go-project",
			filename:            "go.mod",
			fileContent:         "module github.com/acme/test",
			expectedProjectType: "go",
			expectedProjectName: "github.com/acme/test",
		},
		{
			name:                "python project",
			folder:              "python-project",
			filename:            "pyproject.toml",
			fileContent:         "[project]\nname = \"py-test\"\n",
			expectedProjectType: "python",
			expectedProjectName: "py-test",
		},
		{
			name:                "python project with requirements.txt",
			folder:              "python-project",
			filename:            "requirements.txt",
			fileContent:         "py-test",
			expectedProjectType: "python",
			expectedProjectName: "python-project",
		}, {
			name:                "rust project",
			folder:              "rust-project",
			filename:            "Cargo.toml",
			fileContent:         "[package]\nname = \"rust-test\"\n",
			expectedProjectType: "rust",
			expectedProjectName: "rust-test",
		}, {
			name:                "java project",
			folder:              "java-project",
			filename:            "pom.xml",
			fileContent:         "<project><name>java-test</name></project>",
			expectedProjectType: "java",
			expectedProjectName: "java-test",
		}, {
			name:                "java project with build.gradle",
			folder:              "java-project",
			filename:            "build.gradle",
			fileContent:         "project {\nname = \"java-test\"\n}",
			expectedProjectType: "java",
			expectedProjectName: "java-project",
		}, {
			name:                "js project",
			folder:              "js-project",
			filename:            "package.json",
			fileContent:         "{\"name\":\"js-test\"}",
			expectedProjectType: "js",
			expectedProjectName: "js-test",
		},
		{
			name:                "unknown project",
			folder:              "unknown-project",
			filename:            "unknown.txt",
			fileContent:         "unknown",
			expectedProjectType: "unknown",
			expectedProjectName: "unknown-project",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			folder := path.Join(t.TempDir(), test.folder)
			require.NoError(t, os.MkdirAll(folder, 0700))
			require.NoError(t, os.WriteFile(path.Join(folder, test.filename), []byte(test.fileContent), 0600))
			projectType, projectName := getProjectType(folder)
			require.Equal(t, test.expectedProjectType, projectType)
			require.Equal(t, test.expectedProjectName, projectName)
		})
	}
}
