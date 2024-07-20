package compile

import (
	"fmt"
	"log"
	"os"

	"github.com/agentio/q/pkg/encoding"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile SERVICE DESCRIPTOR",
		Short: "Compile a Service Configuration for an API",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			{
				bytes, err := os.ReadFile(args[0])
				if err != nil {
					return err
				}

				message, err := encoding.UnmarshalYaml(bytes)
				if err != nil {
					return err
				}

				fmt.Fprintf(cmd.OutOrStdout(), "%+v\n", message)
			}

			{
				bytes, err := os.ReadFile(args[1])
				if err != nil {
					return err
				}
				var descriptors descriptorpb.FileDescriptorSet
				if err := proto.Unmarshal(bytes, &descriptors); err != nil {
					log.Fatalln("Failed to parse descriptors:", err)
				}
			}
			return nil
		},
	}
	return cmd
}
