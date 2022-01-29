package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"time"
)

type User struct {
	ID         uint      `gorm:"primaryKey"`
	TgId       uint      `gorm:"embedded"`
	Username   string    `gorm:"embedded"`
	Avatar     string    `gorm:"embedded"`
	LocationId uint      `gorm:"embedded"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	Deleted    bool      `gorm:"embedded"`
}

func GetOrCreateUser(update tgbotapi.Update) User {

	result := User{
		TgId:     uint(update.Message.From.ID),
		Username: update.Message.From.UserName,
	}
	err := config.Db.Where(&User{TgId: uint(update.Message.From.ID)}).FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func UpdateUser(update tgbotapi.Update, UserStruct User) User {
	var err error

	err = config.Db.Where(&User{TgId: uint(update.Message.From.ID)}).Updates(UserStruct).Error
	if err != nil {
		panic(err)
	}

	res := GetOrCreateUser(update)
	return res
}
