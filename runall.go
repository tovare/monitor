package uptime

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Test Definitions
type TestResult struct {
	URL        string        `firestore:"url,omitempty"`
	StatusCode int           `firestore:"lastresult,omitempty"`
	Tested     time.Time     `firestore:"testedtime,omitempty"`
	Success    bool          `firestore:"success"`
	Duration   time.Duration `firestore:"duration,omitempty"`
	DurationMS int64         `firestore:"durationms,omitempty"`
}

type TestMap map[string]TestResult

// tests i prepopulated with all tests. Data will be overwritten by contents in the database.
// if any changes are done to tests it will only new tests will be added to the databae.
var tests = TestMap{
	"tovarecom": {
		URL:        "https://www.tovare.com/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"tovarecom-hybrids": {
		URL:        "https://tovare.com/2020/hybrids-start/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"tovarecom-dashboard": {
		URL:        "https://tovare.com/dashboard/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"alleyoop": {
		URL:        "https://alleyoop.no/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"navno": {
		URL:        "https://nav.no/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"arbeidsplassen": {
		URL:        "https://arbeidsplassen.nav.no/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
}

// RunTest runs all the tests and compares to prior test. If
// the state of a service changes from prior test an email is
// sent.
//
// gcloud functions deploy RunTests --memory=128 --runtime go113 --trigger-topic monitor
//
// The trigger is sent every 5 minutes using the Clooud Trigger feature
// in Google Cloud.
//
func RunTests(ctx context.Context, m PubSubMessage) (err error) {

	// Run all tests
	for i := range tests {
		start := time.Now()
		t := TestURL(tests[i])
		t.Duration = time.Since(start)
		t.DurationMS = int64(t.Duration / 1000000)
		t.Tested = start
		tests[i] = t
	}
	// After completing all the tests we need to figure out
	// if any states has changed.
	oldTests, err := ReadDatabase(ctx, tests)
	changes := make(TestMap)
	for k := range tests {
		if _, ok := oldTests[k]; ok {
			// This is not the first time a test has been runned.
			if oldTests[k].Success != tests[k].Success {
				// The prior result has changed. add to list
				// of tests which requre reporting.
				changes[k] = tests[k]
			}
		}
	}
	// Write the latest test to the database.
	err = WriteToDatabase(ctx, tests)
	if err != nil {
		return
	}

	if len(changes) > 0 {
		// Report on the list of changes by sending an
		// email for each of them.
		err = SendAlertEmail(ctx)
		if err != nil {
			return
		}
	}
	return
}

// TestURL runs a single test and returns the test-results.
func TestURL(test TestResult) TestResult {

	resp, err := http.Get(test.URL)
	if err != nil {
		if resp != nil {
			test.StatusCode = resp.StatusCode
		} else {
			test.StatusCode = -1
		}
		test.Success = false
		return test
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		test.Success = true
	} else {
		test.StatusCode = resp.StatusCode
		test.Success = false
	}
	_, err = ioutil.ReadAll(resp.Body)

	// TODO: Extend with content check.
	return test
}

// WriteToDatabase adds all the tests to the database.
// this function will replace existing entries with the same name.
func WriteToDatabase(ctx context.Context, tests TestMap) (err error) {

	client, err := firestore.NewClient(ctx, "homepage-961")
	if err != nil {
		return
	}
	defer client.Close()

	for k, v := range tests {
		_, err = client.Collection("monitoring").Doc(k).Set(ctx, v)
		if err != nil {
			return
		}
	}

	return
}

// ReadDatabase into a map structure whos entries are still listed in the
// active tests.
func ReadDatabase(ctx context.Context, tests TestMap) (m TestMap, err error) {

	m = make(TestMap)

	client, err := firestore.NewClient(ctx, "homepage-961")
	if err != nil {
		return
	}
	defer client.Close()

	for k := range tests {
		var doc *firestore.DocumentSnapshot
		doc, err = client.Collection("monitoring").Doc(k).Get(ctx)
		if err != nil {
			return
		}
		if doc.Exists() {
			var res TestResult
			// Add a new document.
			err = doc.DataTo(&res)
			if err != nil {
				return
			}
			m[k] = res
		}
	}
	return
}
