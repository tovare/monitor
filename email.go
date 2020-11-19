package uptime

import (
	"context"
	"net/smtp"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// SendAlertEmail when the status of a service changes.
// I needed to allow insecure apps to make this happen.
func SendAlertEmail(ctx context.Context) (err error) {

	password, err := GetPasswordFromSecrets(ctx)
	if err != nil {
		return
	}

	msg := `To: mail@tovare.com
Subject: Alert notification

This is the message

̀̀̀ `

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", "mail@tovare.com", password, "smtp.gmail.com"),
		"mail@tovare.com", []string{"mail@tovare.com"}, []byte(msg))
	if err != nil {
		return
	}
	return
}

// GetPasswordFromSecrets returns my gmail password.
// We only need the secret when sending emails, which is rare.
func GetPasswordFromSecrets(ctx context.Context) (string, error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/908565461144/secrets/gmail-password/versions/1",
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), err
}
