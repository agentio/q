package servicemanagement

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-management",
		Short: "Manage service descriptions with the Google Service Management API",
	}
	cmd.AddCommand(createServiceCmd())
	cmd.AddCommand(createServiceConfigCmd())
	cmd.AddCommand(createServiceRolloutCmd())
	cmd.AddCommand(deleteServiceCmd())
	cmd.AddCommand(generateConfigReportCmd())
	cmd.AddCommand(getOperationCmd())
	cmd.AddCommand(getServiceCmd())
	cmd.AddCommand(getServiceConfigCmd())
	cmd.AddCommand(getServiceRolloutCmd())
	cmd.AddCommand(listServiceConfigsCmd())
	cmd.AddCommand(listServiceRolloutsCmd())
	cmd.AddCommand(listServicesCmd())
	cmd.AddCommand(submitConfigSourceCmd())
	cmd.AddCommand(undeleteServiceCmd())
	return cmd
}
