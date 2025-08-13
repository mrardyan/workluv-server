package infrastructure

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"
	"strings"

	"go-server/pkg/config"
)

// EmailService represents the email service interface
type EmailService interface {
	SendEmail(to []string, subject, body string) error
	SendTemplateEmail(to []string, subject, templateName string, data interface{}) error
	SendHTMLEmail(to []string, subject, htmlBody, textBody string) error
	IsEnabled() bool
}

// SMTPEmailService implements EmailService using SMTP
type SMTPEmailService struct {
	config *config.EmailConfig
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
		return NewSMTPEmailService(&cfg.Email)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", cfg.Email.Provider)
	}
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService(cfg *config.EmailConfig) (*SMTPEmailService, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("email service is disabled")
	}

	var auth smtp.Auth
	if cfg.SMTPUsername != "" && cfg.SMTPPassword != "" {
		auth = smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
	}

	service := &SMTPEmailService{
		config: cfg,
		auth:   auth,
	}

	log.Printf("SMTP email service initialized: host=%s, port=%d, tls=%t, ssl=%t",
		cfg.SMTPHost, cfg.SMTPPort, cfg.UseTLS, cfg.UseSSL)

	return service, nil
}

// IsEnabled returns whether the email service is enabled
func (s *SMTPEmailService) IsEnabled() bool {
	return s.config.Enabled
}

// SendEmail sends a plain text email
func (s *SMTPEmailService) SendEmail(to []string, subject, body string) error {
	return s.SendHTMLEmail(to, subject, "", body)
}

// SendHTMLEmail sends an email with both HTML and text content
func (s *SMTPEmailService) SendHTMLEmail(to []string, subject, htmlBody, textBody string) error {
	if !s.config.Enabled {
		log.Printf("Email service is disabled, skipping email send for subject: %s", subject)
		return nil
	}

	msg := s.buildEmailMessage(to, subject, htmlBody, textBody)

	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	// Use standard SMTP with STARTTLS support (most compatible)
	return smtp.SendMail(addr, s.auth, s.config.FromAddress, to, []byte(msg))
}

// SendTemplateEmail sends an email using a template
func (s *SMTPEmailService) SendTemplateEmail(to []string, subject, templateName string, data interface{}) error {
	if !s.config.Enabled {
		log.Printf("Email service is disabled, skipping template email send for subject: %s, template: %s", subject, templateName)
		return nil
	}

	htmlBody, textBody, err := s.renderTemplate(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	return s.SendHTMLEmail(to, subject, htmlBody, textBody)
}

// sendWithTLS sends email using TLS connection
func (s *SMTPEmailService) sendWithTLS(addr, msg string) error {
	host := strings.Split(addr, ":")[0]

	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false, // Always verify certificates in production
	}

	// Connect to server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server with TLS: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate if credentials are provided
	if s.auth != nil {
		if err := client.Auth(s.auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender
	if err := client.Mail(s.config.FromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	recipients := strings.Split(msg, "To: ")[1]
	recipients = strings.Split(recipients, "\r\n")[0]
	toAddresses := strings.Split(recipients, ", ")

	for _, addr := range toAddresses {
		addr = strings.TrimSpace(addr)
		if addr != "" {
			if err := client.Rcpt(addr); err != nil {
				return fmt.Errorf("failed to set recipient %s: %w", addr, err)
			}
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf("Email sent successfully to %d recipients", len(toAddresses))
	return nil
}

// buildEmailMessage builds the complete email message
func (s *SMTPEmailService) buildEmailMessage(to []string, subject, htmlBody, textBody string) string {
	var msg bytes.Buffer

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromAddress))
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

// renderTemplate renders an email template
func (s *SMTPEmailService) renderTemplate(templateName string, data interface{}) (htmlBody, textBody string, err error) {
	if s.config.TemplateDir == "" {
		return "", "", fmt.Errorf("template directory not configured")
	}

	// Try to load HTML template
	htmlPath := filepath.Join(s.config.TemplateDir, templateName+".html")
	htmlTemplate, err := template.ParseFiles(htmlPath)
	if err == nil {
		var htmlBuf bytes.Buffer
		if err := htmlTemplate.Execute(&htmlBuf, data); err != nil {
			return "", "", fmt.Errorf("failed to execute HTML template: %w", err)
		}
		htmlBody = htmlBuf.String()
	}

	// Try to load text template
	textPath := filepath.Join(s.config.TemplateDir, templateName+".txt")
	textTemplate, err := template.ParseFiles(textPath)
	if err == nil {
		var textBuf bytes.Buffer
		if err := textTemplate.Execute(&textBuf, data); err != nil {
			return "", "", fmt.Errorf("failed to execute text template: %w", err)
		}
		textBody = textBuf.String()
	}

	// If neither template exists, return an error
	if htmlBody == "" && textBody == "" {
		return "", "", fmt.Errorf("no templates found for %s (looked for %s.html and %s.txt)", templateName, templateName, templateName)
	}

	return htmlBody, textBody, nil
}

// DisabledEmailService is a no-op implementation when email is disabled
type DisabledEmailService struct{}

func (d *DisabledEmailService) SendEmail(to []string, subject, body string) error {
	log.Printf("Email service disabled, would send email to: %v, subject: %s", to, subject)
	return nil
}

func (d *DisabledEmailService) SendTemplateEmail(to []string, subject, templateName string, data interface{}) error {
	log.Printf("Email service disabled, would send template email to: %v, subject: %s, template: %s", to, subject, templateName)
	return nil
}

func (d *DisabledEmailService) SendHTMLEmail(to []string, subject, htmlBody, textBody string) error {
	log.Printf("Email service disabled, would send HTML email to: %v, subject: %s", to, subject)
	return nil
}

func (d *DisabledEmailService) IsEnabled() bool {
	return false
}

// Helper functions for common email operations

// SendWelcomeEmail sends a welcome email to a new user (using verification_email template)
func SendWelcomeEmail(emailService EmailService, to, username, verificationURL string) error {
	data := map[string]interface{}{
		"Username":        username,
		"VerificationURL": verificationURL,
	}
	return emailService.SendTemplateEmail([]string{to}, "Welcome to Workluv - Verify Your Email", "verification_email", data)
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(emailService EmailService, to, username, resetURL string) error {
	data := map[string]interface{}{
		"Username": username,
		"ResetURL": resetURL,
	}
	return emailService.SendTemplateEmail([]string{to}, "Password Reset Request", "password_reset", data)
}

// SendVerificationEmail sends an email verification email
func SendVerificationEmail(emailService EmailService, to, username, verificationURL string) error {
	data := map[string]interface{}{
		"Username":        username,
		"VerificationURL": verificationURL,
	}
	return emailService.SendTemplateEmail([]string{to}, "Please Verify Your Email", "verification_email", data)
}

// SendSecurityAlertEmail sends a security alert email for suspicious activity
func SendSecurityAlertEmail(emailService EmailService, to, username, deviceInfo string) error {
	data := map[string]interface{}{
		"Username":   username,
		"DeviceInfo": deviceInfo,
	}
	return emailService.SendTemplateEmail([]string{to}, "Security Alert - New Login Detected", "security_alert", data)
}

// SendTwoFactorSetupEmail sends an email for two-factor authentication setup
func SendTwoFactorSetupEmail(emailService EmailService, to, username string) error {
	data := map[string]interface{}{
		"Username": username,
	}
	return emailService.SendTemplateEmail([]string{to}, "Two-Factor Authentication Setup", "two_factor_setup", data)
}
