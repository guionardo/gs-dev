package plugins_register

import (
	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/build"
	"github.com/guionardo/gs-dev/plugins/dev"
	initshell "github.com/guionardo/gs-dev/plugins/init"
	"github.com/guionardo/gs-dev/plugins/install"
	"github.com/guionardo/gs-dev/plugins/pad"
	plugins_setup "github.com/guionardo/gs-dev/plugins/plugins"
	"github.com/guionardo/gs-dev/plugins/todo"
	"github.com/guionardo/gs-dev/plugins/url"
)

func GetRegisteredPlugins(output *outputfile.OutputFile) []plugins.CliPlugin {
	return []plugins.CliPlugin{
		dev.Constructor(output),
		dev.FavConstructor(output),
		plugins_setup.Constructor(output),
		initshell.Constructor(output),
		install.Constructor(output),
		url.Constructor(output),
		pad.Constructor(output),
		todo.Constructor(output),
		build.Constructor(output),
	}
}
