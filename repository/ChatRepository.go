package repository

import (
	"fmt"
	"project0/config"
	t "time"
)

type Chat struct {
	ID        uint `gorm:"primaryKey"`
	ExpiredAt t.Time
}

func CreateChat(EndTime t.Time) Chat {
	result := Chat{ExpiredAt: EndTime}
	err := config.Db.Create(&result).Error
	if err != nil {
		fmt.Println("Чатик не создался ¯ \\ _ (ツ) _ / ¯ ")
	}
	return result
}
