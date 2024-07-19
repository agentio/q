package apikeys

import (
	"fmt"

	"github.com/agentio/q/pkg/client"
	"github.com/spf13/cobra"

	"cloud.google.com/go/apikeys/apiv2/apikeyspb"
	"google.golang.org/protobuf/encoding/protojson"
)

func listKeysCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "list-keys PROJECT",
		Short: "List keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := args[0]

			c, ctx, err := client.ApiKeysClient(cmd.Context(), project)
			if err != nil {
				return err
			}

			nextPageToken := ""
			for {
				response, err := c.ListKeys(ctx, &apikeyspb.ListKeysRequest{
					Parent:    "projects/" + project + "/locations/global",
					PageSize:  2,
					PageToken: nextPageToken,
				})
				if err != nil {
					return err
				}
				if output == "json" {
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
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}
