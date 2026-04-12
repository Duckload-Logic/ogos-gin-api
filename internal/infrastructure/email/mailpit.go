package email

import (
	"context"
	"fmt"
	"net"
	"net/smtp"

	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

type MailPit struct {
	host string
	port int
}

func NewMailPit(host string, port int) (*MailPit, error) {
	return &MailPit{
		host: host,
		port: port,
	}, nil
}

func (m *MailPit) SendEmail(
	ctx context.Context,
	to, subject, body string,
) (bool, error) {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	from := constants.FromEmail()

	// Construct the email message in standard SMTP format (using \r\n)
	msg := fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body

	err := smtp.SendMail(addr, nil, from, []string{to}, []byte(msg))
	if err != nil {
		// Check for Network Timeouts/Connection errors
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return false, fmt.Errorf("[MailPit] Connection timeout: %w", err)
		}

		// Check for SMTP Protocol Errors (e.g., 550 Invalid Recipient)
		// These typically contain the 3-digit SMTP status code
		return false, fmt.Errorf("[MailPit] SMTP Protocol Error: %w", err)
	}

	return true, nil
}

func (m *MailPit) SendOTP(
	ctx context.Context,
	to, otp string,
) (bool, error) {
	return m.SendEmail(ctx, to, "Verification Code", OTP_TEMPLATE(otp))
}
