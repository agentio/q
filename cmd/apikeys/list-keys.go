package apikeys

import (
	"github.com/spf13/cobra"
)

func listKeysCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "list-keys",
		Short: "List keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}
