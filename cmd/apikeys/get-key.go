package apikeys

import (
	"fmt"

	"cloud.google.com/go/apikeys/apiv2/apikeyspb"
	"github.com/agent-kit/q/pkg/client"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func getKeyCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "get-key",
		Short: "Get key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, ctx, err := client.ApiKeysClient(cmd.Context())
			if err != nil {
				return err
			}
			response, err := c.GetKey(ctx, &apikeyspb.GetKeyRequest{
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
