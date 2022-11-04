package models

import (
	"crypto/md5"
	"fmt"
	"net/url"
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
	Email      string `json:"-"`
	MailNotice bool   `json:"-"`
}

func getGravatarAvatar(email string) string {
	data := []byte(strings.ToLower(email))
	str := fmt.Sprintf("%x", md5.Sum(data))
	path, err := url.JoinPath(GRAVATAR, str)
	if err != nil {
		panic("invalid gravatar url.")
	}
	return path
}

func InsertMessage(m Message) (Message, error) {
	f := Message{
		Avatar:     getGravatarAvatar(m.Email),
		Date:       time.Now().Unix(),
		Name:       m.Name,
		Content:    m.Content,
		Site:       m.Site,
		Reply:      m.Reply,
		Email:      m.Email,
		MailNotice: m.MailNotice,
	}

	if err := db.Create(&f).Error; err != nil {
		return f, err
	}

	return f, nil
}
