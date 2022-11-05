package models

import (
	"html/template"
	"os"
	"strconv"
)

var (
	DSN          = os.Getenv("DSN")
	ADMIN_EMAIL  = os.Getenv("ADMIN_EMAIL")
	ADMIN_SECRET = os.Getenv("ADMIN_SECRET")

	MAIL_PORT, _     = strconv.Atoi(os.Getenv("MAIL_PORT"))
	MAIL_HOST        = os.Getenv("MAIL_HOST")
	MAIL_SENDER_NAME = os.Getenv("MAIL_SENDER_NAME")
	MAIL_ACCOUNT     = os.Getenv("MAIL_ACCOUNT")
	MAIL_PASSWORD    = os.Getenv("MAIL_PASSWORD")
	MAIL_SUBJECT     = os.Getenv("MAIL_SUBJECT")

	MailTemplate *template.Template
)

func init() {
	var err error
	MailTemplate, err = template.ParseFiles("mail.html")
	if err != nil {
		panic(err)
	}
}
