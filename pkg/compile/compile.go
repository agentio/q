package compile

import (
	"log"
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
	c.Http = &annotations.Http{}
	c.Backend = &serviceconfig.Backend{}
	for _, api := range c.Apis {
		typeMap := AddAPIDetailFromDescriptors(api, c, d)
		log.Printf("%+v", typeMap)
	}
	c.Quota = &serviceconfig.Quota{}
	c.Authentication = &serviceconfig.Authentication{}
	c.Types = []*typepb.Type{}
	allTypes := CollectTypesFromDescriptors(d)

	//log.Printf("%+v", allTypes)

	for _, v := range allTypes {
		c.Types = append(c.Types, v)
	}
}

func cardinalityForLabel(l *descriptorpb.FieldDescriptorProto_Label) typepb.Field_Cardinality {
	if l == nil {
		return typepb.Field_CARDINALITY_OPTIONAL
	} else {
		switch *l {
		case descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL:
			return typepb.Field_CARDINALITY_OPTIONAL
		case descriptorpb.FieldDescriptorProto_LABEL_REPEATED:
			return typepb.Field_CARDINALITY_REPEATED
		default:
			return typepb.Field_CARDINALITY_UNKNOWN
		}
	}
}

func CollectTypesFromDescriptors(d *descriptorpb.FileDescriptorSet) map[string]*typepb.Type {
	types := make(map[string]*typepb.Type)

	for _, file := range d.File {
		for _, message := range file.MessageType {
			log.Printf("IN %+v", message)

			fields := []*typepb.Field{}
			for _, f := range message.Field {
				field := &typepb.Field{
					Kind:        typepb.Field_Kind(*f.Type),
					Cardinality: cardinalityForLabel(f.Label),
					Name:        *f.Name,
					Number:      *f.Number,
					JsonName:    *f.JsonName,
				}
				if field.Kind == typepb.Field_TYPE_MESSAGE {
					field.TypeUrl = "type.googleapis.com/" + strings.TrimLeft(*f.TypeName, ".")
				}
				f.Options.ProtoReflect().Range(func(ext protoreflect.FieldDescriptor, v protoreflect.Value) bool {
					if string(ext.FullName()) == "google.api.field_behavior" {
						i := v.List().Get(0).Enum()
						var s string
						switch i {
						case 1:
							s = "OPTIONAL"
						case 2:
							s = "REQUIRED"
						default:
							s = "UNKNOWN"
						}
						//log.Printf("WHAT IS THIS %T %+v", v.List().Get(0), v.List().Get(0))
						a, _ := anypb.New(wrapperspb.String(s))
						field.Options = append(field.Options, &typepb.Option{
							Name:  string(ext.FullName()),
							Value: a,
						})
					} else if string(ext.FullName()) == "google.api.resource_reference" {
						//log.Printf("WHAT IS THIS %T %+v", v.Message().Interface(), v.Message().Interface())
						a, _ := anypb.New(v.Message().Interface())
						field.Options = append(field.Options, &typepb.Option{
							Name:  string(ext.FullName()),
							Value: a,
						})

					}
					return true
				})
				fields = append(fields, field)
			}

			t := &typepb.Type{
				Name:   *file.Package + "." + *message.Name,
				Fields: fields,
			}
			if file.Name != nil {
				t.SourceContext = &sourcecontextpb.SourceContext{
					FileName: *file.Name,
				}
			}
			if file.Syntax != nil {
				switch *file.Syntax {
				case "proto2":
					t.Syntax = typepb.Syntax_SYNTAX_PROTO2
				case "proto3":
					t.Syntax = typepb.Syntax_SYNTAX_PROTO3
				case "editions":
					t.Syntax = typepb.Syntax_SYNTAX_EDITIONS
				}
			}

			log.Printf("OUT %+v", t)
			types[t.Name] = t
		}
	}
	return types
}

func AddAPIDetailFromDescriptors(api *apipb.Api, c *serviceconfig.Service, d *descriptorpb.FileDescriptorSet) map[string]bool {
	typeMap := make(map[string]bool)

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
					// note the input/output types as interesting
					if method.InputType != nil {
						typeMap[strings.TrimLeft(*method.InputType, ".")] = true
					}
					if method.OutputType != nil {
						typeMap[strings.TrimLeft(*method.OutputType, ".")] = true
					}
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
	return typeMap
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
