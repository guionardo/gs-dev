package install

import (
	"context"
	"io"

	ctxdata "github.com/guionardo/gs-dev/internal/context"
	installservice "github.com/guionardo/gs-dev/internal/services/install"
)

type InstallStruct struct {
	service *installservice.InstallService

	Uninstall bool `flag:"uninstall,u" description:"uninstall bindings"`
}

func (i *InstallStruct) Run(ctx context.Context, output io.Writer) error {
	if i.Uninstall {
		return i.service.Uninstall()
	}

	return i.service.Install()
}
func (i *InstallStruct) Setup(ctx context.Context) error {
	cd, err := ctxdata.GetCommandContextData(ctx)
	if err != nil {
		return err
	}

	i.service = installservice.NewInstallService(cd.RootConfig)

	return nil
}
