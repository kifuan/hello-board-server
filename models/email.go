package models

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// Parses HTML body with given data.
func parseEmailBody(tpl *template.Template, data any) (string, error) {
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	return buf.String(), nil
}

func sendEmail(subject, to, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(MailAccount, MailSenderName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(MailHost, MailPort, MailAccount, MailPassword)
	if err := d.DialAndSend(m); err != nil {
		logrus.Warnf("failed to send email to %s: %s.", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logrus.Info("sent email to %s successfully.", to)
	return nil
}

// Sends an email notice to specified id with reply content.
// It does not handle -1 as id.
func sendEmailNoticeIfAllowed(replyName string, content string, id int) error {
	if id == -1 {
		return nil
	}

	replyMessage, err := GetFullMessage(id)
	if err != nil {
		return err
	}

	if !replyMessage.MailNotice {
		return nil
	}

	body, err := parseEmailBody(MailReplyTemplate, map[string]any{
		"name":      replyMessage.Name,
		"content":   content,
		"key":       GenerateUnsubscribeEmailKey(replyMessage.Email),
		"site":      Site,
		"replyName": replyName,
	})

	if err != nil {
		return err
	}

	return sendEmail(MailSubject, replyMessage.Email, body)
}

// Sends notice to owner, if notice is enabled
// and the current message replies to root.
func sendOwnerNoticeIfEnabledAndRoot(m Message) error {
	if !OwnerNotice {
		return nil
	}

	if m.Reply != -1 {
		return nil
	}

	body, err := parseEmailBody(MailOwnerTemplate, map[string]any{
		"owner":   OwnerName,
		"site":    Site,
		"name":    m.Name,
		"content": m.Content,
	})

	if err != nil {
		return err
	}

	return sendEmail(OwnerSubject, OwnerEmail, body)
}
