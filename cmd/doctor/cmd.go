package doctor

import (
	"fmt"

	"github.com/agentio/q/pkg/gcloud"
	"github.com/spf13/cobra"
)

var verbose bool

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Verify necessary dependencies and configuration",
		RunE:  action,
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "display intermediate information")
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Let the doctor's recommendations stand out.
	fmt.Printf("Checking gcloud configuration...\n")
	info, err := gcloud.GetInfo(verbose)
	if err != nil {
		return err
	}
	// there must be a logged-in account
	account, err := info.Account()
	if err != nil {
		return err
	}
	fmt.Printf("account = %s\n", account)
	// a project must be set
	project, err := info.Project()
	if err != nil {
		return err
	}
	fmt.Printf("project = %s\n", project)
	// a run region must be set
	runregion, err := info.RunRegion()
	if err != nil {
		return err
	}
	fmt.Printf("run/region = %s\n", runregion)
	fmt.Printf("Everything looks good!\n")
	return nil
}
