package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestXxx(t *testing.T) {
	cmd := cobra.Command{}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {

	})
}
