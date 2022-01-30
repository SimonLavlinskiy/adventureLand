package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Location struct {
	ID       uint `gorm:"primaryKey"`
	UserTgId uint
	UserID   uint
	AxisX    uint64 `gorm:"embedded"`
	AxisY    uint64 `gorm:"embedded"`
	Maps     string `gorm:"embedded"`
}

func GetOrCreateLocation(update tgbotapi.Update) Location {
	res := GetOrCreateUser(update)

	result := Location{
		UserTgId: uint(update.Message.From.ID),
		AxisX:    5,
		AxisY:    5,
		Maps:     "Main",
	}

	err := config.Db.Where(&Location{UserID: res.ID}).FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func UpdateLocation(update tgbotapi.Update, LocationStruct Location) Location {
	var err error

	err = config.Db.Where(&Location{UserTgId: uint(update.Message.From.ID)}).Updates(LocationStruct).Error
	if err != nil {
		panic(err)
	}

	res := GetOrCreateLocation(update)
	fmt.Println(res)

	return res
}
