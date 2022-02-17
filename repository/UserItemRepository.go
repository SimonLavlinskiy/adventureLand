package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type UserItem struct {
	ID     int  `gorm:"primaryKey"`
	Count  *int `gorm:"embedded"`
	UserId int  `gorm:"embedded"`
	User   User
	ItemId int `gorm:"embedded"`
	Item   Item
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
		Count:  &countItem,
		UserId: userId,
		ItemId: int(item.ID),
	}
	err := config.Db.Preload("Item").Where(UserItem{UserId: userId, ItemId: int(item.ID)}).FirstOrCreate(&result).Error
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

func AddUserItemCount(update tgbotapi.Update, userItem UserItem, cellule Cellule, updateUserMoney int, instrumentView string) (string, int) {
	userTgId := GetUserTgId(update)
	resUser := GetUser(User{TgId: userTgId})
	userId := int(resUser.ID)
	countAfterUserGetItem := *cellule.CountItem - 1
	var sumCountItemResult int

	if instrumentView == "👋" {
		sumCountItemResult = *userItem.Count + 1

		UpdateUserItem(User{ID: uint(userId)}, UserItem{ID: userItem.ID, Count: &sumCountItemResult})
		UpdateUser(update, User{Money: &updateUserMoney})
		UpdateCellule(cellule.ID, Cellule{CountItem: &countAfterUserGetItem})

		return "Ok", 1
	} else {
		for _, instrument := range cellule.Item.Instruments {
			if instrumentView == instrument.Good.View {
				sumCountItemResult = *userItem.Count + *instrument.CountResultItem

				UpdateUserItem(User{ID: uint(userId)}, UserItem{ID: userItem.ID, Count: &sumCountItemResult})
				UpdateUser(update, User{Money: &updateUserMoney})
				UpdateCellule(cellule.ID, Cellule{CountItem: &countAfterUserGetItem})

				return "Ok", *instrument.CountResultItem
			}
		}
	}
	return "err", 0
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
