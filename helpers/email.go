package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/txsvc/stdlib/v2"
)

const (
	MailgunEmailDomainENV = "MAILGUN_EMAIL_DOMAIN"
	MailgunApiKeyENV      = "MAILGUN_API_KEY"
)

// MailgunSimpleEmail is a minimal email implementation
func MailgunSimpleEmail(sender, recipient, subject, body string) error {
	domain := stdlib.GetString(MailgunEmailDomainENV, "")
	if domain == "" {
		return fmt.Errorf("invalid email configuration")
	}
	apiKey := stdlib.GetString(MailgunApiKeyENV, "")
	if apiKey == "" {
		return fmt.Errorf("invalid email configuration")
	}

	mg := mailgun.NewMailgun(domain, apiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	message := mg.NewMessage(sender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
