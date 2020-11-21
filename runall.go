package uptime

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// tests i prepopulated with all tests. Data will be overwritten by contents in the database.
// if any changes are done to tests it will only new tests will be added to the databae.
var tests = TestMap{
	"tovarecom": {
		Name:       "tovarecom",
		URL:        "https://tovare.com/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"tovarecom-hybrids": {
		Name:       "tovarecom-hybrids",
		URL:        "https://tovare.com/2020/hybrids-start/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"tovarecom-dashboard": {
		Name:       "tovarecom-dashboard",
		URL:        "https://tovare.com/dashboard/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"alleyoop": {
		Name:       "alleyoop",
		URL:        "https://alleyoop.no/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"navno": {
		Name:       "navno",
		URL:        "https://www.nav.no/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"arbeidsplassen": {
		Name:       "arbeidsplassen",
		URL:        "https://arbeidsplassen.nav.no/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"idebanken": {
		Name:       "idebanken",
		URL:        "http://idebanken.no",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"navnofamilie": {
		Name:       "navnofamilie",
		URL:        "http://familie.nav.no",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"data.nav.no": {
		Name:       "data.nav.no",
		URL:        "http://data.nav.no",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"detsombetyrnoe": {
		Name:       "detsombetyrnoe",
		URL:        "http://detsombetyrnoe.no",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"deterdinpensjon": {
		Name:       "deterdinpensjon",
		URL:        "https://www.deterdinpensjon.no",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"nais.io": {
		Name:       "nais.io",
		URL:        "https://nais.io/",
		StatusCode: 200,
		Tested:     time.Now(),
		Success:    true,
		Duration:   0,
	},
	"memu.no": {
		Name:       "memu.no",
		URL:        "https://memu.no",
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
	// Report on the list of changes by sending an
	// email for each of them.

	if len(changes) > 0 {
		err = SendAlertEmail(ctx)
		if err != nil {
			return
		}
	}

	// Stream data to BigQuery
	err = StreamToBigQuery(ctx, tests)
	if err != nil {
		return err
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
		test.ErrorMsg = err.Error()
		return test
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		test.Success = true
	} else {
		test.StatusCode = resp.StatusCode
		test.Success = false
		test.ErrorMsg = http.StatusText(resp.StatusCode)
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

func StreamToBigQuery(ctx context.Context, tests TestMap) error {
	projectID := "homepage-961"
	datasetID := "monitor"
	tableID := "uptime"

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	inserter := client.Dataset(datasetID).Table(tableID).Inserter()

	testresults := make([]*TestResult, 0)
	for _, v := range tests {
		x := v
		testresults = append(testresults, &x)
	}
	if err := inserter.Put(ctx, testresults); err != nil {
		return err
	}
	return nil
}
