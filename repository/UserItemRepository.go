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

func (ui UserItem) GetUserItem() UserItem {
	zero := 0
	result := UserItem{
		UserId:       int(ui.User.ID),
		ItemId:       ui.ItemId,
		Count:        &zero,
		CountUseLeft: ui.CountUseLeft,
	}
	err := config.Db.
		Preload("Item").
		Preload("User").
		Where(ui).
		FirstOrCreate(&result).
		Error

	if err != nil {
		panic(err)
	}

	return result
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

func (ui UserItem) UpdateUserItem(user User) {
	err := config.Db.Where(UserItem{UserId: int(user.ID), ID: ui.ID}).Updates(ui).Error
	if err != nil {
		panic(err)
	}
}

func (ui UserItem) EatItem(update tgbotapi.Update, user User) string {
	userItemCount := ui.Count

	if *userItemCount > 0 {
		itemHeal := uint(*ui.Item.Healing)
		itemSatiety := uint(*ui.Item.Satiety)
		itemDamage := uint(*ui.Item.Damage)

		*ui.Count = *ui.Count - 1
		user.Health = user.Health + itemHeal - itemDamage
		user.Satiety = user.Satiety + itemSatiety

		User{
			Health:  user.Health,
			Satiety: user.Satiety,
		}.UpdateUser(update)

		UserItem{
			ID:    ui.ID,
			Count: ui.Count,
		}.UpdateUserItem(user)
	}

	message := fmt.Sprintf("🍽 Ты съел 1 %s", ui.Item.View)

	return message
}

func (ui UserItem) GetFullDescriptionOfUserItem() string {
	userItem := ui.GetUserItem()
	var fullDescriptionUserItem string
	if userItem.Item.IsInventory == true {
		fullDescriptionUserItem = fmt.Sprintf("%s *%s* - %d шт.\n*Сила*: + %d💥\n", userItem.Item.View, userItem.Item.Name, *userItem.Count, *userItem.Item.Damage)
	} else if userItem.Item.IsBackpack == true {
		fullDescriptionUserItem = fmt.Sprintf("%s *%s* - %dшт.\n*Здоровье*: +%d ♥️️\n*Сытость*: +%d  \U0001F9C3\n", userItem.Item.View, userItem.Item.Name, *userItem.Count, *userItem.Item.Healing, *userItem.Item.Satiety)
	}
	itemDescription := "Описания нет("

	if userItem.Item.Description != nil {
		itemDescription = fmt.Sprintf("*Описание*: %s", *userItem.Item.Description)
	}

	return fullDescriptionUserItem + itemDescription
}

func UpdateUserInstrument(update tgbotapi.Update, user User, instrument Item) (string, string) {
	userItem := UserItem{ItemId: int(instrument.ID), UserId: int(user.ID)}.GetUserItem()

	c := *userItem.CountUseLeft - 1
	if c > 0 {
		UserItem{ID: userItem.ID, CountUseLeft: &c}.UpdateUserItem(user)
		return "Ok", "Ok"
	}

	zeroValue := 0

	if *userItem.Count > 1 {
		userItemCount := *userItem.Count - 1
		countUseLeft := userItem.Item.CountUse
		UserItem{
			ID:           userItem.ID,
			CountUseLeft: countUseLeft,
			Count:        &userItemCount,
		}.UpdateUserItem(user)
	} else {
		UserItem{
			ID:           userItem.ID,
			CountUseLeft: &zeroValue,
			Count:        &zeroValue,
		}.UpdateUserItem(user)

		if user.LeftHandId != nil && *user.LeftHandId == int(userItem.Item.ID) {
			SetNullUserField(update, "left_hand_id")
		}
		if user.RightHandId != nil && *user.RightHandId == int(userItem.Item.ID) {
			SetNullUserField(update, "right_hand_id")
		}
	}

	return "Not ok", fmt.Sprintf("\n💥 Инструмент «%s %s» был сломан! 💥", userItem.Item.View, userItem.Item.Name)
}
