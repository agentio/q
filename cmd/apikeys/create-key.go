package apikeys

import (
	"errors"
	"fmt"

	"cloud.google.com/go/apikeys/apiv2/apikeyspb"
	"github.com/agentio/q/pkg/client"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func createKeyCmd() *cobra.Command {
	var format string
	var parent string
	var service string
	var keyid string
	var displayName string
	cmd := &cobra.Command{
		Use:   "create-key",
		Short: "Create key",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, ctx, err := client.ApiKeysClient(cmd.Context())
			if err != nil {
				return err
			}
			if parent == "" {
				return errors.New("--parent must be specified")
			}
			if service == "" {
				return errors.New("--service must be specified")
			}
			response, err := c.CreateKey(ctx, &apikeyspb.CreateKeyRequest{
				Parent: parent,
				Key: &apikeyspb.Key{
					DisplayName: displayName,
					Restrictions: &apikeyspb.Restrictions{
						ApiTargets: []*apikeyspb.ApiTarget{
							{
								Service: service,
							},
						},
					},
				},
				KeyId: keyid,
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
	cmd.Flags().StringVar(&parent, "parent", "", "parent (projects/PROJECTNAME)")
	cmd.Flags().StringVar(&service, "service", "", "service to be used with this key")
	cmd.Flags().StringVar(&keyid, "keyid", "", "key id")
	cmd.Flags().StringVar(&displayName, "display-name", "", "display name")
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
