package servicecontrol

import (
	"fmt"
	"math"
	"time"

	"cloud.google.com/go/servicecontrol/apiv1/servicecontrolpb"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DefaultAuthScopes reports the default set of authentication scopes to use with this package.
func DefaultAuthScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
}

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

func mockCmd() *cobra.Command {
	var apikey string
	var format string
	var operation string
	var service string
	cmd := &cobra.Command{
		Use:   "mock",
		Short: "mock",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			clientOpts := defaultGRPCClientOptions()
			connPool, err := gtransport.DialPool(ctx, clientOpts...)
			if err != nil {
				return err
			}
			c := servicecontrolpb.NewServiceControllerClient(connPool.Conn())
			now := time.Now()
			timestamp := timestamppb.New(now)
			response, err := c.Check(ctx, &servicecontrolpb.CheckRequest{
				ServiceName: service,
				Operation: &servicecontrolpb.Operation{
					StartTime:     timestamp,
					OperationName: operation,
					ConsumerId:    "api_key:" + apikey,
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
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	cmd.Flags().StringVar(&apikey, "apikey", "", "API key")
	cmd.Flags().StringVar(&operation, "operation", "", "Operation name")
	cmd.Flags().StringVar(&service, "service", "", "Service name")
	return cmd
}
