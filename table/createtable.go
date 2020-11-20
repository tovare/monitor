package main

import (
	"context"
	"fmt"
	"log"
	"time"

	uptime "github.com/tovare/monitor"

	"cloud.google.com/go/bigquery"
)

// createTableExplicitSchema demonstrates creating a new BigQuery table and specifying a schema.
func createTableExplicitSchema() error {
	projectID := "homepage-961"
	datasetID := "monitor"
	tableID := "testlog"

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	sampleSchema, err := bigquery.InferSchema(uptime.TestResult{})

	/*	sampleSchema := bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "url", Type: bigquery.StringFieldType},
		{Name: "statuscode", Type: bigquery.IntegerFieldType},
		{Name: "testedtime", Type: bigquery.DateTimeFieldType},
		{Name: "success", Type: bigquery.BooleanFieldType},
		{Name: "duration", Type: bigquery.IntegerFieldType},
		{Name: "durationms", Type: bigquery.IntegerFieldType},
	} */

	metaData := &bigquery.TableMetadata{
		Schema:         sampleSchema,
		ExpirationTime: time.Now().AddDate(2, 0, 0), // Table will be automatically deleted in 1 year.
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("Creating table schema. Use once only.")
	err := createTableExplicitSchema()
	if err != nil {
		log.Fatal(err)
	}
}
