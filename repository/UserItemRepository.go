package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type UserItem struct {
	ID           int  `gorm:"primaryKey"`
	Count        *int `gorm:"embedded"`
	CountUseLeft *int `gorm:"embedded"`
	UserId       int  `gorm:"embedded"`
	User         User
	ItemId       int `gorm:"embedded"`
	Item         Item
}

func GetUserItem(userItem UserItem) (UserItem, error) {
	var result UserItem
	err := config.Db.Preload("Item").Preload("User").Where(userItem).First(&result).Error

	if err != nil {
		fmt.Println("Item not found")
	}

	return result, err
}

func GetOrCreateUserItem(update tgbotapi.Update, item Item) UserItem {
	userTgId := GetUserTgId(update)
	resUser := GetUser(User{TgId: userTgId})
	userId := int(resUser.ID)
	countItem := 0
	result := UserItem{
		Count:        &countItem,
		UserId:       userId,
		ItemId:       int(item.ID),
		CountUseLeft: item.CountUse,
	}
	err := config.Db.
		Preload("Item").
		Preload("Item.Instruments").
		Where(UserItem{UserId: userId, ItemId: int(item.ID)}).
		FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func GetUserItems(userId uint) []UserItem {
	var result []UserItem

	err := config.Db.
		Preload("Item").
		Preload("User").
		Where(UserItem{UserId: int(userId)}).
		Where("count > 0").
		Find(&result).
		Error
	if err != nil {
		panic(err)
	}

	return result
}

func GetBackpackItems(userId uint) []UserItem {
	userItems := GetUserItems(userId)

	var backpackUserItem []UserItem

	for _, userItem := range userItems {
		if userItem.Item.IsBackpack == true {
			backpackUserItem = append(backpackUserItem, userItem)
		}
	}

	return backpackUserItem
}

func GetInventoryItems(userId uint) []UserItem {
	userItems := GetUserItems(userId)

	var inventoryUserItem []UserItem

	for _, userItem := range userItems {
		if userItem.Item.IsInventory == true {
			inventoryUserItem = append(inventoryUserItem, userItem)
		}
	}

	return inventoryUserItem
}

func UpdateUserItem(user User, userItem UserItem) {
	err := config.Db.Where(UserItem{UserId: int(user.ID), ID: userItem.ID}).Updates(userItem).Error
	if err != nil {
		panic(err)
	}
}

func UpdateModelsWhenUserGetItem(update tgbotapi.Update, user User, userItem UserItem, cellule Cellule, instrument *Instrument, sumCountItemResult int) {
	itemCost := 0
	if cellule.Item.Cost != nil {
		itemCost = *cellule.Item.Cost
	}

	updateUserMoney := *user.Money - itemCost

	UpdateUserItem(User{ID: user.ID}, UserItem{ID: userItem.ID, Count: &sumCountItemResult, CountUseLeft: userItem.Item.CountUse})
	UpdateUser(update, User{Money: &updateUserMoney})

	countAfterUserGetItem := *cellule.CountItem - 1

	if instrument == nil {
		UpdateCellule(cellule.ID, Cellule{CountItem: &countAfterUserGetItem})
	} else {
		fmt.Println("121")
		//UpdateCelluleWhenUserUseIt(cellule, *instrument)
	}
}

func EatItem(update tgbotapi.Update, user User, userItem UserItem) string {
	userItemCount := userItem.Count

	if *userItemCount > 0 {
		itemHeal := userItem.Item.Healing
		itemSatiety := userItem.Item.Satiety
		itemDamage := userItem.Item.Damage
		*userItemCount = *userItemCount - 1

		userHealth := user.Health + uint(*itemHeal) - uint(*itemDamage)
		userSatiety := user.Satiety + uint(*itemSatiety)

		userUpdate := User{
			Health:  userHealth,
			Satiety: userSatiety,
		}

		userItemUpdate := UserItem{
			ID:    userItem.ID,
			Count: userItemCount,
		}

		UpdateUser(update, userUpdate)
		UpdateUserItem(user, userItemUpdate)
	}
	message := "🍽 Ты съел 1 " + userItem.Item.View + ""
	return message
}

func GetFullDescriptionOfUserItem(userItem UserItem) string {
	userItem, _ = GetUserItem(userItem)
	var fullDescriptionUserItem string
	if userItem.Item.IsInventory == true {
		fullDescriptionUserItem = userItem.Item.View + " *" + userItem.Item.Name + "* - " + ToString(*userItem.Count) + " шт.\n" +
			"*Сила*: +" + ToString(*userItem.Item.Damage) + "💥\n"
	} else if userItem.Item.IsBackpack == true {
		fullDescriptionUserItem = userItem.Item.View + " *" + userItem.Item.Name + "* - " + ToString(*userItem.Count) + " шт.\n" +
			"*Здоровье*: +" + ToString(*userItem.Item.Healing) + " ♥️️\n" +
			"*Сытость*: +" + ToString(*userItem.Item.Satiety) + " \U0001F9C3 \n"
	}
	itemDescription := "Описания нет("

	if userItem.Item.Description != nil {
		itemDescription = "*Описание*: " + *userItem.Item.Description
	}

	return fullDescriptionUserItem + itemDescription
}

func UpdateUserInstrument(update tgbotapi.Update, user User, instrument Item) string {
	userItem, _ := GetUserItem(UserItem{ItemId: int(instrument.ID), UserId: int(user.ID)})

	c := *userItem.CountUseLeft - 1
	if c > 0 {
		UpdateUserItem(user, UserItem{ID: userItem.ID, CountUseLeft: &c})
		return ""
	}
	countsValue := 0

	if *userItem.Count > 1 {
		userItemCount := *userItem.Count - 1
		countUseLeft := userItem.Item.CountUse
		UpdateUserItem(user, UserItem{ID: userItem.ID, CountUseLeft: countUseLeft, Count: &userItemCount})
	} else {
		UpdateUserItem(user, UserItem{ID: userItem.ID, CountUseLeft: &countsValue, Count: &countsValue})

		switch int(userItem.Item.ID) {
		case *user.LeftHandId:
			SetNullUserField(update, "left_hand_id")
		case *user.RightHandId:
			SetNullUserField(update, "right_hand_id")
		}
	}
	return "\n\n💥 Инструмент " + userItem.Item.View + " " + userItem.Item.Name + " был сломан! 💥"

}
