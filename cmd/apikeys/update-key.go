package apikeys

import (
	"fmt"
	"os"

	"cloud.google.com/go/apikeys/apiv2/apikeyspb"
	"github.com/agentio/q/pkg/client"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func updateKeyCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "update-key FILE",
		Short: "Update key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, ctx, err := client.ApiKeysClient(cmd.Context())
			if err != nil {
				return err
			}
			b, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			var key apikeyspb.Key
			err = protojson.Unmarshal(b, &key)
			if err != nil {
				return err
			}
			response, err := c.UpdateKey(ctx, &apikeyspb.UpdateKeyRequest{
				Key: &key,
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
