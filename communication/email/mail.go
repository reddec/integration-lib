package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"io"
	"net"
	"strconv"
	"strings"
)

// Default mail configuration from env variables
//
//     MAIL_SERVER - mandatory, address with port (usually: 25 for non-TLS, 465 for TLS) of SMTP server
//     MAIL_LOGIN  - mandatory, login/sender email address
func Default() *Mail {
	cfg := &Mail{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	if cfg.Server == "" {
		panic("SMTP server (MAIL_SERVER) required")
	}
	if cfg.Login == "" {
		panic("SMTP login/sender required (MAIL_LOGIN)")
	}
	return cfg
}

// Mail configuration
type Mail struct {
	Server   string   `env:"MAIL_SERVER"`              // SMTP server (mandatory)
	TLS      bool     `env:"MAIL_TLS"`                 // Enable TLS connection
	Login    string   `env:"MAIL_LOGIN"`               // Login name and sender email (mandatory)
	Password string   `env:"MAIL_PASSWORD"`            // Send password (if no password defined, AUTH will not be used)
	To       []string `env:"MAIL_TO" envSeparator:":"` // Default receivers
}

// Send plain/text message to default addresses (see SendTextToContext)
func (mail Mail) SendText(subject, text string) error {
	return mail.SendTextContext(context.Background(), subject, text)
}

// Send plain/text message to default addresses (see SendTextToContext)
func (mail Mail) SendTextContext(ctx context.Context, subject, text string) error {
	return mail.SendTextToContext(ctx, subject, text, mail.To)
}

// Send plain/text message to custom addresses (see SendTextToContext)
func (mail Mail) SendTextTo(subject, text string, to []string) error {
	return mail.SendTextToContext(context.Background(), subject, text, to)
}

// Send plain/text message to custom addresses. Addresses will be visible in To header
func (mail Mail) SendTextToContext(ctx context.Context, subject, text string, to []string) error {
	return mail.SendContext(ctx, []byte(text), subject, "text/plain", to)
}

// Send arbitrary content over SMTP server. Addresses will be visible in To header
func (mail Mail) SendContext(ctx context.Context, data []byte, subject, contentType string, to []string) error {
	host, _, err := net.SplitHostPort(mail.Server)
	if err != nil {
		return fmt.Errorf("mail: parse server: %w", err)
	}
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", mail.Server)
	if err != nil {
		return fmt.Errorf("mail: dial to server: %w", err)
	}
	defer conn.Close()
	if mail.TLS {
		conn = tls.Client(conn, &tls.Config{
			ServerName: host,
		})
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("mail: create mail client: %w", err)
	}
	defer client.Close()
	if mail.Password != "" {
		auth := sasl.NewPlainClient("", mail.Login, mail.Password)
		err = client.Auth(auth)
		if err != nil {
			return fmt.Errorf("mail: authorize: %w", err)
		}
	}

	err = client.Mail(mail.Login, nil)
	if err != nil {
		return fmt.Errorf("mail: start mail: %w", err)
	}

	for _, addr := range to {
		err = client.Rcpt(addr)
		if err != nil {
			return fmt.Errorf("mail: define recepeint %s: %w", addr, err)
		}
	}
	buffer := &bytes.Buffer{}
	buffer.WriteString("Subject: " + subject + "\r\n")
	buffer.WriteString("To: " + strings.Join(to, ",") + "\r\n")
	buffer.WriteString("From: " + mail.Login + "\r\n")
	buffer.WriteString("Content-Type: " + contentType + "\r\n")
	buffer.WriteString("Content-Length: " + strconv.Itoa(len(data)) + "\r\n")
	buffer.WriteString("\r\n")
	buffer.Write(data)
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("mail: initialize stream: %w", err)
	}
	_, err = io.Copy(w, buffer)
	if err != nil {
		_ = w.Close()
		return fmt.Errorf("mail: write message: %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("mail: flush message: %w", err)
	}
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("mail: close client: %w", err)
	}
	return nil
}
