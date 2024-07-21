package compile

import (
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/label"
	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/sourcecontextpb"
	"google.golang.org/protobuf/types/known/typepb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func AddDetailFromDescriptors(c *serviceconfig.Service, d *descriptorpb.FileDescriptorSet) {
	c.Types = []*typepb.Type{}
	c.Http = &annotations.Http{}
	c.Backend = &serviceconfig.Backend{}
	for _, api := range c.Apis {
		AddAPIDetailFromDescriptors(api, c, d)
	}
	c.Quota = &serviceconfig.Quota{}
	c.Authentication = &serviceconfig.Authentication{}
}

func AddAPIDetailFromDescriptors(api *apipb.Api, c *serviceconfig.Service, d *descriptorpb.FileDescriptorSet) {

	for _, file := range d.File {
		for _, service := range file.Service {
			fullName := (*file.Package) + "." + *(service.Name)
			if fullName == api.Name {
				parts := strings.Split(*file.Package, ".")
				api.Version = parts[len(parts)-1]
				api.SourceContext = &sourcecontextpb.SourceContext{
					FileName: *file.Name,
				}
				if *file.Syntax == "proto3" {
					api.Syntax = typepb.Syntax_SYNTAX_PROTO3
				}
				for _, method := range service.Method {
					// TODO: backend rules should only be added if they aren't already user-specified
					c.Backend.Rules = append(c.Backend.Rules, &serviceconfig.BackendRule{
						Selector: *file.Package + "." + *service.Name + "." + *method.Name,
					})
					options := []*typepb.Option{}
					method.Options.ProtoReflect().Range(func(ext protoreflect.FieldDescriptor, v protoreflect.Value) bool {
						var value *anypb.Any
						k := ext.Kind().String()
						cardinality := ext.Cardinality().String()
						if k == "string" && cardinality == "repeated" {
							value, _ = anypb.New(wrapperspb.String(v.List().Get(0).String()))
						} else if k == "message" {
							value, _ = anypb.New(v.Message().Interface())
						}
						options = append(options, &typepb.Option{
							Name:  string(ext.FullName()),
							Value: value,
						})
						// collect http rules into the "http" section of the service config
						if string(ext.FullName()) == "google.api.http" {

							switch h := v.Message().Interface().(type) {
							case *annotations.HttpRule:
								h.Selector = (*file.Package) + "." + (*service.Name) + "." + (*method.Name)
								c.Http.Rules = append(c.Http.Rules, h)
							default:
								// nothing
							}
						}
						return true
					})
					apiMethod := &apipb.Method{
						Name:            *(method.Name),
						RequestTypeUrl:  "type.googleapis.com/" + strings.TrimLeft(*(method.InputType), "."),
						ResponseTypeUrl: "type.googleapis.com/" + strings.TrimLeft(*(method.OutputType), "."),
						Options:         options,
					}
					api.Methods = append(api.Methods, apiMethod)
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
