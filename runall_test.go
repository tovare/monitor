package uptime

import (
	"context"
	"testing"
)

func TestRunTests(t *testing.T) {
	err := RunTests(context.Background(), PubSubMessage{[]byte("hello")})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteToDatabase(t *testing.T) {
	// Slight risk of a lost test.
	m, err := ReadDatabase(context.Background(), tests)
	err = WriteToDatabase(context.Background(), m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadToDatabase(t *testing.T) {
	m, err := ReadDatabase(context.Background(), tests)
	if err != nil {
		t.Fatal(err)
	}
	for k := range m {
		if m[k].URL != tests[k].URL {
			t.Fatalf("Failure on %v", k)
		}
	}
}

func TestGetPasswordFromSecrets(t *testing.T) {
	pass, err := GetPasswordFromSecrets(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(pass) < 5 {
		t.Fatal("Password to short")
	}
}

func TestSendAlertEmail(t *testing.T) {
	err := SendAlertEmail(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
