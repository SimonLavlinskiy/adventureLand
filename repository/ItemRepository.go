package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Item struct {
	ID              uint         `gorm:"primaryKey"`
	Name            string       `gorm:"embedded"`
	Description     *string      `gorm:"embedded"`
	View            string       `gorm:"embedded"`
	Type            string       `gorm:"embedded"`
	Cost            *int         `gorm:"embedded"`
	Healing         *int         `gorm:"embedded"`
	Damage          *int         `gorm:"embedded"`
	Satiety         *int         `gorm:"embedded"`
	CanTake         bool         `gorm:"embedded"`
	Instruments     []Instrument `gorm:"many2many:instrument_item;"`
	DressType       *string      `gorm:"embedded"`
	IsBackpack      bool         `gorm:"embedded"`
	IsInventory     bool         `gorm:"embedded"`
	MaxCountUserHas *int         `gorm:"embedded"`
}

type InstrumentItem struct {
	ItemID       int `gorm:"primaryKey"`
	InstrumentID int `gorm:"primaryKey"`
}

func UserGetItem(update tgbotapi.Update, LocationStruct Location, char []string) string {
	resultCell := GetCellule(Cellule{MapsId: *LocationStruct.MapsId, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY})

	if resultCell.ItemID != nil && (resultCell.Item.IsBackpack == true || resultCell.Item.IsInventory == true) {
		status, res := UserGetItemUpdateModels(update, resultCell, char[0])

		if status == "Ok" {
			return "Ты получил " + char[2] + " " + res + "шт."
		} else {
			return res
		}
	}

	return "0"
}

func UserGetItemUpdateModels(update tgbotapi.Update, resultCell Cellule, instrumentView string) (string, string) {
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
				return "Not ok", err
			}

			return "Ok", ToString(countUserGetItem)
		}
		return "Not ok", "Не хватает деняк!"
	}
	return "Not ok", "У тебя уже есть такой!"
}

func canUserTakeItem(item UserItem) bool {
	if item.Item.MaxCountUserHas == nil || *item.Count < *item.Item.MaxCountUserHas {
		return true
	}
	return false
}
