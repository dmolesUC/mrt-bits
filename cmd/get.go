package cmd

import (
	"fmt"
	dl "github.com/dmolesUC3/mrt-bits/internal/download"
	"github.com/spf13/cobra"
	"os"
)

const (
	usageGet     = "get <key>"
	shortDescGet = "Get a bitstream from the cloud"
	longDescGet  = shortDescGet + "\n\n" + "Gets a bitstream from the cloud and writes it to stdout."

	flagBucket  = "bucket"
	usageBucket = "bucket or container (required)"

	flagOutput  = "output"
	usageOutput = "write to specified file instead of stdout"

	flagRemoteName  = "remote-name"
	usageRemoteName = "write output to file named based on the remote key"
)

type get struct {
	output     string
	remoteName bool
	bucket     string
}

func (g *get) get(key string) (int, error) {
	svc, err := flags.Service()
	if err != nil {
		return 0, err
	}
	download := dl.NewDownload(svc, g.bucket, key)
	if g.remoteName {
		if g.output == "" {
			return download.ToRemoteFile()
		}
		return 0, fmt.Errorf("%#v and %#v arguments cannot be specified together", flagOutput, flagRemoteName)
	}
	if g.output == "" {
		return download.To(os.Stdout)
	}
	return download.ToFile(g.output)
}

func (g *get) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   usageGet,
		Short: shortDescGet,
		Long:  longDescGet,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := g.get(args[0])
			return err
		},
	}
	// TODO: support DNS-based addressing
	cmd.Flags().StringVarP(&g.bucket, flagBucket, "b", "", usageBucket)
	_ = cmd.MarkFlagRequired(flagBucket)

	cmd.Flags().StringVarP(&g.output, flagOutput, "o", "", usageOutput)
	cmd.Flags().BoolVarP(&g.remoteName, flagRemoteName, "O", false, usageRemoteName)
	return cmd
}

func init() {
	rootCmd.AddCommand((&get{}).command())
}
