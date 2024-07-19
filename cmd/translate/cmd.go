package translate

import (
	"fmt"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func Cmd() *cobra.Command {
	var source string
	var target string
	var format string
	var parent string
	cmd := &cobra.Command{
		Use:   "translate TEXT",
		Short: "Translate with the Google Cloud Translation API",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := translate.NewTranslationClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()
			response, err := c.TranslateText(ctx, &translatepb.TranslateTextRequest{
				SourceLanguageCode: source,
				TargetLanguageCode: target,
				Contents:           args,
				Parent:             parent,
			})
			if err != nil {
				return err
			}
			if format == "json" {
				b, err := protojson.Marshal(response)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(b))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&source, "source", "en-us", "source language")
	cmd.Flags().StringVar(&target, "target", "es-mx", "target language")
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	cmd.Flags().StringVarP(&parent, "parent", "p", "", "parent project (format: projects/PROJECTID)")
	return cmd
}
