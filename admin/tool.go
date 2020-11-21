// This is used oncly once.
package main

import (
	"context"
	"fmt"
	"time"

	uptime "github.com/tovare/monitor"

	"cloud.google.com/go/bigquery"
)

// createDataset demonstrates creation of a new dataset using an explicit destination location.
// This is used once only.
func createDataset(projectID, datasetID string) error {

	fmt.Print("Creating BigQuery dataset")

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

// createTableExplicitSchema demonstrates creating a new BigQuery table and specifying a schema.
func createTableExplicitSchema(projectID, datasetID, tableID string) error {

	fmt.Println("Creating table schema. Use once only.")

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	sampleSchema, err := bigquery.InferSchema(uptime.TestResult{})

	metaData := &bigquery.TableMetadata{
		Schema:         sampleSchema,
		ExpirationTime: time.Now().AddDate(2, 0, 0), // Table will be automatically deleted in 2 years.
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	return nil
}

func main() {
	// createDataset("homepage-961", "monitor")
	createTableExplicitSchema("homepage-961", "monitor", "uptime")

}
