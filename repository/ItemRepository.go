package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Item struct {
	ID          uint     `gorm:"primaryKey"`
	Name        string   `gorm:"embedded"`
	View        string   `gorm:"embedded"`
	Type        string   `gorm:"embedded"`
	CanTake     bool     `gorm:"embedded"`
	CanTakeWith []string `gorm:"type:text[]"`
	Healing     *int     `gorm:"embedded"`
	Damage      *int     `gorm:"embedded"`
	Satiety     *int     `gorm:"embedded"`
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

	if resultCell.ItemID != nil {
		switch resultCell.Item.Type {
		case "food":
			UserGetFood(update, resultCell)
		}
	} else {
		return 0
	}

	return *resultCell.CountItem
}

func UserGetFood(update tgbotapi.Update, resultCell Cellule) {
	countAfterUserGetItem := *resultCell.CountItem - 1

	res := GetOrCreateUserItem(update, *resultCell.Item)

	AddUserItemCount(update, res, resultCell)

	err := config.Db.Where(&Cellule{ID: resultCell.ID}).Updates(Cellule{CountItem: &countAfterUserGetItem}).Error
	if err != nil {
		panic(err)
	}
}
