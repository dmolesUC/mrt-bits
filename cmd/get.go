package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/internal/quietly"
	"github.com/dmolesUC3/mrt-bits/service"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
)

const (
	usageGet     = "get <key>"
	shortDescGet = "Get a bitstream from the cloud"
	longDescGet  = shortDescGet + "\n\n" + "Gets a bitstream from the cloud and writes it to stdout."

	flagBucket  = "bucket"
	usageBucket = "bucket or container (required)"

	flagOutput = "output"
	usageOutput = "write to specified file instead of stdout"

	flagRemoteName = "remote-name"
	usageRemoteName = "write output to file named based on the remote key"

	bufsize = 512 * 1024
)

func init() {
	bucket := ""
	output := ""
	remoteName := false

	cmd := &cobra.Command{
		Use:   usageGet,
		Short: shortDescGet,
		Long:  longDescGet,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			if remoteName {
				if output != "" {
					return fmt.Errorf("\"%s\" and \"%s\" arguments cannot be specified together", flagOutput, flagRemoteName)
				}
				return downloadToFile(bucket, key, path.Base(key))
			}
			if output != "" {
				return downloadToFile(bucket, key, output)
			}
			return download(bucket, key)
		},
	}
	cmd.Flags().StringVarP(&bucket, flagBucket, "b", "", usageBucket)
	_ = cmd.MarkFlagRequired(flagBucket)

	cmd.Flags().StringVarP(&output, flagOutput, "o", "", usageOutput)
	cmd.Flags().BoolVarP(&remoteName, flagRemoteName, "O", false, usageRemoteName)

	rootCmd.AddCommand(cmd)
}

func download(bucket, key string) error {
	svc, err := flags.Service()
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stderr, "Writing %#v from service %s, bucket %#v to stdout\n", key, svc, bucket)

	downloadTo := downloader(bucket, key)
	return downloadTo(os.Stdout, svc)
}

func downloadToFile(bucket, key, filename string) error {
	svc, err := flags.Service()
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stderr, "Writing %#v from service %s, bucket %#v to %s\n", key, svc, bucket, filename)


	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		defer quietly.Close(file)
		if err != nil {
			return err
		}

		downloadTo := downloader(bucket, key)
		return downloadTo(file, svc)
	}
	return fmt.Errorf("file %#v already exists", filename)
}

func downloader(bucket, key string) (downloadTo func(out io.WriteCloser, svc service.Service) error) {
	return func(out io.WriteCloser, svc service.Service) error {
		_, body, err := svc.Get(bucket, key)
		defer quietly.Close(body)
		if err != nil {
			return err
		}
		total := 0
		defer func() {
			_, _ = fmt.Fprintf(os.Stderr, "%d bytes downloaded\n", total)
		}()
		buffer := make([]byte, bufsize)
		for {
			n, err := body.Read(buffer)
			if n > 0 {
				total += n
				_, err2 := out.Write(buffer[:n])
				if err2 != nil {
					return err2
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}
