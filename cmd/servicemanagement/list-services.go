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
	var project string
	var output string
	cmd := &cobra.Command{
		Use:   "list-services",
		Short: "List services",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()

			response := c.ListServices(ctx, &servicemanagementpb.ListServicesRequest{
				ProducerProjectId: project,
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
	cmd.Flags().StringVarP(&project, "project", "p", "", "producer project")
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}
