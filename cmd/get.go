package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/service"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const (
	usageGet     = "get <key>"
	shortDescGet = "Get a bitstream from the cloud"
	longDescGet  = shortDescGet + "\n\n" + "Gets a bitstream from the cloud and writes it to stdout."

	flagBucket  = "bucket"
	usageBucket = "bucket or container (required)"

	bufsize = 512 * 1024
)

func init() {
	bucket := ""

	cmd := &cobra.Command{
		Use:   usageGet,
		Short: shortDescGet,
		Long:  longDescGet,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(bucket, args[0])
		},
	}
	// TODO: add -o/--output and -O/--remote-name
	cmd.Flags().StringVarP(&bucket, flagBucket, "b", "", usageBucket)
	_ = cmd.MarkFlagRequired(flagBucket)

	rootCmd.AddCommand(cmd)
}

func get(bucket, key string) error {
	svc, err := flags.Service()
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stderr, "Getting %#v from service %s, bucket %#v\n", key, svc, bucket)

	_, body, err := svc.Get(bucket, key)
	defer service.CloseQuietly(body)
	if err != nil {
		return err
	}
	buffer := make([]byte, bufsize)
	for {
		n, err := body.Read(buffer)
		if n > 0 {
			_, err2 := os.Stdout.Write(buffer[:n])
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
