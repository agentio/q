package compile

import (
	"log"
	"strings"

	"google.golang.org/genproto/googleapis/api/label"
	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/typepb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func AddDetailFromDescriptors(c *serviceconfig.Service, d *descriptorpb.FileDescriptorSet) {
	for _, api := range c.Apis {

		AddAPIDetailFromDescriptors(api, d)

	}
}

func AddAPIDetailFromDescriptors(api *apipb.Api, d *descriptorpb.FileDescriptorSet) {
	log.Printf("%s", api.Name)

	for _, file := range d.File {
		for _, service := range file.Service {
			fullName := (*file.Package) + "." + *(service.Name)
			if fullName == api.Name {
				log.Printf("-- %s", fullName)

				for _, method := range service.Method {

					log.Printf("options: %+v", method.Options)

					options := []*typepb.Option{}
					method.Options.ProtoReflect().Range(func(ext protoreflect.FieldDescriptor, v protoreflect.Value) bool {
						//log.Printf("ranging %T %T ext=%+v value=%+v", ext, v, ext, v)

						var value *anypb.Any
						log.Printf("%s", ext.FullName())
						k := ext.Kind().String()
						cardinality := ext.Cardinality().String()
						if k == "string" && cardinality == "repeated" {
							log.Printf("list of strings %s", v.List().Get(0).String())

							value, _ = anypb.New(wrapperspb.String(v.List().Get(0).String()))

						} else if k == "message" {
							log.Printf("message %T %+v", v.Message().Interface(), v.Message().Interface())

							value, _ = anypb.New(v.Message().Interface())

						}
						options = append(options, &typepb.Option{
							Name:  string(ext.FullName()),
							Value: value,
						})

						return true
					})

					apiMethod := &apipb.Method{
						Name:            *(method.Name),
						RequestTypeUrl:  "type.googleapis.com/" + strings.TrimLeft(*(method.InputType), "."),
						ResponseTypeUrl: "type.googleapis.com/" + strings.TrimLeft(*(method.OutputType), "."),
						Options:         options,
					}
					api.Methods = append(api.Methods, apiMethod)

					log.Printf("---- %s", *(method.Name))

				}

			}
		}

	}

}

func AddCommonEndpointsSettings(c *serviceconfig.Service) {
	c.Control = &serviceconfig.Control{
		Environment: "servicecontrol.googleapis.com",
	}
	c.Logs = []*serviceconfig.LogDescriptor{
		{
			Name: "endpoints_log",
		},
	}
	c.Metrics = []*metric.MetricDescriptor{
		{
			Name: "serviceruntime.googleapis.com/api/consumer/request_count",
			Type: "serviceruntime.googleapis.com/api/consumer/request_count",
			Labels: []*label.LabelDescriptor{
				{Key: "/credential_id"},
				{Key: "/protocol"},
				{Key: "/response_code"},
				{Key: "/response_code_class"},
				{Key: "/status_code"},
			},
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_INT64,
		},
		{
			Name: "serviceruntime.googleapis.com/api/consumer/total_latencies",
			Type: "serviceruntime.googleapis.com/api/consumer/total_latencies",
			Labels: []*label.LabelDescriptor{
				{Key: "/credential_id"},
			},
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_DISTRIBUTION,
		},
		{
			Name: "serviceruntime.googleapis.com/api/producer/request_count",
			Type: "serviceruntime.googleapis.com/api/producer/request_count",
			Labels: []*label.LabelDescriptor{
				{Key: "/protocol"},
				{Key: "/response_code"},
				{Key: "/response_code_class"},
				{Key: "/status_code"},
			},
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_INT64,
		},
		{
			Name:       "serviceruntime.googleapis.com/api/producer/total_latencies",
			Type:       "serviceruntime.googleapis.com/api/producer/total_latencies",
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_DISTRIBUTION,
		},
		{
			Name: "serviceruntime.googleapis.com/api/consumer/quota_used_count",
			Type: "serviceruntime.googleapis.com/api/consumer/quota_used_count",
			Labels: []*label.LabelDescriptor{
				{Key: "/credential_id"},
				{Key: "/quota_group_name"},
			},
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_INT64,
		},
		{
			Name: "serviceruntime.googleapis.com/api/consumer/request_sizes",
			Type: "serviceruntime.googleapis.com/api/consumer/request_sizes",
			Labels: []*label.LabelDescriptor{
				{Key: "/credential_id"},
			},
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_DISTRIBUTION,
		},
		{
			Name: "serviceruntime.googleapis.com/api/consumer/response_sizes",
			Type: "serviceruntime.googleapis.com/api/consumer/response_sizes",
			Labels: []*label.LabelDescriptor{
				{Key: "/credential_id"},
			},
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_DISTRIBUTION,
		},
		{
			Name:       "serviceruntime.googleapis.com/api/producer/request_overhead_latencies",
			Type:       "serviceruntime.googleapis.com/api/producer/request_overhead_latencies",
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_DISTRIBUTION,
		},
		{
			Name:       "serviceruntime.googleapis.com/api/producer/backend_latencies",
			Type:       "serviceruntime.googleapis.com/api/producer/backend_latencies",
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_DISTRIBUTION,
		},
		{
			Name:       "serviceruntime.googleapis.com/api/producer/request_sizes",
			Type:       "serviceruntime.googleapis.com/api/producer/request_sizes",
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_DISTRIBUTION,
		},
		{
			Name:       "serviceruntime.googleapis.com/api/producer/response_sizes",
			Type:       "serviceruntime.googleapis.com/api/producer/response_sizes",
			MetricKind: metric.MetricDescriptor_DELTA,
			ValueType:  metric.MetricDescriptor_DISTRIBUTION,
		},
	}
	c.MonitoredResources = []*monitoredres.MonitoredResourceDescriptor{
		{
			Type: "api",
			Labels: []*label.LabelDescriptor{
				{Key: "cloud.googleapis.com/location"},
				{Key: "cloud.googleapis.com/uid"},
				{Key: "serviceruntime.googleapis.com/api_version"},
				{Key: "serviceruntime.googleapis.com/api_method"},
				{Key: "serviceruntime.googleapis.com/consumer_project"},
				{Key: "cloud.googleapis.com/project"},
				{Key: "cloud.googleapis.com/service"},
			},
		},
	}
	c.Logging = &serviceconfig.Logging{
		ProducerDestinations: []*serviceconfig.Logging_LoggingDestination{
			{
				MonitoredResource: "api",
				Logs:              []string{"endpoints_log"},
			},
		},
	}
	c.Monitoring = &serviceconfig.Monitoring{
		ProducerDestinations: []*serviceconfig.Monitoring_MonitoringDestination{
			{
				MonitoredResource: "api",
				Metrics: []string{
					"serviceruntime.googleapis.com/api/producer/request_count",
					"serviceruntime.googleapis.com/api/producer/total_latencies",
					"serviceruntime.googleapis.com/api/producer/request_overhead_latencies",
					"serviceruntime.googleapis.com/api/producer/backend_latencies",
					"serviceruntime.googleapis.com/api/producer/request_sizes",
					"serviceruntime.googleapis.com/api/producer/response_sizes",
				},
			},
		},
		ConsumerDestinations: []*serviceconfig.Monitoring_MonitoringDestination{
			{
				MonitoredResource: "api",
				Metrics: []string{
					"serviceruntime.googleapis.com/api/consumer/request_count",
					"serviceruntime.googleapis.com/api/consumer/quota_used_count",
					"serviceruntime.googleapis.com/api/consumer/total_latencies",
					"serviceruntime.googleapis.com/api/consumer/request_sizes",
					"serviceruntime.googleapis.com/api/consumer/response_sizes",
				},
			},
		},
	}
	c.SystemParameters = &serviceconfig.SystemParameters{}
}
