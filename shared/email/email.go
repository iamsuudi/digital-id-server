package email

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

// Config holds email configuration
type Config struct {
	APIKey    string
	FromEmail string
	FromName  string
}

// EmailParams contains all parameters for sending an email
type EmailParams struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	HTML        string
	Text        string // Optional plain text version
	ReplyTo     string
	Attachments []*resend.Attachment
}

// Service represents the email service
type Service struct {
	client *resend.Client
	config *Config
}

// NewService creates a new email service
func NewService(config *Config) (*Service, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("email API key is required")
	}

	if config.FromEmail == "" {
		config.FromEmail = "onboarding@resend.dev"
	}

	if config.FromName == "" {
		config.FromName = "Oict Digital ID"
	}

	client := resend.NewClient(config.APIKey)

	return &Service{
		client: client,
		config: config,
	}, nil
}

// Send sends an email with the given parameters
func (s *Service) Send(params *EmailParams) (*resend.SendEmailResponse, error) {
	if len(params.To) == 0 {
		return nil, fmt.Errorf("at least one recipient is required")
	}

	if params.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}

	if params.HTML == "" && params.Text == "" {
		return nil, fmt.Errorf("either HTML or text content is required")
	}

	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)

	req := &resend.SendEmailRequest{
		From:        from,
		To:          params.To,
		Cc:          params.Cc,
		Bcc:         params.Bcc,
		Subject:     params.Subject,
		Html:        params.HTML,
		Text:        params.Text,
		ReplyTo:     params.ReplyTo,
		Attachments: params.Attachments,
	}

	return s.client.Emails.Send(req)
}

// DefaultService returns a service initialized with environment variables
func DefaultService() (*Service, error) {
	return NewService(&Config{
		APIKey:    os.Getenv("RESEND_API_KEY"),
		FromEmail: os.Getenv("RESEND_FROM_EMAIL"),
		FromName:  os.Getenv("RESEND_FROM_NAME"),
	})
}
