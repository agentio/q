package serviceusage

import (
	"fmt"

	serviceusage "cloud.google.com/go/serviceusage/apiv1"
	"cloud.google.com/go/serviceusage/apiv1/serviceusagepb"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/encoding/protojson"
)

func listServicesCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "list-services",
		Short: "List services",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := serviceusage.NewClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()

			response := c.ListServices(ctx, &serviceusagepb.ListServicesRequest{
				Parent: args[0],
			})
			if output == "json" {
				fmt.Fprintf(cmd.OutOrStdout(), "[")
			}
			first := true
			for {
				s, err := response.Next()
				if err == iterator.Done {
					break
				} else if err != nil {
					return err
				}
				if output == "json" {
					if first {
						first = false
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), ",")
					}
					b, err := protojson.Marshal(s)
					if err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%s", string(b))
				}
			}
			if output == "json" {
				fmt.Fprintf(cmd.OutOrStdout(), "]\n")
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}