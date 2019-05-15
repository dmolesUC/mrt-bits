package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
	"github.com/dmolesUC3/mrt-bits/operations"
	"github.com/dmolesUC3/mrt-bits/service"
	"github.com/spf13/cobra"
	"os"
)

const (
	usageGet     = "get <container> <key>"
	shortDescGet = "Get a bitstream from the cloud"
	longDescGet  = shortDescGet + "\n\n" + "Gets a bitstream from the cloud and writes it to stdout."

	flagOutput  = "output"
	usageOutput = "write to specified file instead of stdout"

	flagRemoteName  = "remote-name"
	usageRemoteName = "write output to file named based on the remote key"

	flagArchive  = "archive"
	usageArchive = "treat key as a prefix, and return a ZIP archive of all objects with that prefix"
)

type get struct {
	output     string
	remoteName bool
	archive    bool
}

func (g *get) get(container, key string) (int, error) {
	svc, err := flags.Service()
	if err != nil {
		return 0, err
	}
	if g.archive {
		return g.downloadArchive(svc, container, key)
	}
	return g.download(svc, container, key)
}

func (g *get) download(svc service.Service, container, key string) (int, error) {
	download := operations.NewDownload(svc, container, key)
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

func (g *get) downloadArchive(svc service.Service, container, prefix string) (int, error) {
	if g.remoteName {
		return 0, fmt.Errorf("%#v and %#v arguments cannot be specified together", flagArchive, flagRemoteName)
	}
	archive := operations.NewZipArchive(svc, container, prefix)

	// TODO: only if verbose
	size, count, err := archive.Size()
	if err != nil {
		return 0, err
	}
	quietly.Fprintf(os.Stderr, "reading %d files (%d bytes expected)", count, size)

	var out *os.File
	if g.output == "" {
		out = os.Stdout
	} else if _, err = os.Stat(g.output); os.IsNotExist(err) {
		// TODO: don't create file till we know we've got something
		//       (and/or quietly delete file)
		out, err = os.Create(g.output)
		defer quietly.Close(out)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("file %#v already exists", g.output)
	}

	return archive.To(out)
}

func (g *get) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   usageGet,
		Short: shortDescGet,
		Long:  longDescGet,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := g.get(args[0], args[1])
			return err
		},
	}
	// TODO: support DNS-based addressing
	cmd.Flags().StringVarP(&g.output, flagOutput, "o", "", usageOutput)
	cmd.Flags().BoolVarP(&g.remoteName, flagRemoteName, "O", false, usageRemoteName)
	cmd.Flags().BoolVarP(&g.archive, flagArchive, "a", false, usageArchive)
	return cmd
}

func init() {
	rootCmd.AddCommand((&get{}).command())
}
