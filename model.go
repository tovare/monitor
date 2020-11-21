package uptime

import (
	"time"

	"cloud.google.com/go/bigquery"
)

// Test Definitions
type TestResult struct {
	Name       string        `firestore:"name" bigquery:"name"`
	URL        string        `firestore:"url" bigquery:"url"`
	StatusCode int           `firestore:"statuscode" bigquery:"statuscode"`
	Tested     time.Time     `firestore:"testedtime" bigquery:"testedtime"`
	Success    bool          `firestore:"success" bigquery:"success"`
	Duration   time.Duration `firestore:"duration" bigquery:"duration"`
	DurationMS int64         `firestore:"durationms" bigquery:"durationms"`
	ErrorMsg   string        `firestore:"errormsg" bigquery:"errormsg"`
	TestString string        `firestore:"teststring" bigquery:"-"`
}

type TestMap map[string]TestResult

// Save implements the BigQuery ValueSaver interface and uses a best effort
// de-duplicator. The list below needs to be in sync with the above since I
// opted not to map at runtime through introspection.
func (i *TestResult) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"name":       i.Name,
		"url":        i.URL,
		"statuscode": i.StatusCode,
		"testedtime": i.Tested,
		"success":    i.Success,
		"duration":   i.Duration,
		"durationms": i.DurationMS,
		"errormsg":   i.ErrorMsg,
	}, bigquery.NoDedupeID, nil
}
