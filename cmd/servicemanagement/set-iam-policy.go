package servicemanagement

import (
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/iam/apiv1/iampb"
	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func setIamPolicyCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "set-iam-policy RESOURCE FILE",
		Short: "Set iam policy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			b, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}
			var policy iampb.Policy
			err = json.Unmarshal(b, &policy)
			if err != nil {
				return err
			}
			response, err := c.SetIamPolicy(ctx, &iampb.SetIamPolicyRequest{
				Resource: args[0],
				Policy:   &policy,
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
