package monitoring

import (
	"fmt"
	"strconv"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spf13/cobra"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
)

func writeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write PROJECT METRIC VALUE",
		Args:  cobra.ExactArgs(3),
		Short: "write test values to the Cloud Monitoring API",
		Long: `sample metrics: 
custom.googleapis.com/stores/daily_sales
custom.googleapis.com/stores/temp`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Creates a client.
			client, err := monitoring.NewMetricClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}

			// Sets your Google Cloud Platform project ID.
			projectID := args[0]
			metric := args[1]
			value := args[2]
			d, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}

			// Prepares an individual data point
			dataPoint := &monitoringpb.Point{
				Interval: &monitoringpb.TimeInterval{
					EndTime: &timestamppb.Timestamp{
						Seconds: time.Now().Unix(),
					},
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_DoubleValue{
						DoubleValue: d,
					},
				},
			}

			// Writes time series data.
			if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
				Name: fmt.Sprintf("projects/%s", projectID),
				TimeSeries: []*monitoringpb.TimeSeries{
					{
						Metric: &metricpb.Metric{
							Type:   metric,
							Labels: map[string]string{},
						},
						Resource: &monitoredrespb.MonitoredResource{
							Type: "global",
							Labels: map[string]string{
								"project_id": projectID,
							},
						},
						MetricKind: metricpb.MetricDescriptor_GAUGE,
						Points: []*monitoringpb.Point{
							dataPoint,
						},
					},
				},
			}); err != nil {
				return fmt.Errorf("failed to write time series data: %v", err)
			}

			// Closes the client and flushes the data to Stackdriver.
			if err := client.Close(); err != nil {
				return fmt.Errorf("failed to close client: %v", err)
			}

			fmt.Printf("Done writing time series data.\n")
			return nil
		},
	}
	return cmd
}
