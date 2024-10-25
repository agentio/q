package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	logging "cloud.google.com/go/logging/apiv2"
	"cloud.google.com/go/logging/apiv2/loggingpb"
	"google.golang.org/protobuf/proto"

	appengine "google.golang.org/genproto/googleapis/appengine/logging/v1"

	"github.com/spf13/cobra"
)

func readCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "read PROJECT SERVICE",
		Short: "read log entries with the Cloud Logging API",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := cmd.Context()
			c, err := logging.NewClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()

			project := args[0]
			service := url.PathEscape(args[1])

			iter := c.ListLogEntries(ctx, &loggingpb.ListLogEntriesRequest{
				ResourceNames: []string{"projects/" + project},
				Filter:        `logName = "projects/` + project + `/logs/` + service + `"`,
				OrderBy:       "timestamp desc",
			})
			count := 0
			for {
				entry, err := iter.Next()
				if err != nil {
					log.Printf("%s", err)
					break
				}
				b, err := json.MarshalIndent(entry, "", "  ")
				if err != nil {
					log.Printf("%s", err)
					break
				}
				fmt.Printf("%s\n", string(b))
				switch v := entry.Payload.(type) {
				case *loggingpb.LogEntry_ProtoPayload:
					{
						if v.ProtoPayload.GetTypeUrl() == "type.googleapis.com/google.appengine.logging.v1.RequestLog" {
							var payload appengine.RequestLog
							err = proto.Unmarshal(v.ProtoPayload.GetValue(), &payload)
							if err != nil {
								log.Printf("%s", err)
							}
							b, err := json.MarshalIndent(&payload, "", "  ")
							if err != nil {
								log.Printf("%s", err)
								break
							}
							fmt.Printf("%s\n", string(b))
						}
					}
				default:

				}
				count += 1
				if count == limit {
					break
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 100, "maximum number of entries to return")
	return cmd
}
