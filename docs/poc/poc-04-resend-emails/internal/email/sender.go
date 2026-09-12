package email

import (
	"fmt"

	"github.com/resend/resend-go/v4"
)

func Send(client *resend.Client, sender, receiver, subject, htmlBody, reply string) error {
	params := &resend.SendEmailRequest{
		From:    sender,
		To:      []string{receiver},
		Subject: subject,
		Html:    htmlBody,
		ReplyTo: reply,
		Headers: map[string]string{
			"Importance":        "high",
			"X-Priority":        "1",
			"X-MSMail-Priority": "High",
		},
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("An error was expected sending an email to %s: %w", receiver, err)
	}

	fmt.Printf("[%s] Email sent successfully...\n", sent.Id)
	return nil
}
