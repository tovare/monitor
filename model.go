package uptime

import (
	"time"

	"cloud.google.com/go/bigquery"
)

// Test Definitions
type TestResult struct {
	Name       string        `firestore:"name,omitempty" bigquery:"name"`
	URL        string        `firestore:"url,omitempty" bigquery:"url,nullable"`
	StatusCode int           `firestore:"statuscode,omitempty" bigquery:"statuscode,nullable"`
	Tested     time.Time     `firestore:"testedtime,omitempty" bigquery:"testedtime,nullable"`
	Success    bool          `firestore:"success" bigquery:"success,nullable"`
	Duration   time.Duration `firestore:"duration,omitempty" bigquery:"duration,nullable"`
	DurationMS int64         `firestore:"durationms,omitempty" bigquery:"durationms,nullable"`
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
	}, bigquery.NoDedupeID, nil
}
