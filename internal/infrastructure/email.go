package infrastructure

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"
	"strings"

	"workluv/pkg/config"
)

// EmailService represents the email service interface
type EmailService interface {
	SendEmail(ctx context.Context, to []string, subject, body string) error
	SendTemplateEmail(ctx context.Context, to []string, subject, templateName string, data interface{}) error
	SendHTMLEmail(ctx context.Context, to []string, subject, htmlBody, textBody string) error
	IsEnabled() bool
}

// SMTPEmailService implements EmailService using SMTP
type SMTPEmailService struct {
	config *config.Config
	auth   smtp.Auth
}

// EmailMessage represents an email message structure
type EmailMessage struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	HTMLBody    string
	Attachments []EmailAttachment
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename string
	Content  []byte
	MimeType string
}

// NewEmailService creates a new email service based on configuration
func NewEmailService(cfg *config.Config) (EmailService, error) {
	if !cfg.Email.Enabled {
		return &DisabledEmailService{}, nil
	}

	switch strings.ToLower(cfg.Email.Provider) {
	case "smtp":
		return NewSMTPEmailService(cfg)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", cfg.Email.Provider)
	}
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService(cfg *config.Config) (*SMTPEmailService, error) {
	if !cfg.Email.Enabled {
		return nil, fmt.Errorf("email service is disabled")
	}

	var auth smtp.Auth
	if cfg.Email.SMTPUsername != "" && cfg.Email.SMTPPassword != "" {
		auth = smtp.PlainAuth("", cfg.Email.SMTPUsername, cfg.Email.SMTPPassword, cfg.Email.SMTPHost)
	}

	service := &SMTPEmailService{
		config: cfg,
		auth:   auth,
	}

	log.Printf("SMTP email service initialized: host=%s, port=%d, tls=%t, ssl=%t",
		cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.UseTLS, cfg.Email.UseSSL)

	return service, nil
}

// IsEnabled returns whether the email service is enabled
func (s *SMTPEmailService) IsEnabled() bool {
	return s.config.Email.Enabled
}

// SendEmail sends a plain text email
func (s *SMTPEmailService) SendEmail(ctx context.Context, to []string, subject, body string) error {
	return s.SendHTMLEmail(ctx, to, subject, "", body)
}

// SendHTMLEmail sends an email with both HTML and text content
func (s *SMTPEmailService) SendHTMLEmail(ctx context.Context, to []string, subject, htmlBody, textBody string) error {
	if !s.config.Email.Enabled {
		log.Printf("Email service is disabled, skipping email send for subject: %s", subject)
		return nil
	}

	msg := s.buildEmailMessage(to, subject, htmlBody, textBody)
	addr := fmt.Sprintf("%s:%d", s.config.Email.SMTPHost, s.config.Email.SMTPPort)

	// Use improved connection method based on configuration
	if s.config.Email.UseSSL {
		return s.sendSMTPWithSSL(ctx, addr, to[0], msg)
	} else if s.config.Email.UseTLS {
		return s.sendSMTPWithTLS(ctx, addr, to[0], msg)
	} else {
		// Use standard SMTP (fallback)
		return smtp.SendMail(addr, s.auth, s.config.Email.FromAddress, to, []byte(msg))
	}
}

// SendTemplateEmail sends an email using a template
func (s *SMTPEmailService) SendTemplateEmail(ctx context.Context, to []string, subject, templateName string, data interface{}) error {
	if !s.config.Email.Enabled {
		log.Printf("Email service is disabled, skipping template email send for subject: %s, template: %s", subject, templateName)
		return nil
	}

	htmlBody, textBody, err := s.renderTemplate(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	return s.SendHTMLEmail(ctx, to, subject, htmlBody, textBody)
}

// sendSMTPWithSSL sends email using SSL connection
func (s *SMTPEmailService) sendSMTPWithSSL(ctx context.Context, addr, to, message string) error {
	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName: s.config.Email.SMTPHost,
	}

	// Connect with SSL
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect with SSL: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.Email.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	if s.auth != nil {
		if err := client.Auth(s.auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// Send email
	if err := client.Mail(s.config.Email.FromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	log.Printf("Email sent successfully via SSL to: %s", to)
	return nil
}

// sendSMTPWithTLS sends email using STARTTLS
func (s *SMTPEmailService) sendSMTPWithTLS(ctx context.Context, addr, to, message string) error {
	// Connect to SMTP server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: s.config.Email.SMTPHost,
	}

	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate
	if s.auth != nil {
		if err := client.Auth(s.auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// Send email
	if err := client.Mail(s.config.Email.FromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	log.Printf("Email sent successfully via TLS to: %s", to)
	return nil
}

// buildEmailMessage builds the complete email message
func (s *SMTPEmailService) buildEmailMessage(to []string, subject, htmlBody, textBody string) string {
	var msg bytes.Buffer

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.Email.FromName, s.config.Email.FromAddress))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")

	// Content type based on whether HTML is provided
	if htmlBody != "" && textBody != "" {
		// Multipart email with both HTML and text
		boundary := "boundary123456789"
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundary))

		// Text part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		msg.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		msg.WriteString(textBody)
		msg.WriteString("\r\n\r\n")

		// HTML part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		msg.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		msg.WriteString(htmlBody)
		msg.WriteString("\r\n\r\n")

		msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else if htmlBody != "" {
		// HTML only
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		msg.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		msg.WriteString(htmlBody)
	} else {
		// Plain text only
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		msg.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		msg.WriteString(textBody)
	}

	return msg.String()
}

// renderTemplate renders an email template with improved error handling
func (s *SMTPEmailService) renderTemplate(templateName string, data interface{}) (htmlBody, textBody string, err error) {
	if s.config.Email.TemplateDir == "" {
		return "", "", fmt.Errorf("template directory not configured")
	}

	var htmlErr, textErr error

	// Try to load and render HTML template
	htmlPath := filepath.Join(s.config.Email.TemplateDir, templateName+".html")
	if htmlTemplate, err := template.ParseFiles(htmlPath); err == nil {
		var htmlBuf bytes.Buffer
		if err := htmlTemplate.Execute(&htmlBuf, data); err != nil {
			htmlErr = fmt.Errorf("failed to execute HTML template %s: %w", htmlPath, err)
		} else {
			htmlBody = htmlBuf.String()
		}
	} else {
		htmlErr = fmt.Errorf("failed to parse HTML template %s: %w", htmlPath, err)
	}

	// Try to load and render text template
	textPath := filepath.Join(s.config.Email.TemplateDir, templateName+".txt")
	if textTemplate, err := template.ParseFiles(textPath); err == nil {
		var textBuf bytes.Buffer
		if err := textTemplate.Execute(&textBuf, data); err != nil {
			textErr = fmt.Errorf("failed to execute text template %s: %w", textPath, err)
		} else {
			textBody = textBuf.String()
		}
	} else {
		textErr = fmt.Errorf("failed to parse text template %s: %w", textPath, err)
	}

	// If neither template loaded successfully, return detailed error
	if htmlBody == "" && textBody == "" {
		if htmlErr != nil && textErr != nil {
			return "", "", fmt.Errorf("no templates found for %s - HTML error: %v, Text error: %v", templateName, htmlErr, textErr)
		}
		return "", "", fmt.Errorf("no templates found for %s (looked for %s.html and %s.txt)", templateName, templateName, templateName)
	}

	// Log any partial failures but continue
	if htmlErr != nil {
		log.Printf("Warning: HTML template failed for %s: %v", templateName, htmlErr)
	}
	if textErr != nil {
		log.Printf("Warning: Text template failed for %s: %v", templateName, textErr)
	}

	return htmlBody, textBody, nil
}

// DisabledEmailService is a no-op implementation when email is disabled
type DisabledEmailService struct{}

func (d *DisabledEmailService) SendEmail(ctx context.Context, to []string, subject, body string) error {
	log.Printf("Email service disabled, would send email to: %v, subject: %s", to, subject)
	return nil
}

func (d *DisabledEmailService) SendTemplateEmail(ctx context.Context, to []string, subject, templateName string, data interface{}) error {
	log.Printf("Email service disabled, would send template email to: %v, subject: %s, template: %s", to, subject, templateName)
	return nil
}

func (d *DisabledEmailService) SendHTMLEmail(ctx context.Context, to []string, subject, htmlBody, textBody string) error {
	log.Printf("Email service disabled, would send HTML email to: %v, subject: %s", to, subject)
	return nil
}

func (d *DisabledEmailService) IsEnabled() bool {
	return false
}

// Helper functions for common email operations

// SendWelcomeEmail sends a welcome email to a new user (using verification_email template)
func SendWelcomeEmail(ctx context.Context, emailService EmailService, cfg *config.Config, to, username, token string) error {
	// Create verification URL using client URL
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", cfg.Server.ClientURL, token)

	data := map[string]interface{}{
		"Username":        username,
		"VerificationURL": verificationURL,
		"Token":           token,
	}
	return emailService.SendTemplateEmail(ctx, []string{to}, "Welcome to Workluv - Verify Your Email", "verification_email", data)
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(ctx context.Context, emailService EmailService, cfg *config.Config, to, username, token string) error {
	// Create reset URL using server host
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", cfg.Server.Host, token)

	data := map[string]interface{}{
		"Username": username,
		"ResetURL": resetURL,
		"Token":    token,
	}
	return emailService.SendTemplateEmail(ctx, []string{to}, "Reset Your Password", "password_reset", data)
}

// SendVerificationEmail sends an email verification email
func SendVerificationEmail(ctx context.Context, emailService EmailService, cfg *config.Config, to, username, token string) error {
	// Create verification URL using client URL
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", cfg.Server.ClientURL, token)

	data := map[string]interface{}{
		"Username":        username,
		"VerificationURL": verificationURL,
		"Token":           token,
	}
	return emailService.SendTemplateEmail(ctx, []string{to}, "Verify Your Email Address", "verification_email", data)
}

// SendSecurityAlertEmail sends a security alert email for suspicious activity
func SendSecurityAlertEmail(ctx context.Context, emailService EmailService, to, username, alertType, deviceInfo string) error {
	subject := fmt.Sprintf("Security Alert: %s", alertType)

	data := map[string]interface{}{
		"Username":   username,
		"AlertType":  alertType,
		"DeviceInfo": deviceInfo,
	}
	return emailService.SendTemplateEmail(ctx, []string{to}, subject, "security_alert", data)
}

// SendTwoFactorSetupEmail sends an email for two-factor authentication setup
func SendTwoFactorSetupEmail(ctx context.Context, emailService EmailService, to, username string) error {
	data := map[string]interface{}{
		"Username": username,
	}
	return emailService.SendTemplateEmail(ctx, []string{to}, "Two-Factor Authentication Setup", "two_factor_setup", data)
}
