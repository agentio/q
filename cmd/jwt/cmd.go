package jwt

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "Read, verify, and generate JSON Web Tokens",
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(generateCmd())
	cmd.AddCommand(verifyCmd())
	return cmd
}
