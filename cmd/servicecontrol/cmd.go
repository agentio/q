package servicecontrol

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-control",
		Short: "Control API services with the Google Service Control API",
	}
	return cmd
}
