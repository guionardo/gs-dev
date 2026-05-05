package url

import (
	"context"
	"fmt"
	"io"

	urlservice "github.com/guionardo/gs-dev/internal/services/url"
)

type UrlStruct struct {
	service *urlservice.URLService

	Directory string `flag:"directory,d" description:"repository root to get the URL from\n(can be a relative, absolute and ~ resolved paths)" default:"."`
	JustShow  bool   `flag:"just-show,j" description:"just-show the URL without opening it"`
}

func (us *UrlStruct) Run(ctx context.Context, w io.Writer) error {
	if us.JustShow {
		url, err := us.service.GetRemoteHttpURL(us.Directory)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(w, url)

		return err
	}

	return us.service.OpenRemoteURL(us.Directory)
}

func (us *UrlStruct) Setup(ctx context.Context) error {
	us.service = urlservice.NewURLService()
	return nil
}
