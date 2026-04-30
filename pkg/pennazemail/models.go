package pennazemail

import "fmt"

// AmazonSqsNotification - Represents the bounce or complaint notification stored in Amazon SQS.
type AmazonSqsNotification struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// AmazonSesComplaintNotification represents an Amazon SES complaint notification.
type AmazonSesNotification struct {
	// NotificationType - A string that holds the type of notification represented by the JSON object.
	// Possible values are Bounce, Complaint, or Delivery.
	NotificationType string `json:"notificationType"`
	// Mail - A JSON object that contains information about the original mail to which the notification pertains.
	Mail string `json:"mail"`
	// Bounce - This field is present only if the notificationType is Bounce and contains a JSON object
	// that holds information about the bounce.
	Bounce AmazonSesBounce `json:"bounce"`
	// Complaint - This field is present only if the notificationType is Complaint and contains a JSON
	// object that holds information about the complaint.
	Complaint AmazonSesComplaint `json:"amazonSesComplaint"`
	// Delivery - This field is present only if the notificationType is Delivery and contains a JSON
	// object that holds information about the delivery. For more information
	Delivery AmazonSesDelivery `json:"delivery"`
}

// AmazonSesComplaint represents meta data for the complaint notification from Amazon SES.
type AmazonSesComplaint struct {
	// ComplainedRecipients is A list that contains information about recipients that may have
	// been responsible for the complaint.
	ComplainedRecipients []AmazonSesRecipient `json:"complainedRecipients"`
	// Timestamp The date and time when the ISP sent the complaint notification, in ISO 8601 format.
	// The date and time in this field might not be the same as the date and time when Amazon SES
	// received the notification.
	Timestamp string `json:"timestamp"`
	MessageID string `json:"messageID"`
	// FeedbackID is A unique ID associated with the complaint.
	FeedbackID string `json:"feedbackID"`
	// ComplaintSubType field can either be null or OnAccountSuppressionList. If the value is
	// OnAccountSuppressionList, Amazon SES accepted the message, but didn't attempt to send it
	// because it was on the account-level suppression list.
	ComplaintSubType string `json:"complaintSubType"`
	// UserAgent - The value of the User-Agent field from the feedback report. This indicates
	// the name and version of the system that generated the report.  This field is only available
	// if the feedback report is attached to the complaint
	UserAgent string `json:"userAgent"`
	// ComplaintFeedbackType - The value of the Feedback-Type field from the feedback report
	// received from the ISP. This contains the type of feedback. This field is only available
	// if the feedback report is attached to the complaint
	ComplaintFeedbackType string `json:"complaintFeedbackType"`
	// Arrival Date - The value of the Arrival-Date or Received-Date field from the feedback
	// report (in ISO8601 format). This field may be absent in the report (and therefore also
	// absent in the JSON). This field is only available if the feedback report is attached to
	// the complaint
	ArrivalDate string `json:"arrivalDate"`
}

// AmazonSesBounce - Represents meta data for the bounce notification from Amazon SES.
type AmazonSesBounce struct {
	// BounceType - The type of bounce, as determined by Amazon SES.
	BounceType string `json:"bounceType"`
	// BounceSubType - The subtype of the bounce, as determined by Amazon SES.
	BounceSubType string `json:"bounceSubType"`
	// Timestamp - The date and time at which the bounce was sent (in ISO8601 format).
	// Note that this is the time at which the notification was sent by the ISP, and not
	// the time at which it was received by Amazon SES.
	Timestamp string `json:"timestamp"`
	// BouncedRecipienets - A list that contains information about the recipients of the
	// original mail that bounced.
	BouncedRecipients []AmazonSesRecipient `json:"bouncedRecipients"`
	// FeedbackID is A unique ID associated with the complaint.
	FeedbackID string `json:"feedbackID"`
	// RemoteMtaIp - The IP address of the MTA to which Amazon SES attempted to deliver the email.
	// This field is only available if Amazon SES was able to contact the remote Message Transfer
	// Authority (MTA)
	RemoteMtaIp string `json:"remoteMtaIp"`
	// ReportingMTA - The value of the Reporting-MTA field from the DSN. This is the value of the
	// MTA that attempted to perform the delivery, relay, or gateway operation described in the DSN.
	// This field is only available if a delivery status notification (DSN) was attached to the bounce.
	ReportingMTA string `json:"reportingMTA"`
}

// AmazonSesRecipient - Represents the email address of recipients that bounced
// / when sending from Amazon SES.
type AmazonSesRecipient struct {
	// EmailAddress - The email address of the recipient. If a DSN is available, this is
	// the value of the Final-Recipient field from the DSN.
	EmailAddress string `json:"emailAddress"`
	// Action - The value of the Action field from the DSN. This indicates the action
	// performed by the Reporting-MTA as a result of its attempt to deliver the message
	// to this recipient.
	Action string `json:"action"`
	// Status - The value of the Status field from the DSN. This is the per-recipient
	// transport-independent status code that indicates the delivery status of the message.
	Status string `json:"status"`
	// DiagnosticCode - The status code issued by the reporting MTA. This is the value of the
	// Diagnostic-Code field from the DSN. This field may be absent in the DSN (and therefore
	// also absent in the JSON).
	DiagnosticCode string `json:"diagnosticCode"`
}

// AmazonSesDelivery - Represents meta data for the bounce notification from Amazon SES.
type AmazonSesDelivery struct {
	// Timestamp - The date and time at which the bounce was sent (in ISO8601 format).
	// Note that this is the time at which the notification was sent by the ISP, and not
	// the time at which it was received by Amazon SES.
	Timestamp string `json:"timestamp"`
	// ProcessingTimeMillis - The time in milliseconds between when Amazon SES accepted the request
	// from the sender to passing the message to the recipient's mail server.
	ProcessingTimeMillis int `json:"processingTimeMillis"`
	// Recipients - A list of the intended recipients of the email to which the delivery notification applies.
	Recipients []string `json:"recipients"`
	// SmtpResponse - The SMTP response message of the remote ISP that accepted the email from Amazon
	// SES. This message varies by email, by receiving mail server, and by receiving ISP.
	SmtpResponse string `json:"smtpResponse"`
	// ReportingMta -  The hostname of the Amazon SES mail server that sent the mail.
	ReportingMta string `json:"reportingMta"`
	// RemoteMtaIp - The IP address of the MTA to which Amazon SES delivered the email.
	RemoteMtaIp string `json:"remoteMtaIp"`
}

func GetBounceDescription(bounceType string, bounceSubType string) string {
	if bounceType == "Undetermined" && bounceSubType == "Undetermined" {
		return `The recipient's email provider sent a bounce message. The bounce message didn't contain enough information for Amazon SES to determine the reason for the bounce. The bounce email, which was sent to the address in the Return-Path header of the email that resulted in the bounce, might contain additional information about the issue that caused the email to bounce.`
	}

	if bounceType == "Permanent" && bounceSubType == "General" {
		return `The recipient's email provider sent a hard bounce message, but didn't specify the reason for the hard bounce. Important: When you receive this type of bounce notification, you should immediately remove the recipient's email address from your mailing list. Sending messages to addresses that produce hard bounces can have a negative impact on your reputation as a sender. If you continue sending email to addresses that produce hard bounces, we might pause your ability to send additional email.`
	}

	if bounceType == "Permanent" && bounceSubType == "NoEmail" {
		return `The intended recipient's email provider sent a bounce message indicating that the email address doesn't exist. Important: When you receive this type of bounce notification, you should immediately remove the recipient's email address from your mailing list. Sending messages to addresses that don't exist can have a negative impact on your reputation as a sender. If you continue sending email to addresses that don't exist, we might pause your ability to send additional email.`
	}

	if bounceType == "Permanent" && bounceSubType == "Suppressed" {
		return `The recipient's email address is on the Amazon SES suppression list because it has a recent history of producing hard bounces. To override the global suppression list, see 'Using the Amazon SES account-level suppression list'.`
	}

	if bounceType == "Permanent" && bounceSubType == "OnAccountSuppressionList" {
		return `Amazon SES has suppressed sending to this address because it is on the account-level suppression list. This does not count toward your bounce rate metric.`
	}

	if bounceType == "Transient" && bounceSubType == "General" {
		return `The recipient's email provider sent a general bounce message. You might be able to send a message to the same recipient in the future if the issue that caused the message to bounce is resolved. Note: If you send an email to a recipient who has an active automatic response rule (such as an "out of the office" message), you might receive this type of notification. Even though the response has a notification type of Bounce, Amazon SES doesn't count automatic responses when it calculates the bounce rate for your account.`
	}

	if bounceType == "Transient" && bounceSubType == "MailboxFull" {
		return `The recipient's email provider sent a bounce message because the recipient's inbox was full. You might be able to send to the same recipient in the future when the mailbox is no longer full.`
	}

	if bounceType == "Transient" && bounceSubType == "MessageTooLarge" {
		return `The recipient's email provider sent a bounce message because message you sent was too large. You might be able to send a message to the same recipient if you reduce the size of the message.`
	}

	if bounceType == "Transient" && bounceSubType == "ContentRejected" {
		return `The recipient's email provider sent a bounce message because the message you sent contains content that the provider doesn't allow. You might be able to send a message to the same recipient if you change the content of the message.`
	}

	if bounceType == "Transient" && bounceSubType == "AttachmentRejected" {
		return `The recipient's email provider sent a bounce message because the message contained an unacceptable attachment. For example, some email providers may reject messages with attachments of a certain file type, or messages with very large attachments. You might be able to send a message to the same recipient if you remove or change the content of the attachment.`
	}

	return fmt.Sprintf("Could not determine a description of the bounce for bounceType (%s) and bounceSubType (%s).", bounceType, bounceSubType)
}

// EmailAttachment contains details for attaching a file to an email message
type EmailAttachment struct {
	Filename string
}

// Recipients specifies the different types of recipients of an email message
type Recipients struct {
	ToEmails  []string
	CcEmails  []string
	BccEmails []string
}
