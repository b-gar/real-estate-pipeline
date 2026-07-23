package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

func StreamToBigQuery[T any](ctx context.Context, projectID, datasetID, tableID string, data []T, structType any) error {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Failed to create BigQuery client: %w", err)
	}
	defer client.Close()
	table := client.Dataset(datasetID).Table(tableID)
	if _, err := table.Metadata(ctx); err != nil {
		schema, err := bigquery.InferSchema(structType)
		if err != nil {
			return fmt.Errorf("Failed to infer schema: %w", err)
		}
		if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
			return fmt.Errorf("Failed to create table: %w", err)
		}
		fmt.Printf("Created new BigQuery table: %s.%s\n", datasetID, tableID)
	}
	inserter := table.Inserter()
	if err := inserter.Put(ctx, data); err != nil {
		return fmt.Errorf("Failed to upload data to BigQuery: %w", err)
	}
	fmt.Printf("Successfully uploaded %d records to %s.%s\n", len(data), datasetID, tableID)
	return nil
}
