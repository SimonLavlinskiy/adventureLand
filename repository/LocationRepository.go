package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Location struct {
	ID       uint `gorm:"primaryKey"`
	UserTgId uint
	UserID   uint
	AxisX    int    `gorm:"embedded"`
	AxisY    int    `gorm:"embedded"`
	Map      string `gorm:"embedded"`
}

func GetOrCreateMyLocation(update tgbotapi.Update) Location {
	resUser := GetOrCreateUser(update)

	result := Location{
		UserTgId: uint(update.Message.From.ID),
		AxisX:    3,
		AxisY:    2,
		Map:      "Ekaterensky",
	}

	err := config.Db.Where(&Location{UserID: resUser.ID}).FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func UpdateLocation(update tgbotapi.Update, LocationStruct Location) Location {
	myLocation := GetOrCreateMyLocation(update)

	if myLocation.Map != LocationStruct.Map {
		resCellule := GetCellule(Cellule{Map: myLocation.Map, AxisX: LocationStruct.AxisX, AxisY: LocationStruct.AxisY})
		resTeleport := GetTeleport(Teleport{ID: resCellule.TeleportId})
		LocationStruct = Location{
			AxisX: resTeleport.StartX,
			AxisY: resTeleport.StartY,
			Map:   resTeleport.Map,
		}
	}

	var result Cellule
	var err error

	err = config.Db.First(&result, &Cellule{Map: LocationStruct.Map, AxisX: LocationStruct.AxisX, AxisY: LocationStruct.AxisY}).Error
	if err != nil {
		if err.Error() == "record not found" {
			return myLocation
		}
		panic(err)
	}

	if !result.CanStep && result.Type == "teleport" {
		err = config.Db.Where(&Location{UserTgId: uint(update.Message.From.ID)}).Updates(LocationStruct).Error
		if err != nil {
			panic(err)
		}
	} else if !result.CanStep {
		return myLocation
	} else {
		err = config.Db.Where(&Location{UserTgId: uint(update.Message.From.ID)}).Updates(LocationStruct).Error
		if err != nil {
			panic(err)
		}
	}

	myLocation = GetOrCreateMyLocation(update)
	return myLocation
}
