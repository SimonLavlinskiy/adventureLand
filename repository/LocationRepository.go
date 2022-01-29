package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Location struct {
	ID     uint   `gorm:"primaryKey"`
	UserId uint   `gorm:"embedded"`
	AxisX  uint64 `gorm:"embedded"`
	AxisY  uint64 `gorm:"embedded"`
	Maps   string `gorm:"embedded"`
}

func GetOrCreateLocation(update tgbotapi.Update) Location {
	res := GetOrCreateUser(update)

	result := Location{
		UserId: res.ID,
		AxisX:  5,
		AxisY:  5,
		Maps:   "Main",
	}

	err := config.Db.Where(&Location{ID: uint(res.LocationId)}).FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}
