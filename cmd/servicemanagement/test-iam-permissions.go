package servicemanagement

import (
	"fmt"

	"cloud.google.com/go/iam/apiv1/iampb"
	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func testIamPermissionsCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "test-iam-permissions RESOURCE PERMISSION",
		Short: "Test iam permission",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			response, err := c.TestIamPermissions(ctx, &iampb.TestIamPermissionsRequest{
				Resource:    args[0],
				Permissions: []string{args[1]},
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
