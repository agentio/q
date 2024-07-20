package apikeys

import (
	"fmt"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/agent-kit/q/pkg/client"
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
			c, ctx, err := client.ApiKeysLROClient(cmd.Context())
			if err != nil {
				return err
			}
			response, err := c.GetOperation(ctx, &longrunningpb.GetOperationRequest{
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
