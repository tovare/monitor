package uptime

import "time"

// Test Definitions
type TestResult struct {
	Name       string        `firestore:"name,omitempty"`
	URL        string        `firestore:"url,omitempty"`
	StatusCode int           `firestore:"statuscode,omitempty"`
	Tested     time.Time     `firestore:"testedtime,omitempty"`
	Success    bool          `firestore:"success"`
	Duration   time.Duration `firestore:"duration,omitempty"`
	DurationMS int64         `firestore:"durationms,omitempty"`
}

type TestMap map[string]TestResult
