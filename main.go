package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mongodb/ftdc"
)

func parseFTDC(filePath string) error {
	ctx := context.Background()
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()
	iterator := ftdc.ReadChunks(ctx, file)

	encoder := json.NewEncoder(os.Stdout)

	for iterator.Next() {
		chunk := iterator.Chunk()

		for _, metric := range chunk.Metrics {
			metricData := struct {
				Key    string  `json:"key"`
				Values []int64 `json:"values"`
			}{
				Key:    metric.Key(),
				Values: metric.Values,
			}

			if err := encoder.Encode(metricData); err != nil {
				return fmt.Errorf("failed to encode metric to JSON: %w", err)
			}
		}
	}
	if err := iterator.Err(); err != nil {
		return fmt.Errorf("iterator error during parsing: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		// Usage information should go to stderr.
		fmt.Fprintln(os.Stderr, "Usage: go run main.go <path-to-ftdc-file>")
		fmt.Fprintln(os.Stderr, "Example: go run main.go diagnostics.data/metrics.2023-10-27T10-15-00Z-00000")
		os.Exit(1)
	}

	filePath := os.Args[1]

	if err := parseFTDC(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
