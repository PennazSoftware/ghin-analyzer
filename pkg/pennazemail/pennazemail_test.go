//go:build integration
// +build integration

// To run the integration tests, run the following command:  go test -tags=integration
package pennazemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendRawEmail_INT(t *testing.T) {
	sesClient := NewAwsSes()

	body := "<html><body>This is my test email. It should have an attachment!</body></html>"

	recipient := Recipients{
		ToEmails: []string{"danpenn@msn.com"},
		CcEmails: []string{"dan@pennaz.com"},
	}

	err := sesClient.SendEmailWithAttachments(recipient, "Test Email", body, "Dan Penn", "dan@pennaz.com", []EmailAttachment{})
	assert.NoError(t, err)
}
