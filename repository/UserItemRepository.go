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

func EatItem(update tgbotapi.Update, user User, userItem UserItem) string {
	userItemCount := userItem.Count

	if *userItemCount > 0 {
		itemHeal := uint(*userItem.Item.Healing)
		itemSatiety := uint(*userItem.Item.Satiety)
		itemDamage := uint(*userItem.Item.Damage)

		*userItem.Count = *userItem.Count - 1
		user.Health = user.Health + itemHeal - itemDamage
		user.Satiety = user.Satiety + itemSatiety

		UpdateUser(update, User{
			Health:  user.Health,
			Satiety: user.Satiety,
		})

		UpdateUserItem(user, UserItem{
			ID:    userItem.ID,
			Count: userItem.Count,
		})
	}

	message := "🍽 Ты съел 1 " + userItem.Item.View + ""

	return message
} // Горжусь ею)

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

func UpdateUserInstrument(update tgbotapi.Update, user User, instrument Item) (string, string) {
	userItem, _ := GetUserItem(UserItem{ItemId: int(instrument.ID), UserId: int(user.ID)})

	c := *userItem.CountUseLeft - 1
	if c > 0 {
		UpdateUserItem(user, UserItem{ID: userItem.ID, CountUseLeft: &c})
		return "Ok", ""
	}

	zeroValue := 0

	if *userItem.Count > 1 {
		userItemCount := *userItem.Count - 1
		countUseLeft := userItem.Item.CountUse
		UpdateUserItem(user,
			UserItem{
				ID:           userItem.ID,
				CountUseLeft: countUseLeft,
				Count:        &userItemCount,
			})
	} else {
		UpdateUserItem(user,
			UserItem{
				ID:           userItem.ID,
				CountUseLeft: &zeroValue,
				Count:        &zeroValue,
			})

		if user.LeftHandId != nil && *user.LeftHandId == int(userItem.Item.ID) {
			SetNullUserField(update, "left_hand_id")
		}
		if user.RightHandId != nil && *user.RightHandId == int(userItem.Item.ID) {
			SetNullUserField(update, "right_hand_id")
		}
	}

	return "Not ok", "💥 Инструмент «" + userItem.Item.View + " " + userItem.Item.Name + "» был сломан! 💥"
}
