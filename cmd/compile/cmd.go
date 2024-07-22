package compile

import (
	"fmt"
	"os"

	"github.com/agent-kit/q/pkg/compile"
	"github.com/agent-kit/q/pkg/encoding"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/encoding/protojson"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile SERVICE DESCRIPTOR",
		Short: "Compile a Service Configuration for an API",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := readServiceConfig(args[0])
			if err != nil {
				return err
			}
			compile.CompileDescriptor(config, args[1])
			compile.AddCommonEndpointsSettings(config)
			bytes, err := protojson.Marshal(config)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", string(bytes))
			return nil
		},
	}
	return cmd
}

func readServiceConfig(filename string) (*serviceconfig.Service, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	message, err := encoding.UnmarshalYaml(bytes)
	if err != nil {
		return nil, err
	}
	switch v := message.(type) {
	case *serviceconfig.Service:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type %t", v)
	}
}
