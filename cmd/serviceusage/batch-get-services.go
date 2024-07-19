package serviceusage

import "github.com/spf13/cobra"

func batchGetServicesCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "batch-get-services",
		Short: "Batch get services",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
