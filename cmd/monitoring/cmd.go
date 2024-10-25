package monitoring

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Monitor services with the Cloud Monitoring API",
	}
	cmd.AddCommand(readCmd())
	cmd.AddCommand(writeCmd())
	return cmd
}
