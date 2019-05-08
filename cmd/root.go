package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/service"
	"github.com/spf13/cobra"
)

// ------------------------------------------------------------
// Exported symbols

func Execute() error {
	return rootCmd.Execute()
}

func AddCommand(c *cobra.Command) {
	rootCmd.AddCommand(c)
}

// ------------------------------------------------------------
// Unexported symbols

const (
	usageRoot = "mrt-bits <command> [<args>]"
	shortDescRoot = "Merritt bitstream service (experimental)"
)

var longDescRoot = fmt.Sprintf(
	"%s\n\nmrt-bits supports the following environment variables:\n\n%s",
	shortDescRoot, service.AllEnvs,
)

var rootCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use: usageRoot,
		Short: shortDescRoot,
		Long:  longDescRoot,
		SilenceUsage: true,
	}
	flags.AddTo(cmd.PersistentFlags())
	return cmd
}()
