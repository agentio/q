package monitoring

import (
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func listTimeSeries() *cobra.Command {
	var format string
	var filter string
	var d int
	cmd := &cobra.Command{
		Use:   "list-time-series PROJECT METRIC",
		Args:  cobra.ExactArgs(2),
		Short: "read test values from the Cloud Monitoring API",
		Long: `sample metrics:
custom.googleapis.com/stores/daily_sales
serviceruntime.googleapis.com/api/request_count
iam.googleapis.com/service_account/key/authn_events_count
serviceruntime.googleapis.com/api/producer/response_sizes
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Creates a client.
			c, err := monitoring.NewMetricClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}
			defer c.Close()

			// Sets your Google Cloud Platform project ID.
			projectID := args[0]
			metric := args[1]

			response := c.ListTimeSeries(ctx, &monitoringpb.ListTimeSeriesRequest{
				Name: "projects/" + projectID,
				Interval: &monitoringpb.TimeInterval{
					EndTime:   timestamppb.Now(),
					StartTime: timestamppb.New(time.Now().Add(-time.Duration(d) * time.Second)),
				},
				View:   monitoringpb.ListTimeSeriesRequest_FULL,
				Filter: `metric.type = "` + metric + `"` + filter,
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
	cmd.Flags().StringVar(&filter, "filter", "", "additional filter expression")
	cmd.Flags().IntVarP(&d, "duration", "d", 3600, "duration of time to query (to now)")
	return cmd
}
