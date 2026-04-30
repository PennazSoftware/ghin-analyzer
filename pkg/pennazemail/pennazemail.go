package pennazemail

import (
	"PennazSoftware/ghin-analyzer/pkg/hcutil"
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"gopkg.in/gomail.v2"
)

const (
	EmailVerificationTemplateName string = "TBDVerificationTemplate"
	AwsUsWest2                    string = "us-west-2"
	DefaultSenderEmail            string = "dan@pennaz.com"
)

// AWSSES is the interface for interracting with AWS SES (Simple Email Service)
type AWSSes interface {
	IsEmailVerified(email string) bool
	SendEmail(recipients []string, subject string, htmlBody string, sender string) error
	SendEmailWithAttachments(recipient Recipients, subject string, htmlBody string, fromName string, fromEmail string, attachments []EmailAttachment) error
	SendCustomVerificationEmail(participantEmail string, participantName string, organizerName string) error
	VerifyEmailIdentify(email string) error
}

type Client struct {
	sesClient *ses.Client
	context   context.Context
}

// Verifies the Client struct is implementing all interfaces
var _ AWSSes = (*Client)(nil)

func NewAwsSes() *Client {
	ctx := context.TODO()

	// Load the SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"), // Optional: Specify region explicitly
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	clientService := ses.NewFromConfig(cfg)

	return &Client{
		sesClient: clientService,
		context:   ctx,
	}
}

// IsEmailVerified checks the list of SES verified email addresses to see if the provided email has
// been verified by SES or not.
func (c *Client) IsEmailVerified(email string) bool {
	input := &ses.GetIdentityVerificationAttributesInput{
		Identities: []string{
			email,
		},
	}

	result, err := c.sesClient.GetIdentityVerificationAttributes(c.context, input)
	if err != nil {
		log.Printf("failed to GetIdentityVerificationAttributes for %s. %+v", email, err)
		return false
	}

	if len(result.VerificationAttributes) == 0 {
		return false
	}

	if result.VerificationAttributes[email].VerificationStatus == "Success" {
		return true
	}

	return false
}

// VerifyEmailIdentify adds an email address to the list of identities for your Amazon SES account in the
// current AWS region and attempts to verify it. As a result of executing this operation, a verification
// email is sent to the specified address.
func (c *Client) VerifyEmailIdentify(email string) error {
	input := &ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	}

	_, err := c.sesClient.VerifyEmailIdentity(c.context, input)
	if err != nil {
		log.Printf("failed to VerifyEmailIdentify for %s. %+v", email, err)
	}

	log.Printf("successfully sent email verification mail to %s", email)

	return err
}

// SendCustomVerificationEmail adds an email address to the list of identities for your Amazon SES account
// and sends a custom email verification message to verifiy the email identity.
func (c *Client) SendCustomVerificationEmail(participantEmail string, participantName string, organizerName string) error {
	input := &ses.SendCustomVerificationEmailInput{
		EmailAddress: aws.String(participantEmail),
		TemplateName: aws.String(EmailVerificationTemplateName),
	}

	output, err := c.sesClient.SendCustomVerificationEmail(c.context, input)
	if err != nil {
		log.Printf("failed to SendCustomVerificationEmail for %s. %+v", participantEmail, err)
	}

	log.Printf("Custom Verification Email Result: %+v", output)

	return err
}

// SendEmail composes an email message and immediately queues it for sending.
func (c *Client) SendEmail(recipients []string, subject string, htmlBody string, sender string) error {
	// Compose the input
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			CcAddresses: []string{},
			ToAddresses: recipients,
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// Add the message to the send queue
	_, err := c.sesClient.SendEmail(c.context, input)

	// Display error messages if they occur.
	if err != nil {
		log.Printf("failed to send email - %+v", err)
		return err
	}

	log.Printf("email successfully sent to %+v", recipients)

	return nil
}

func (c *Client) SendEmailWithAttachments(recipient Recipients, subject string, htmlBody string, fromName string, fromEmail string, attachments []EmailAttachment) error {
	// create raw message
	msg := gomail.NewMessage()

	// set to section
	var recipients []string
	for _, r := range recipient.ToEmails {
		recipient := r
		recipients = append(recipients, recipient)
	}

	// Set to emails
	msg.SetHeader("To", recipient.ToEmails...)

	// cc mails mentioned
	if len(recipient.CcEmails) != 0 {
		// Need to add cc mail IDs also in recipient list
		for _, r := range recipient.CcEmails {
			recipient := r
			recipients = append(recipients, recipient)
		}
		msg.SetHeader("cc", recipient.CcEmails...)
	}

	// bcc mails mentioned
	if len(recipient.BccEmails) != 0 {
		// Need to add bcc mail IDs also in recipient list
		for _, r := range recipient.BccEmails {
			recipient := r
			recipients = append(recipients, recipient)
		}
		msg.SetHeader("bcc", recipient.BccEmails...)
	}

	msg.SetAddressHeader("From", fromEmail, fromName)
	msg.SetHeader("To", recipient.ToEmails...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	// If attachments exists
	if len(attachments) != 0 {
		for _, f := range attachments {
			msg.Attach(f.Filename)
		}
	}

	// create a new buffer to add raw data
	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	// create new raw message
	message := types.RawMessage{Data: emailRaw.Bytes()}

	input := &ses.SendRawEmailInput{Source: &fromEmail, Destinations: recipients, RawMessage: &message}

	// Add the message to the send queue
	_, err := c.sesClient.SendRawEmail(c.context, input)

	// Display error messages if they occur.
	if err != nil {
		log.Printf("failed to send email - %+v", err)
		return err
	}

	log.Printf("email successfully sent to %+v", hcutil.ObjectToJSON(recipients))

	return nil
}
