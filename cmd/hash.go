package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/operations"
	"github.com/spf13/cobra"
)

const (
	usageHash = "hash <container> <key>"
	shortDescHash = "Get the digest of a bitstream in the cloud"
	longDescHash  = shortDescGet + "\n\n" + "Hashes a bitstream from the cloud with the specified algorithm."

	flagAlgorithm = "algorithm"
	usageAlgorithm = "hash algorithm (sha256 or md5)"
)

type hash struct {
	algorithm string
}

func (h *hash) hash(container, key string) ([]byte, error) {
	svc, err := flags.Service()
	if err != nil {
		return nil, err
	}
	info := operations.NewInfo(svc, container, key)
	switch h.algorithm {
	case "sha256":
		return info.SHA256()
	case "md5":
		return info.MD5()
	default:
		return nil, fmt.Errorf("unknown digest algorithm: %#v", h.algorithm)
	}
}

func (h *hash) printHash(container, key string) (error) {
	digest, err := h.hash(container, key)
	if err != nil {
		return err
	}
	fmt.Printf("%x\n", digest)
	return nil
}

func (h *hash) command() *cobra.Command {
	cmd := &cobra.Command {
		Use: usageHash,
		Short: shortDescHash,
		Long: longDescHash,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.printHash(args[0], args[1])
		},
	}
	cmd.Flags().StringVarP(&h.algorithm, flagAlgorithm, "a", "sha256", usageAlgorithm)
	return cmd
}

func init() {
	rootCmd.AddCommand((&hash{}).command())
}