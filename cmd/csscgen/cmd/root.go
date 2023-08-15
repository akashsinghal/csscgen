package cmd

import (
	"github.com/spf13/cobra"
)

const (
	use       = "csscgen"
	shortDesc = "csscgen is a tool used to generate k8s templates & the artifacts for supply chain load testing"
)

var Root = New(use, shortDesc)

func New(use, short string) *cobra.Command {
	root := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	root.AddCommand(NewCmdGenK8s(use))

	return root
}
