package models

import (
	"project0/config"
	t "time"
)

type Chat struct {
	ID        uint `gorm:"primaryKey"`
	ExpiredAt t.Time
	Deleted   bool `gorm:"default:false"`
}

func (chat Chat) DeleteChat() {
	err := config.Db.Model(Chat{}).
		Where(&Chat{ID: chat.ID}).
		Update("deleted", true).
		Error
	if err != nil {
		panic(err)
	}
}
