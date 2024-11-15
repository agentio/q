package servicemanagement

import (
	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"github.com/spf13/cobra"
)

func deleteOperationCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "delete-operation OPERATION",
		Short: "Delete operation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()
			err = c.LROClient.DeleteOperation(ctx, &longrunningpb.DeleteOperationRequest{
				Name: args[0],
			})
			return err
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
