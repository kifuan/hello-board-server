package models

import (
	"bytes"
	"fmt"

	"gopkg.in/gomail.v2"
)

// Sends an email notice to specified id with reply content.
// It does not handle -1 as id.
func sendEmailNotice(replyName string, content string, id int) error {
	if id == -1 {
		return nil
	}

	message, err := GetFullMessage(id)
	if err != nil {
		return err
	}

	if !message.MailNotice {
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(MailAccount, MailSenderName))
	m.SetHeader("To", message.Email)
	m.SetHeader("Subject", MailSubject)

	body, err := parseEmailBody(map[string]any{
		"name":      message.Name,
		"content":   content,
		"key":       GenerateUnsubscribeEmailKey(message.Email),
		"site":      Site,
		"replyName": replyName,
	})

	if err != nil {
		return err
	}

	m.SetBody("text/html", body)
	d := gomail.NewDialer(MailHost, MailPort, MailAccount, MailPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// Parses HTML body with given data.
func parseEmailBody(data any) (string, error) {
	buf := new(bytes.Buffer)
	if err := MailTemplate.Execute(buf, data); err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	return buf.String(), nil
}
