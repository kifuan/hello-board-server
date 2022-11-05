package models

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Message struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Avatar     string `json:"avatar"`
	Date       int64  `json:"date"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	Site       string `json:"site"`
	Reply      int    `json:"reply"`
	Email      string `json:"email,omitempty"`
	MailNotice bool   `json:"mailNotice,omitempty"`
	Owner      bool   `json:"owner,omitempty"`
}

// Generates unsubscribe key for validation.
func (m Message) GenerateUnsubscribeKey() string {
	str := fmt.Sprintf("%d%s%s", m.ID, m.Email, UnsubscribeSalt)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// Initializes the message from uploading.
func (m *Message) initFromUploading() {
	data := []byte(strings.ToLower(m.Email))
	m.Avatar = fmt.Sprintf("%x", md5.Sum(data))
	m.Date = time.Now().UnixMilli()
}

// Makes it reply to a root message.
func (m *Message) fixReply() error {
	if m.Reply == -1 {
		return nil
	}
	m.Content = fmt.Sprintf("Reply to #%d: %s", m.Reply, m.Content)

	reply := *m
	var err error
	for reply.Reply != -1 {
		reply, err = GetFullMessage(reply.Reply)
		if err != nil {
			return fmt.Errorf("error fixing reply: %w", err)
		}
	}
	m.Reply = reply.ID
	return nil
}

// Inserts the message to database.
func InsertMessage(m *Message) error {
	// Check for the email first.
	if m.Email == AdminEmail {
		return errors.New("Don't try to use my email!")
	}

	if m.Email == AdminSecret {
		m.Email = AdminEmail
		m.Owner = true
	}

	// Initializes it.
	m.initFromUploading()

	if err := sendEmailNotice(m.Name, m.Content, m.Reply); err != nil {
		// We don't return this error because
		// current message has no problem. It is
		// due to the reply message.
		logrus.Error(err)
	}

	// Fix the reply id, letting it reply to the root message.
	if err := m.fixReply(); err != nil {
		return err
	}

	if err := db.Create(m).Error; err != nil {
		return err
	}

	return nil
}

func getRepliesFor(id int, page int) (messages []Message, err error) {
	d := db.Select("id, avatar, date, name, content, site, reply, owner").Where("reply=?", id).Order("date DESC")
	// Only limit root messages.
	if id == -1 {
		d = d.Offset(page * PageSize).Limit(PageSize)
	}

	if err = d.Find(&messages).Error; err != nil {
		err = fmt.Errorf("failed to get replies for %d: %w", id, err)
		return messages, err
	}
	return messages, nil
}

// Gets messages for specified page.
func GetMessages(page int) (messages []Message, err error) {
	root, err := getRepliesFor(-1, page)
	if err != nil {
		return messages, err
	}
	for _, r := range root {
		replies, err := getRepliesFor(r.ID, page)
		if err != nil {
			return messages, err
		}
		messages = append(messages, replies...)
	}
	return append(root, messages...), nil
}

// Gets all messages, without Email and MailNotice fields.
func GetAllMessages() (messages []Message, err error) {
	if err = db.Select("id, avatar, date, name, content, site, reply, owner").Order("date DESC").Find(&messages).Error; err != nil {
		return messages, fmt.Errorf("failed to get all messages: %w", err)
	}
	return messages, nil
}

// Gets full message, including Email and MailNotice.
func GetFullMessage(id int) (m Message, err error) {
	if err = db.Where("id=?", id).Find(&m).Error; err != nil {
		err = fmt.Errorf("failed to find message with id %d: %w", id, err)
		return m, err
	}
	return m, nil
}

// Unsubscribes mail notice.
func UnsubscribeMailNotice(id int) error {
	if err := db.Model(&Message{}).Where("id=?", id).Update("mail_notice", false).Error; err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	return nil
}

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
		"id":        message.ID,
		"key":       message.GenerateUnsubscribeKey(),
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

func CountMessages() (count int64, err error) {
	if err = db.Model(&Message{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count: %w", err)
	}
	return count, nil
}
