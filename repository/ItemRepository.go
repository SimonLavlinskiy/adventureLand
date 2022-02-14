package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Item struct {
	ID            uint    `gorm:"primaryKey"`
	Name          string  `gorm:"embedded"`
	View          string  `gorm:"embedded"`
	Type          string  `gorm:"embedded"`
	CanTake       bool    `gorm:"embedded"`
	CanTakeWith   *Item   `gorm:"foreignKey:CanTakeWithId"`
	CanTakeWithId *int    `gorm:"embedded"`
	Healing       *int    `gorm:"embedded"`
	Damage        *int    `gorm:"embedded"`
	Satiety       *int    `gorm:"embedded"`
	Cost          *int    `gorm:"embedded"`
	DressType     *string `gorm:"embedded"`
}

func UserGetItem(update tgbotapi.Update, LocationStruct Location, char []string) string {
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
		case "food", "pick", "axe":
			res := UserGetItemUpdateModels(update, resultCell)
			if res != "Ok" {
				return res
			}
		}
	} else {
		return "0"
	}
	return "Ты взял " + char[2] + " 1шт.\nВ ячейке: " + ToString(*resultCell.CountItem-1) + " шт."
}

func UserGetItemUpdateModels(update tgbotapi.Update, resultCell Cellule) string {
	countAfterUserGetItem := *resultCell.CountItem - 1
	user := GetUser(User{TgId: uint(update.Message.From.ID)})

	resUserItem := GetOrCreateUserItem(update, *resultCell.Item)
	if canUserTakeItem(resUserItem) {
		if *user.Money >= *resultCell.Item.Cost {
			updateUserMoney := *user.Money - *resultCell.Item.Cost
			AddUserItemCount(update, resUserItem, resultCell, updateUserMoney)
			UpdateCellule(resultCell.ID, Cellule{CountItem: &countAfterUserGetItem})
			return "Ok"
		}
		return "Не хватает деняк!"
	}
	return "У тебя уже есть такой!"
}

func canUserTakeItem(item UserItem) bool {
	if (item.Item.Type == "pick" || item.Item.Type == "axe") && *item.Count < 1 {
		return true
	} else if item.Item.Type == "food" {
		return true
	}
	return false
}
