package models

import (
	"html/template"
	"os"
	"strconv"
)

var (
	Site            = os.Getenv("SITE")
	DSN             = os.Getenv("DSN")
	UnsubscribeSalt = os.Getenv("UNSUBSCRIBE_SALT")
	PageSize, _     = strconv.Atoi(os.Getenv("PAGE_SIZE"))

	OwnerEmail     = os.Getenv("OWNER_EMAIL")
	OwnerSecret    = os.Getenv("OWNER_SECRET")
	OwnerNotice, _ = strconv.ParseBool(os.Getenv("OWNER_NOTICE"))
	OwnerName      = os.Getenv("OWNER_NAME")

	MailPort, _       = strconv.Atoi(os.Getenv("MAIL_PORT"))
	MailHost          = os.Getenv("MAIL_HOST")
	MailSenderName    = os.Getenv("MAIL_SENDER_NAME")
	MailAccount       = os.Getenv("MAIL_ACCOUNT")
	MailPassword      = os.Getenv("MAIL_PASSWORD")
	MailSubject       = os.Getenv("MAIL_SUBJECT")
	MailReplyTemplate *template.Template
	MailOwnerTemplate *template.Template
)

func init() {
	var err error
	MailReplyTemplate, err = template.ParseFiles("reply.html")
	if err != nil {
		panic(err)
	}
	MailOwnerTemplate, err = template.ParseFiles("owner.html")
	if err != nil {
		panic(err)
	}
}
