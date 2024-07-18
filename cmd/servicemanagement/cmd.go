package servicemanagement

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-management",
		Short: "Publish and manage information about services with the Google Service Management API",
	}
	return cmd
}
