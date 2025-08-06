package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mongodb/ftdc"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:     "ftdc-reader [path-to-ftdc-file]",
	Short:   "A CLI tool to parse MongoDB FTDC files.",
	Long:    `A CLI tool to parse MongoDB FTDC files.`,
	Example: "ftdc-reader diagnostics.data/metrics.2023-10-27T10-15-00Z-00000 -o JSON",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		return parseFTDC(filePath, outputFormat)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "outputFormat", "o", "JSON", "Output format. Options: JSON, BSON")
}

type metricOutput struct {
	Key    string  `json:"key" bson:"key"`
	Values []int64 `json:"values" bson:"values"`
}

func parseFTDC(filePath string, format string) error {
	ctx := context.Background()
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	iterator := ftdc.ReadChunks(ctx, file)

	// Use a switch to handle different output formats.
	switch strings.ToUpper(format) {
	case "JSON":
		encoder := json.NewEncoder(os.Stdout)
		for iterator.Next() {
			chunk := iterator.Chunk()
			for _, metric := range chunk.Metrics {
				data := metricOutput{
					Key:    metric.Key(),
					Values: metric.Values,
				}
				if err := encoder.Encode(data); err != nil {
					return fmt.Errorf("failed to encode metric to JSON: %w", err)
				}
			}
		}
	case "BSON":
		for iterator.Next() {
			chunk := iterator.Chunk()
			for _, metric := range chunk.Metrics {
				data := metricOutput{
					Key:    metric.Key(),
					Values: metric.Values,
				}
				bsonBytes, err := bson.Marshal(data)
				if err != nil {
					return fmt.Errorf("failed to marshal metric to BSON: %w", err)
				}
				if _, err := os.Stdout.Write(bsonBytes); err != nil {
					return fmt.Errorf("failed to write BSON to stdout: %w", err)
				}
			}
		}
	default:
		return fmt.Errorf("unsupported output format '%s'. Supported formats are JSON, BSON", format)
	}

	if err := iterator.Err(); err != nil {
		return fmt.Errorf("iterator error during parsing: %w", err)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
