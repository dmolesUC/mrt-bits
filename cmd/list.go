package cmd

import (
	"github.com/dmolesUC3/mrt-bits/operations"
	"github.com/spf13/cobra"
	"os"
)

const (
	usageList = "list <container> [<prefix>]"
	shortDescList = "list objects in the cloud"
	longDescList  = shortDescList + "\n\n" + "Lists objects in the cloud and writes keys to stdout."
)

func init() {

	extract := func(args []string) (string, string) {
		if len(args) < 1 {
			return "", ""
		}
		if len(args) < 2 {
			return args[0], ""
		}
		return args[0], args[1]
	}

	cmd := &cobra.Command{
		Use: usageList,
		Short: shortDescList,
		Long: longDescList,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			container, prefix := extract(args)
			_, err := list(container, prefix)
			return err
		},
	}
	rootCmd.AddCommand(cmd)
}

func list(container, prefix string) (int, error) {
	svc, err := flags.Service()
	if err != nil {
		return 0, nil
	}
	return operations.NewListObjects(svc, container, prefix).To(os.Stdout)
}