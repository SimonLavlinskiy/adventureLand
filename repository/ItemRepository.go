package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Item struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"embedded"`
	View   string `gorm:"embedded"`
	Type   string `gorm:"embedded"`
	CanUse bool   `gorm:"embedded"`
	Count  *int   `gorm:"embedded"`
}

func UserGetItem(update tgbotapi.Update, LocationStruct Location) int {
	var resultCell Cellule
	var err error

	err = config.Db.
		Preload("Item").
		First(&resultCell, &Cellule{Map: LocationStruct.Map, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY}).
		Error
	if err != nil {
		panic(err)
	}

	if resultCell.ItemID != nil && resultCell.Item.Type == "food" {
		countAfterUserGetItem := 0

		res := GetOrCreateUserItem(update, *resultCell.Item)

		AddUserItemCount(update, res, *resultCell.Item)

		err = config.Db.Where(&Item{ID: uint(*resultCell.ItemID)}).Updates(Item{Count: &countAfterUserGetItem}).Error
		if err != nil {
			panic(err)
		}
	} else {
		return 0
	}

	return *resultCell.Item.Count
}
