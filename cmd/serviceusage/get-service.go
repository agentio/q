package serviceusage

import "github.com/spf13/cobra"

func getServiceCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "get-service",
		Short: "Get service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}
