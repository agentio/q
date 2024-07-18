package servicecontrol

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-control",
		Short: "Add access control and telemetry to APIs with the Google Service Control API",
	}
	return cmd
}
