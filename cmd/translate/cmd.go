package translate

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "translate",
		Short: "Translate with the Google Cloud Translation API",
	}
	return cmd
}
