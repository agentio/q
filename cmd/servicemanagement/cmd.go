package servicemanagement

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-management",
		Short: "Manage service descriptions with the Service Management API",
	}
	cmd.AddCommand(cancelOperationCmd())
	cmd.AddCommand(createServiceCmd())
	cmd.AddCommand(createServiceConfigCmd())
	cmd.AddCommand(createServiceRolloutCmd())
	cmd.AddCommand(deleteOperationCmd())
	cmd.AddCommand(deleteServiceCmd())
	cmd.AddCommand(generateConfigReportCmd())
	cmd.AddCommand(getIamPolicyCmd())
	cmd.AddCommand(getOperationCmd())
	cmd.AddCommand(getServiceCmd())
	cmd.AddCommand(getServiceConfigCmd())
	cmd.AddCommand(getServiceRolloutCmd())
	cmd.AddCommand(listOperationsCmd())
	cmd.AddCommand(listServiceConfigsCmd())
	cmd.AddCommand(listServiceRolloutsCmd())
	cmd.AddCommand(listServicesCmd())
	cmd.AddCommand(setIamPolicyCmd())
	cmd.AddCommand(submitConfigSourceCmd())
	cmd.AddCommand(testIamPermissionsCmd())
	cmd.AddCommand(undeleteServiceCmd())
	cmd.AddCommand(waitOperationCmd())
	return cmd
}
