package models

import (
	"html/template"
	"os"
	"strconv"
)

var (
	Site        = os.Getenv("SITE")
	DSN         = os.Getenv("DSN")
	AdminEmail  = os.Getenv("ADMIN_EMAIL")
	AdminSecret = os.Getenv("ADMIN_SECRET")

	UnsubscribeSalt = os.Getenv("UNSUBSCRIBE_SALT")

	MailPort, _    = strconv.Atoi(os.Getenv("MAIL_PORT"))
	MailHost       = os.Getenv("MAIL_HOST")
	MailSenderName = os.Getenv("MAIL_SENDER_NAME")
	MailAccount    = os.Getenv("MAIL_ACCOUNT")
	MailPassword   = os.Getenv("MAIL_PASSWORD")
	MailSubject    = os.Getenv("MAIL_SUBJECT")
	MailTemplate   *template.Template
)

func init() {
	var err error
	MailTemplate, err = template.ParseFiles("mail.html")
	if err != nil {
		panic(err)
	}
}
