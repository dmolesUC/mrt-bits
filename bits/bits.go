package bits

import (
	"github.com/spf13/cobra"
)

// ------------------------------------------------------------
// Exported symbols

func Execute() error {
	return command.Execute()
}

func AddCommand(c *cobra.Command) {
	// TODO: standard flags
	command.AddCommand(c)
}

// ------------------------------------------------------------
// Unexported symbols

const (
	usage     = "bits"
	shortDesc = "Merritt bitstream service (experimental)"
)

var command = &cobra.Command{Use: usage, Short: shortDesc}
