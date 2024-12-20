package apikeys

import (
	"fmt"

	"github.com/agentio/q/pkg/client"
	"github.com/spf13/cobra"

	"cloud.google.com/go/apikeys/apiv2/apikeyspb"
	"google.golang.org/protobuf/encoding/protojson"
)

func listKeysCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "list-keys PROJECT",
		Short: "List keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, ctx, err := client.ApiKeysClient(cmd.Context())
			if err != nil {
				return err
			}
			nextPageToken := ""
			for {
				response, err := c.ListKeys(ctx, &apikeyspb.ListKeysRequest{
					Parent:    "projects/" + args[0] + "/locations/global",
					PageSize:  100,
					PageToken: nextPageToken,
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
				nextPageToken = response.NextPageToken
				if nextPageToken == "" {
					break
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
