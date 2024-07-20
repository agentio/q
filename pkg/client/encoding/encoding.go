package encoding

import (
	"bytes"
	"fmt"
	"io"

	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

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

// styleForJSON sets the style field on a tree of yaml.Nodes for JSON export.
func styleForJSON(node *yaml.Node) {
	// Strip comments, they confuse downstream json-to-proto deserialization.
	node.HeadComment = ""
	node.LineComment = ""
	node.FootComment = ""
	// Adjust style.
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
		styleForJSON(n)
	}
}

// extractType gets the "type" string value from the top of a tree of YAML nodes.
// To avoid confusing downstream deserializers, it removes "type" and its value
// from the tree.
func extractType(node *yaml.Node) string {
	switch node.Kind {
	case yaml.DocumentNode:
		return extractType(node.Content[0])
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == "type" {
				// Remove these two entries and return the type value.
				value := node.Content[i+1].Value
				node.Content = append(node.Content[0:i], node.Content[i+2:]...)
				return value
			}
		}
	default:
		return ""
	}
	return ""
}

func ParseYaml(bytes []byte) (proto.Message, error) {
	var node yaml.Node
	err := yaml.Unmarshal(bytes, &node)
	if err != nil {
		return nil, err
	}
	typeString := extractType(&node)
	styleForJSON(&node)
	b, err := yaml.Marshal(node.Content[0])
	if err != nil {
		return nil, err
	}
	// Handle known types.
	if typeString == "google.api.Service" {
		var service serviceconfig.Service
		err = protojson.Unmarshal(b, &service)
		if err != nil {
			return nil, err
		}
		return &service, nil
	}
	return nil, fmt.Errorf("unsupported type %s", typeString)
}
