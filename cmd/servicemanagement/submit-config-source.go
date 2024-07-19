package servicemanagement

import (
	"fmt"
	"os"

	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"cloud.google.com/go/servicemanagement/apiv1/servicemanagementpb"
	"github.com/spf13/cobra"
)

func submitConfigSourceCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "submit-config-source SERVICE FILE...",
		Short: "Submit config source",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()

			files := []*servicemanagementpb.ConfigFile{}
			for i, f := range args {
				if i == 0 {
					continue
				}
				b, err := os.ReadFile(f)
				if err != nil {
					return err
				}
				// TODO: replace these assumptions with something smarter
				fileType := servicemanagementpb.ConfigFile_SERVICE_CONFIG_YAML
				if i > 1 {
					fileType = servicemanagementpb.ConfigFile_FILE_DESCRIPTOR_SET_PROTO
				}
				files = append(files, &servicemanagementpb.ConfigFile{
					FilePath:     f,
					FileContents: b,
					FileType:     fileType,
				})
			}
			operation, err := c.SubmitConfigSource(ctx, &servicemanagementpb.SubmitConfigSourceRequest{
				ServiceName: args[0],
				ConfigSource: &servicemanagementpb.ConfigSource{
					Files: files,
				},
				ValidateOnly: true,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", operation.Name())
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
