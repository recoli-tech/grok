package grok

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

// Mail ...
type Mail struct {
	MailFrom  string
	MailTo    string
	PlainText string
	HTML      string
	Subject   string
}

// MailProvider ...
type MailProvider interface {
	Send(*Mail) error
}

// CreateMailProvider ...
func CreateMailProvider(settings *MailSettings) MailProvider {
	if settings.Provider == "fake" {
		return NewFakeMailProvider(settings.Fake.ShouldReturnError)
	}

	apiKey := settings.SendGrid.APIKey
	if settings.SendGrid.FromEnv {
		apiKey = os.Getenv(settings.SendGrid.APIKeyEnv)
	}

	return NewSendGridMailProvider(apiKey)
}

type sendGridProvider struct {
	client   *sendgrid.Client
	settings *MailSettings
}

// NewSendGridMailProvider ...
func NewSendGridMailProvider(apiKey string) MailProvider {
	return &sendGridProvider{
		client: sendgrid.NewSendClient(apiKey),
	}
}

func (s *sendGridProvider) Send(m *Mail) error {
	message := mail.NewV3MailInit(
		mail.NewEmail("", m.MailFrom),
		m.Subject,
		mail.NewEmail("", m.MailTo),
		mail.NewContent("text/plain", m.PlainText),
	)

	res, err := s.client.Send(message)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusAccepted {
		logrus.
			WithField("response", res).
			Error("sendgrid api error")

		return fmt.Errorf("error sending mail to %s:\n\n%s", m.MailTo, m.PlainText)
	}

	return nil
}

type fakeMailProvider struct {
	shouldReturnError bool
}

// NewFakeMailProvider ...
func NewFakeMailProvider(shouldReturnError bool) MailProvider {
	return &fakeMailProvider{shouldReturnError: shouldReturnError}
}

func (s *fakeMailProvider) Send(*Mail) error {
	if s.shouldReturnError {
		return errors.New("fake mail provider error")
	}

	return nil
}
