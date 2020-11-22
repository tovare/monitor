// This is used oncly once.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	uptime "github.com/tovare/monitor"
	"google.golang.org/api/iterator"

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

// createTestResultTable create a new table of TestResults. The procedure
// will fail if the table allready exists.
func createTestResultTable(projectID, datasetID, tableID string) error {

	fmt.Println("Creating table schema. Use once only.")

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	sampleSchema, err := bigquery.InferSchema(uptime.TestResult{
		Name:       "tovarecom",
		URL:        "https://tovare.com/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
		ErrorMsg:   "",
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Print out
	fmt.Println("The result of the schema interpretation:")
	fmt.Println("--------------------------------")
	for _, fs := range sampleSchema {
		fmt.Println(fs.Name, fs.Type, fs.Required)
	}

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

// printDatasetInfo demonstrates fetching dataset metadata and printing some of it to an io.Writer.
func printDatasetInfo(w io.Writer, projectID, datasetID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	meta, err := client.Dataset(datasetID).Metadata(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Dataset ID: %s\n", datasetID)
	fmt.Fprintf(w, "Description: %s\n", meta.Description)
	fmt.Fprintln(w, "Labels:")
	for k, v := range meta.Labels {
		fmt.Fprintf(w, "\t%s: %s", k, v)
	}
	fmt.Fprintln(w, "Tables:")
	it := client.Dataset(datasetID).Tables(ctx)

	cnt := 0
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		cnt++
		fmt.Fprintf(w, "\t%s\n", t.TableID)
	}
	if cnt == 0 {
		fmt.Fprintln(w, "\tThis dataset does not contain any tables.")
	}
	return nil
}

// findLastRecord finds the last record written to the bigquery table. This can be used
// for incremental updates to the database.
func findLastRecord(projectID, datasetID, tableID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	q := client.Query(`
		SELECT max(testedtime) FROM homepage-961.monitor.uptime `)
	it, err := q.Read(ctx)
	if err != nil {
		return err
	}
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(values[0])
	}
	return nil
}

func main() {
	// createDataset("homepage-961", "monitor")
	// err := createTestResultTable("homepage-961", "monitor", "uptime")
	// err := printDatasetInfo(os.Stdout, "homepage-961", "monitor")
	err := findLastRecord("homepage-961", "monitor", "uptime")
	if err != nil {
		log.Fatal(err)
	}

}
