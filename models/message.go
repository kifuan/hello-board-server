package models

import (
	"crypto/md5"
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

func InsertMessage(m Message) (Message, error) {
	f := Message{
		Avatar:     getGravatarAvatar(m.Email),
		Date:       time.Now().UnixMilli(),
		Name:       m.Name,
		Content:    m.Content,
		Site:       m.Site,
		Reply:      m.Reply,
		Email:      m.Email,
		MailNotice: m.MailNotice,
	}

	// TODO notice the message it replied to.

	if err := db.Create(&f).Error; err != nil {
		return f, err
	}

	return f, nil
}

func GetAllMessages() (messages []Message, err error) {
	if err = db.Select("id, avatar, date, name, content, site, reply").Find(&messages).Error; err != nil {
		return messages, fmt.Errorf("failed to get all messages: %w", err)
	}
	return
}
