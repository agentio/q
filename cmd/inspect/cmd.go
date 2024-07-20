package inspect

import (
	"fmt"
	"log"
	"os"

	"github.com/agentio/q/pkg/encoding"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"gopkg.in/yaml.v3"
)

func Cmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Read API information from compiled file descriptors",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bytes, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			var descriptors descriptorpb.FileDescriptorSet
			if err := proto.Unmarshal(bytes, &descriptors); err != nil {
				log.Fatalln("Failed to parse descriptors:", err)
			}
			if format == "json" {
				s, err := protojson.MarshalOptions{
					UseEnumNumbers:  false,
					EmitUnpopulated: true,
					Indent:          "  ",
					UseProtoNames:   false,
				}.Marshal(&descriptors)
				if err != nil {
					return err
				}
				fmt.Printf("%s", string(s))
			} else {
				n, err := encoding.NodeForMessage(&descriptors)
				if err != nil {
					return err
				}
				encoding.StyleForYAML(n)
				b, err := yaml.Marshal(n)
				if err != nil {
					return err
				}
				fmt.Printf("%s", string(b))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
