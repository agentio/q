package servicemanagement

import (
	"fmt"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func getOperationCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "get-operation OPERATION",
		Short: "Get operation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()
			response, err := c.LROClient.GetOperation(ctx, &longrunningpb.GetOperationRequest{
				Name: args[0],
			})
			if err != nil {
				return err
			}
			if format == "json" {
				b, err := protojson.Marshal(response)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(b))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
