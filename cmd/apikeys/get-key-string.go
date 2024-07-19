package apikeys

import (
	"github.com/spf13/cobra"
)

func getKeyStringCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "get-key-string",
		Short: "Get key string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
