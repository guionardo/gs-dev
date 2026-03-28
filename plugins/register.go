package plugins_register

import (
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/build"
	"github.com/guionardo/gs-dev/plugins/dev"
	gitstats "github.com/guionardo/gs-dev/plugins/git_stats"
	initshell "github.com/guionardo/gs-dev/plugins/init"
	"github.com/guionardo/gs-dev/plugins/install"
	"github.com/guionardo/gs-dev/plugins/pad"
	plugins_setup "github.com/guionardo/gs-dev/plugins/plugins"
	"github.com/guionardo/gs-dev/plugins/todo"
	"github.com/guionardo/gs-dev/plugins/url"
)

func GetRegisteredPlugins() []plugins.CliPlugin {
	return []plugins.CliPlugin{
		dev.Constructor(),
		dev.FavConstructor(),
		plugins_setup.Constructor(),
		initshell.Constructor(),
		install.Constructor(),
		url.Constructor(),
		pad.Constructor(),
		todo.Constructor(),
		build.Constructor(),
		gitstats.Constructor(),
	}
}
