package servicecontrol

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"cloud.google.com/go/servicecontrol/apiv1/servicecontrolpb"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	ltype "google.golang.org/genproto/googleapis/logging/type"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DefaultAuthScopes reports the default set of authentication scopes to use with this package.
func DefaultAuthScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
}

// adapted from https://github.com/googleapis/google-cloud-go/blob/d40fbff9c1984aeed0224a4ac93eb95c5af17126/translate/apiv3/translation_client.go#L98
func defaultGRPCClientOptions() []option.ClientOption {
	return []option.ClientOption{
		internaloption.WithDefaultEndpoint("servicecontrol.googleapis.com:443"),
		internaloption.WithDefaultEndpointTemplate("servicecontrol.UNIVERSE_DOMAIN:443"),
		internaloption.WithDefaultMTLSEndpoint("servicecontrol.mtls.googleapis.com:443"),
		internaloption.WithDefaultUniverseDomain("googleapis.com"),
		internaloption.WithDefaultAudience("https://servicecontrol.googleapis.com/"),
		internaloption.WithDefaultScopes(DefaultAuthScopes()...),
		internaloption.EnableJwtWithScope(),
		internaloption.EnableNewAuthLibrary(),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

// following https://github.com/GoogleCloudPlatform/esp-v2/blob/469de56a070f50300618dfcfeb590ab8055f5c38/tests/utils/service_control_utils.go#L103
type distOptions struct {
	Buckets int64
	Growth  float64
	Scale   float64
}

func createInt64MetricSet(name string, value int64) *servicecontrolpb.MetricValueSet {
	return &servicecontrolpb.MetricValueSet{
		MetricName: name,
		MetricValues: []*servicecontrolpb.MetricValue{
			{
				Value: &servicecontrolpb.MetricValue_Int64Value{Int64Value: value},
			},
		},
	}
}

var (
	timeDistOptions = distOptions{29, 2.0, 1e-6}
	sizeDistOptions = distOptions{8, 10.0, 1}
)

func createDistMetricSet(options *distOptions, name string, value int64) *servicecontrolpb.MetricValueSet {
	buckets := make([]int64, options.Buckets+2)
	fValue := float64(value)
	idx := 0
	if fValue >= options.Scale {
		idx = 1 + int(math.Log(fValue/options.Scale)/math.Log(options.Growth))
		if idx >= len(buckets) {
			idx = len(buckets) - 1
		}
	}
	buckets[idx] = 1
	distValue := servicecontrolpb.Distribution{
		Count:        1,
		BucketCounts: buckets,
		BucketOption: &servicecontrolpb.Distribution_ExponentialBuckets_{
			ExponentialBuckets: &servicecontrolpb.Distribution_ExponentialBuckets{
				NumFiniteBuckets: int32(options.Buckets),
				GrowthFactor:     options.Growth,
				Scale:            options.Scale,
			},
		},
	}
	if value != 0 {
		distValue.Mean = fValue
		distValue.Minimum = fValue
		distValue.Maximum = fValue
	}
	return &servicecontrolpb.MetricValueSet{
		MetricName: name,
		MetricValues: []*servicecontrolpb.MetricValue{
			{
				Value: &servicecontrolpb.MetricValue_DistributionValue{DistributionValue: &distValue},
			},
		},
	}
}

func mockCmd() *cobra.Command {
	var apikey string
	var format string
	var operation string
	var service string
	cmd := &cobra.Command{
		Use:   "mock",
		Short: "Make service control calls for a mock API request",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			clientOpts := defaultGRPCClientOptions()
			connPool, err := gtransport.DialPool(ctx, clientOpts...)
			if err != nil {
				return err
			}
			sc := servicecontrolpb.NewServiceControllerClient(connPool.Conn())
			qc := servicecontrolpb.NewQuotaControllerClient(connPool.Conn())
			now := time.Now()
			op := &servicecontrolpb.Operation{
				OperationId:   uuid.New().String(),
				OperationName: operation,
				StartTime:     timestamppb.New(now),
				ConsumerId:    "api_key:" + apikey,
			}
			var serviceconfig string
			var producerProject string
			var consumerProject string
			{
				// https://github.com/GoogleCloudPlatform/esp-v2/blob/469de56a070f50300618dfcfeb590ab8055f5c38/src/api_proxy/service_control/request_builder.cc#L1051
				log.Printf("calling check")
				response, err := sc.Check(ctx, &servicecontrolpb.CheckRequest{
					ServiceName: service,
					Operation:   op,
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
				serviceconfig = response.ServiceConfigId
				producerProject = "bobadojo"
				if response.CheckInfo.ConsumerInfo != nil {
					consumerProject = fmt.Sprintf("%d", response.CheckInfo.ConsumerInfo.ConsumerNumber)
				} else {
					return nil
				}
			}
			{
				log.Printf("calling allocate quota")
				// https://github.com/GoogleCloudPlatform/esp-v2/blob/469de56a070f50300618dfcfeb590ab8055f5c38/src/api_proxy/service_control/request_builder.cc#L994
				response, err := qc.AllocateQuota(ctx, &servicecontrolpb.AllocateQuotaRequest{
					ServiceName:     service,
					ServiceConfigId: serviceconfig,
					AllocateOperation: &servicecontrolpb.QuotaOperation{
						OperationId: op.OperationId,
						MethodName:  operation,
						ConsumerId:  op.ConsumerId,
						QuotaMode:   servicecontrolpb.QuotaOperation_NORMAL,
						QuotaMetrics: []*servicecontrolpb.MetricValueSet{
							createInt64MetricSet("serviceruntime.googleapis.com/api/consumer/quota_used_count", 1),
						},
					},
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
			}
			{
				log.Printf("constructing operation to report")
				now = time.Now()
				op.EndTime = timestamppb.New(now)
				status := 200
				callerIP := "10.1.1.1"
				// labels are listed here: https://github.com/GoogleCloudPlatform/esp-v2/blob/469de56a070f50300618dfcfeb590ab8055f5c38/src/api_proxy/service_control/request_builder.cc#L541
				op.Labels = map[string]string{
					"/response_code":                                 fmt.Sprintf("%d", status),
					"/response_code_class":                           fmt.Sprintf("%dxx", status/100),
					"/status_code":                                   fmt.Sprintf("%d", status),
					"/protocol":                                      "http",
					"cloud.googleapis.com/location":                  "global",
					"cloud.googleapis.com/project":                   producerProject,
					"cloud.googleapis.com/service":                   service,
					"cloud.googleapis.com/uid":                       op.OperationId,
					"serviceruntime.googleapis.com/api_method":       operation,
					"serviceruntime.googleapis.com/api_version":      "v1",
					"servicecontrol.googleapis.com/caller_ip":        callerIP,
					"serviceruntime.googleapis.com/consumer_project": consumerProject,
					"servicecontrol.googleapis.com/platform":         "Cloud Run",
					"servicecontrol.googleapis.com/service_agent":    "q/0.0.0",
					"servicecontrol.googleapis.com/user_agent":       "q",
				}
				parts := strings.Split(operation, ".")
				api := strings.Join(parts[:len(parts)-1], ".")
				// example log entry creation: https://github.com/GoogleCloudPlatform/esp-v2/blob/469de56a070f50300618dfcfeb590ab8055f5c38/tests/utils/service_control_utils.go#L307
				payload := map[string]interface{}{
					"api_key":              apikey,
					"api_key_state":        "VERIFIED",
					"api_method":           operation,
					"api_name":             api,
					"api_version":          "v1",
					"grpc_status_code":     "OK",
					"http_status_code":     status,
					"log_message":          operation + " is called",
					"producer_project_id":  producerProject,
					"response_code_detail": "via_upstream",
					"service_agent":        "q/0.0.0",
					"service_config_id":    serviceconfig,
					"timestamp":            now.Unix(),
				}
				st, err := structpb.NewStruct(payload)
				if err != nil {
					log.Fatalf("Unable to marshal payload %v", err)
				}
				op.LogEntries = []*servicecontrolpb.LogEntry{
					{
						Name:      "endpoints_log",
						Timestamp: timestamppb.New(now),
						Severity:  ltype.LogSeverity_INFO,
						HttpRequest: &servicecontrolpb.HttpRequest{
							Latency:       &durationpb.Duration{Seconds: 5},
							Protocol:      "grpc",
							RemoteIp:      callerIP,
							RequestMethod: "GET",
							RequestSize:   10,
							RequestUrl:    "/" + api + "/" + parts[len(parts)-1],
							ResponseSize:  10,
							Status:        int32(status),
						},
						Payload: &servicecontrolpb.LogEntry_StructPayload{
							StructPayload: st,
						},
					},
				}
				// metric set values are computed here: https://github.com/GoogleCloudPlatform/esp-v2/blob/469de56a070f50300618dfcfeb590ab8055f5c38/tests/utils/service_control_utils.go#L410
				op.MetricValueSets = []*servicecontrolpb.MetricValueSet{
					createInt64MetricSet("serviceruntime.googleapis.com/api/consumer/request_count", 1),
					createInt64MetricSet("serviceruntime.googleapis.com/api/producer/request_count", 1),
					createInt64MetricSet("serviceruntime.googleapis.com/api/consumer/quota_used_count", 1),
					createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/consumer/total_latencies", 50),
					createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/producer/total_latencies", 50),
					createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/consumer/request_sizes", 100),
					createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/consumer/response_sizes", 100),
					createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/producer/request_overhead_latencies", 50),
					createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/producer/backend_latencies", 50),
					createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/producer/request_sizes", 100),
					createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/producer/response_sizes", 100),
				}
			}
			{
				log.Printf("calling report")
				// https://github.com/GoogleCloudPlatform/esp-v2/blob/469de56a070f50300618dfcfeb590ab8055f5c38/src/api_proxy/service_control/request_builder.cc#L1093
				response, err := sc.Report(ctx, &servicecontrolpb.ReportRequest{
					ServiceName: service,
					Operations:  []*servicecontrolpb.Operation{op},
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
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	cmd.Flags().StringVar(&apikey, "apikey", "", "API key")
	cmd.Flags().StringVar(&operation, "operation", "", "Operation name")
	cmd.Flags().StringVar(&service, "service", "", "Service name")
	return cmd
}
