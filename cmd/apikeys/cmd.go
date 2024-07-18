package apikeys

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-keys",
		Short: "Manage API keys with the Google API Keys API",
	}
	return cmd
}
