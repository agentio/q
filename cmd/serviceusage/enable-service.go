package serviceusage

import "github.com/spf13/cobra"

func enableServiceCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "enable-service",
		Short: "Enable service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}
