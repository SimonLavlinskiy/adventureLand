package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Location struct {
	ID       uint `gorm:"primaryKey"`
	UserTgId uint
	UserID   uint
	AxisX    uint   `gorm:"embedded"`
	AxisY    uint   `gorm:"embedded"`
	Map      string `gorm:"embedded"`
}

func GetOrCreateMyLocation(update tgbotapi.Update) Location {
	res := GetOrCreateUser(update)

	result := Location{
		UserTgId: uint(update.Message.From.ID),
		AxisX:    1,
		AxisY:    1,
		Map:      "Ekaterensky",
	}

	err := config.Db.Where(&Location{UserID: res.ID}).FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func UpdateLocation(update tgbotapi.Update, LocationStruct Location) Location {
	myLocation := GetOrCreateMyLocation(update)

	var result Cellule

	var err error

	err = config.Db.First(&result, &Cellule{Map: LocationStruct.Map, AxisX: LocationStruct.AxisX, AxisY: LocationStruct.AxisY}).Error
	if err != nil {
		if err.Error() == "record not found" {
			return myLocation
		}
		panic(err)
	}

	if !result.CanStep {
		return myLocation
	}

	err = config.Db.Where(&Location{UserTgId: uint(update.Message.From.ID)}).Updates(LocationStruct).Error
	if err != nil {
		panic(err)
	}

	myLocation = GetOrCreateMyLocation(update)
	return myLocation
}
