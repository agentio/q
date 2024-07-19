package apikeys

import (
	"github.com/spf13/cobra"
)

func createKeyCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "create-key",
		Short: "Create key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}
