package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"
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
}

func getGravatarAvatar(email string) string {
	data := []byte(strings.ToLower(email))
	return fmt.Sprintf("%x", md5.Sum(data))
}

func InsertMessage(m *Message) error {
	// Check for the email first.
	if m.Email == ADMIN_EMAIL {
		return errors.New("Don't try to use my email!")
	}

	if m.Email == ADMIN_SECRET {
		m.Email = ADMIN_EMAIL
	}

	// Generate avatar in md5 and time.
	m.Avatar = getGravatarAvatar(m.Email)
	m.Date = time.Now().UnixMilli()

	// TODO notice the message it replied to.

	if err := db.Create(m).Error; err != nil {
		return err
	}

	return nil
}

func GetAllMessages() (messages []Message, err error) {
	if err = db.Select("id, avatar, date, name, content, site, reply").Find(&messages).Error; err != nil {
		return messages, fmt.Errorf("failed to get all messages: %w", err)
	}
	return
}
