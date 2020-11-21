package uptime

import "time"

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
