package monitoring

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Monitor services with the Google Cloud Monitoring API",
	}
	return cmd
}
