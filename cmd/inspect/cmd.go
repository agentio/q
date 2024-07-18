package inspect

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Read API information from compiled file descriptors",
	}
	return cmd
}
