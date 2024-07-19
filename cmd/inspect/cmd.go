package inspect

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

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
				n, err := NodeForMessage(&descriptors)
				if err != nil {
					return err
				}
				StyleForYAML(n)
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

// Prefer this encoder because it uses tighter 2-space indentation.
func yamlEncoder(dst io.Writer) *yaml.Encoder {
	enc := yaml.NewEncoder(dst)
	enc.SetIndent(2)
	return enc
}

// Encode a model as YAML.
func EncodeYAML(model interface{}) ([]byte, error) {
	var b bytes.Buffer
	err := yamlEncoder(&b).Encode(model)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// NodeForMessage converts a proto message for export as a YAML node.
func NodeForMessage(m proto.Message) (*yaml.Node, error) {
	// Marshal the artifact content as JSON using the protobuf marshaller.
	var s []byte
	s, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		Indent:          "  ",
		UseProtoNames:   false,
	}.Marshal(m)
	if err != nil {
		return nil, err
	}
	// Unmarshal the JSON with yaml.v3 so that we can re-marshal it as YAML.
	var doc yaml.Node
	err = yaml.Unmarshal([]byte(s), &doc)
	if err != nil {
		return nil, err
	}
	// The top-level node is a "document" node. We need to marshal the node below it.
	node := doc.Content[0]
	// Restyle the YAML representation so that it will be serialized with YAML defaults.
	StyleForYAML(node)

	return node, nil
}

// StyleForYAML sets the style field on a tree of yaml.Nodes for YAML export.
func StyleForYAML(node *yaml.Node) {
	node.Style = 0
	for _, n := range node.Content {
		StyleForYAML(n)
	}
}

// StyleForJSON sets the style field on a tree of yaml.Nodes for JSON export.
func StyleForJSON(node *yaml.Node) {
	switch node.Kind {
	case yaml.DocumentNode, yaml.SequenceNode, yaml.MappingNode:
		node.Style = yaml.FlowStyle
	case yaml.ScalarNode:
		switch node.Tag {
		case "!!str":
			node.Style = yaml.DoubleQuotedStyle
		default:
			node.Style = 0
		}
	case yaml.AliasNode:
	default:
	}
	for _, n := range node.Content {
		StyleForJSON(n)
	}
}
