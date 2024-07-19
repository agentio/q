package apikeys

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-keys",
		Short: "Manage API keys with the API Keys API",
	}
	cmd.AddCommand(createKeyCmd())
	cmd.AddCommand(deleteKeyCmd())
	cmd.AddCommand(getKeyStringCmd())
	cmd.AddCommand(getKeyCmd())
	cmd.AddCommand(listKeysCmd())
	cmd.AddCommand(lookupKeyCmd())
	cmd.AddCommand(undeleteKeyCmd())
	cmd.AddCommand(updateKeyCmd())
	return cmd
}
