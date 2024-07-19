package servicemanagement

import (
	"fmt"

	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"cloud.google.com/go/servicemanagement/apiv1/servicemanagementpb"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/encoding/protojson"
)

func listServicesCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "list-services PROJECTID",
		Short: "List services",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			response := c.ListServices(ctx, &servicemanagementpb.ListServicesRequest{
				ProducerProjectId: args[0],
			})
			if format == "json" {
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
				if format == "json" {
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
			if format == "json" {
				fmt.Fprintf(cmd.OutOrStdout(), "]\n")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
