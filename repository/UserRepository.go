package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"time"
)

type User struct {
	ID         uint      `gorm:"primaryKey"`
	TgId       uint      `gorm:"embedded"`
	Username   string    `gorm:"embedded"`
	LocationId *uint     `gorm:"embedded"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	Deleted    bool      `gorm:"embedded"`
}

func GetOrCreateUser(update tgbotapi.Update) User {

	result := User{
		TgId:     uint(update.Message.From.ID),
		Username: update.Message.From.UserName,
	}
	fmt.Println(update.Message.From.ID)
	err := config.Db.Where(&User{TgId: uint(update.Message.From.ID)}).FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", result)

	return result
}

func UpdateUsername(update tgbotapi.Update, originalUsername bool) User {
	var err error

	result := User{
		Username: update.Message.From.UserName,
	}
	if originalUsername {
		err = config.Db.Where(&User{TgId: uint(update.Message.From.ID)}).Updates(result).Error
	} else {
		err = config.Db.Where(&User{TgId: uint(update.Message.From.ID)}).Updates(User{Username: "Пися"}).Error
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", result)

	return result
}
