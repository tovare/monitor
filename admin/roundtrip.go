package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

var projectID = "homepage-961"

type TestSchema struct {
	MyTime  time.Time `bigquery:"mytime"`
	MyValue int       `bigquery:"myvalue"`
}

// Roundtrip contains a test of bigquery attempting to do a full roundtrip
// using structs.

func Roundtrip() {

	fmt.Println("Starting the roundtrip test.")
	fmt.Println("Creating the client.")

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("bigquery.NewClient: %v", err)
		return
	}
	defer client.Close()

	// The dataset
	fmt.Println("Creating a dataset if it dosn´t exist allready.")
	if meta, err := client.Dataset("roundtrip").Metadata(ctx); err != nil {
		fmt.Println("THe dataset doesn´t exist so we try to make it.")
		meta := &bigquery.DatasetMetadata{
			Description: "This is a tests for examining bigquery functionality.",
			Location:    "EU", // See https://cloud.google.com/bigquery/docs/locations
		}
		if err := client.Dataset("roundtrip").Create(ctx, meta); err != nil {
			fmt.Print("Failed to create roundtrip: ", err)
			return
		}
		fmt.Println("dataset created.")
	} else {
		fmt.Println("This is the metadata: ", meta)
	}
	// The table
	fmt.Println("Creating a table, but deleting the old one if it exists.")

	schema, err := bigquery.InferSchema(TestSchema{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("The result of the schema interpretation:")
	fmt.Println("--------------------------------")
	for _, fs := range schema {
		fmt.Println(fs.Name, fs.Type, fs.Required)
	}

	metaData := &bigquery.TableMetadata{
		Schema:         schema,
		ExpirationTime: time.Now().AddDate(2, 0, 0), // Table will be automatically deleted in 2 years.
	}
	tableRef := client.Dataset("roundtrip").Table("test")
	// Delete the table if it exists.
	_, err = tableRef.Metadata(ctx)
	if err == nil {
		fmt.Println("deleting old table")
		if err := tableRef.Delete(ctx); err != nil {
			fmt.Println("deleting old table failed:", err)
			return
		}
	}

	if err := tableRef.Create(ctx, metaData); err != nil {
		fmt.Println("Table creation failed: ", err)
		return
	}
	fmt.Println("Table created")

	items2 := []TestSchema{
		{MyTime: time.Now().AddDate(0, 0, 0), MyValue: 10},
		{MyTime: time.Now().AddDate(0, 0, 1), MyValue: 11},
		{MyTime: time.Now().AddDate(0, 0, 2), MyValue: 12},
		{MyTime: time.Now().AddDate(0, 0, 3), MyValue: 13},
	}

	if err := tableRef.Inserter().Put(ctx, items2); err != nil {
		fmt.Println("Failed to insert values", err)
	}

	// Now lets try to roundtrip this stuff.
	q := client.Query(
		`SELECT * FROM homepage-961.roundtrip.test WHERE (mytime) IN 
		   ( SELECT MAX(mytime) FROM homepage-961.roundtrip.test )`)
	it, err := q.Read(ctx)
	if err != nil {
		fmt.Println("Failed to run SQL: ", err)
		return
	}
	for {
		var values TestSchema
		err := it.Next(&values) // Will zero out values.
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Failed in loop: ", err)
			return
		}
		//fmt.Println(" Length of result :", len(values))
		fmt.Println(" First entry time :", values.MyTime)
		fmt.Println("First entry value :", values.MyValue)
	}
}
