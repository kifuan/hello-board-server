package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

func GenerateUnsubscribeEmailKey(email string) string {
	data := []byte(fmt.Sprintf("%s%s", email, UnsubscribeSalt))
	return fmt.Sprintf("%x", md5.Sum(data))
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
	if m.Email == OwnerEmail {
		return errors.New("Don't try to use my email!")
	}

	if m.Email == OwnerSecret {
		m.Email = OwnerEmail
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
func UnsubscribeEmail(email string) error {
	if err := db.Model(&Message{}).Where("email=?", email).Update("mail_notice", false).Error; err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	return nil
}

func CountMessages() (count int64, err error) {
	if err = db.Model(&Message{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count: %w", err)
	}
	return count, nil
}
