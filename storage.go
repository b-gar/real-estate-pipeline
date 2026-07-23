package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
)

func StreamToBigQuery[T any](ctx context.Context, client *bigquery.Client, datasetID, tableID string, data []T, structType any) error {
	table := client.Dataset(datasetID).Table(tableID)
	if _, err := table.Metadata(ctx); err != nil {
		schema, err := bigquery.InferSchema(structType)
		if err != nil {
			return fmt.Errorf("Failed to infer schema: %w", err)
		}
		if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
			return err
		}
		log.Printf("Created new BigQuery table: %s.%s", datasetID, tableID)
	}
	inserter := table.Inserter()
	if err := inserter.Put(ctx, data); err != nil {
		return err
	}
	log.Printf("Successfully streamed %d records to %s.%s", len(data), datasetID, tableID)
	return nil
}
