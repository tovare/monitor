// This is used oncly once.
package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createDataset demonstrates creation of a new dataset using an explicit destination location.
func createDataset(projectID, datasetID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		Location: "EU", // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Print("Creating BigQuery dataset")
	createDataset("homepage-961", "monitor")
}
