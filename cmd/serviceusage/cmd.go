package serviceusage

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-usage",
		Short: "Manage usage of APIs with the Google Service Usage API",
	}
	return cmd
}
