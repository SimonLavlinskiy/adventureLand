package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Item struct {
	ID              uint         `gorm:"primaryKey"`
	Name            string       `gorm:"embedded"`
	View            string       `gorm:"embedded"`
	Type            string       `gorm:"embedded"`
	CanTake         bool         `gorm:"embedded"`
	Instruments     []Instrument `gorm:"many2many:instrument_item;"`
	Healing         *int         `gorm:"embedded"`
	Damage          *int         `gorm:"embedded"`
	Satiety         *int         `gorm:"embedded"`
	Cost            *int         `gorm:"embedded"`
	DressType       *string      `gorm:"embedded"`
	Description     *string      `gorm:"embedded"`
	IsBackpack      bool         `gorm:"embedded"`
	IsInventory     bool         `gorm:"embedded"`
	MaxCountUserHas *int         `gorm:"embedded"`
}

type InstrumentItem struct {
	ItemID       int `gorm:"primaryKey"`
	InstrumentID int `gorm:"primaryKey"`
}

func UserGetItem(update tgbotapi.Update, LocationStruct Location, char []string) string {
	var resultCell Cellule
	var err error

	err = config.Db.
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.ItemsResult").
		Preload("Item.Instruments.NextStageItem").
		First(&resultCell, &Cellule{MapsId: *LocationStruct.MapsId, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY}).
		Error
	if err != nil {
		panic(err)
	}

	if len(resultCell.Item.Instruments) != 0 && char[0] == "👋" {
		return "Не не, меня не наебешь!"
	}

	if resultCell.ItemID != nil && (resultCell.Item.IsBackpack == true || resultCell.Item.IsInventory == true) {
		res := UserGetItemUpdateModels(update, resultCell, char[0])

		return "Ты получил " + char[2] + " " + res + "шт." //\nВ ячейке: " + ToString(*resultCell.CountItem-1) + " шт."
	}

	return "0"
}

func UserGetItemUpdateModels(update tgbotapi.Update, resultCell Cellule, instrumentView string) string {
	userTgId := GetUserTgId(update)
	user := GetUser(User{TgId: userTgId})
	itemCost := 0

	resUserItem := GetOrCreateUserItem(update, *resultCell.Item)
	if canUserTakeItem(resUserItem) {
		if resultCell.Item.Cost == nil || *user.Money >= *resultCell.Item.Cost {
			if resultCell.Item.Cost != nil {
				itemCost = *resultCell.Item.Cost
			}
			updateUserMoney := *user.Money - itemCost
			err, countUserGetItem := AddUserItemCount(update, resUserItem, resultCell, updateUserMoney, instrumentView)

			if err != "Ok" {
				panic(err)
			}

			return ToString(countUserGetItem)
		}
		return "Не хватает деняк!"
	}
	return "У тебя уже есть такой!"
}

func canUserTakeItem(item UserItem) bool {
	if item.Item.MaxCountUserHas == nil || *item.Count < *item.Item.MaxCountUserHas {
		return true
	}
	return false
}
