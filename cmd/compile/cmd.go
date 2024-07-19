package compile

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "Compile a Service Configuration for an API",
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
			return nil
		},
	}
	return cmd
}
