package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	resolveUse = "resolve"
)

type resolveCmdOptions struct {
	configFilePath string
	subject        string
}

func NewCmdGenK8s(argv ...string) *cobra.Command {
	if len(argv) == 0 {
		argv = []string{os.Args[0]}
	}

	eg := fmt.Sprintf(`  # Generates a kubernetes resource template`)

	// var opts resolveCmdOptions

	cmd := &cobra.Command{
		Use:     resolveUse,
		Short:   "Generates a kubernetes resource template",
		Example: eg,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// flags := cmd.Flags()

	// flags.StringVarP(&opts.subject, "subject", "s", "", "Subject Reference")
	// flags.StringVarP(&opts.configFilePath, "config", "c", "", "Config File Path")
	return cmd
}
