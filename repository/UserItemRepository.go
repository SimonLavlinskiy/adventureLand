package repository

import (
	"errors"
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	"strings"
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

func (ui UserItem) UserGetUserItem() UserItem {
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

func GetOrCreateUserItem(user User, item Item) UserItem {
	userId := int(user.ID)
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

func (u User) UserGetResultItem(r Result) {
	ui := UserItem{UserId: int(u.ID), ItemId: int(*r.ItemId)}.UserGetUserItem()
	resItemCount := *ui.Count + int(*r.CountItem)
	u.UpdateUserItem(UserItem{ID: ui.ID, Count: &resItemCount})
}

func (u User) UserGetResultSpecialItem(r Result) {
	ui := UserItem{UserId: int(u.ID), ItemId: int(*r.SpecialItemId)}.UserGetUserItem()
	resItemCount := *ui.Count + int(*r.SpecialItemCount)
	u.UpdateUserItem(UserItem{Count: &resItemCount})
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

func GetUserItemsByType(userId uint, itemType []string) []UserItem {
	var userItems []UserItem

	err := config.Db.
		Preload("Item").
		Preload("User").
		Where(UserItem{UserId: int(userId)}).
		Where("count > 0").
		Find(&userItems).
		Error
	if err != nil {
		panic(err)
	}

	var result []UserItem

	if len(itemType) == 1 && itemType[0] == "other" {
		itemType = strings.Fields(v.GetString("user_location.item_categories.other_categories"))
	}

	for _, i := range userItems {
		for y := range itemType {
			if i.Item.Type == itemType[y] {
				result = append(result, i)
			}
		}
	}

	return result
}

func GetBackpackItems(userId uint, itemType ...string) []UserItem {
	userItems := GetUserItemsByType(userId, itemType)

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

func (u User) UpdateUserItem(ui UserItem) {
	err := config.Db.Where(UserItem{UserId: int(u.ID), ID: ui.ID}).Updates(ui).Error
	if err != nil {
		panic(err)
	}
}

func (ui UserItem) EatItem(user User) string {
	userItemCount := ui.Count

	if *userItemCount > 0 {
		itemHeal := uint(*ui.Item.Healing)
		itemSatiety := uint(*ui.Item.Satiety)
		itemDamage := uint(*ui.Item.Damage)

		*ui.Count = *ui.Count - 1
		user.Health = user.Health + itemHeal - itemDamage
		user.Satiety = user.Satiety + itemSatiety

		User{
			TgId:    user.TgId,
			Health:  user.Health,
			Satiety: user.Satiety,
		}.UpdateUser()

		user.UpdateUserItem(UserItem{
			ID:    ui.ID,
			Count: ui.Count,
		})
	}

	message := fmt.Sprintf("🍽 Ты съел 1 %s", ui.Item.View)

	return message
}

func (ui UserItem) GetFullDescriptionOfUserItem() string {
	userItem := ui.UserGetUserItem()
	var fullDescriptionUserItem string

	switch userItem.Item.Type {
	case "food":
		fullDescriptionUserItem = fmt.Sprintf("%s *%s* - %dшт.\n*Здоровье*: +%d ♥️️\n*Сытость*: +%d  \U0001F9C3\n", userItem.Item.View, userItem.Item.Name, *userItem.Count, *userItem.Item.Healing, *userItem.Item.Satiety)
	case "resource", "sprout", "furniture":
		fullDescriptionUserItem = fmt.Sprintf("%s *%s* - %dшт.\n", userItem.Item.View, userItem.Item.Name, *userItem.Count)
	}

	if userItem.Item.IsInventory == true {
		fullDescriptionUserItem = fmt.Sprintf("%s *%s* - %d шт.\n*Сила*: + %d💥\n", userItem.Item.View, userItem.Item.Name, *userItem.Count, *userItem.Item.Damage)
	}

	itemDescription := "Описания нет("

	if userItem.Item.Description != nil {
		itemDescription = fmt.Sprintf("*Описание*: %s", *userItem.Item.Description)
	}

	return fullDescriptionUserItem + itemDescription
}

func UpdateUserInstrument(user User, instrument Item) (string, error) {
	userItem := UserItem{ItemId: int(instrument.ID), UserId: int(user.ID)}.UserGetUserItem()

	var c int

	if userItem.CountUseLeft != nil {
		c = *userItem.CountUseLeft - 1
	} else {
		c = *userItem.Item.CountUse - 1
	}

	if c > 0 {
		user.UpdateUserItem(UserItem{ID: userItem.ID, CountUseLeft: &c})
		return "Ok", nil
	}

	zeroValue := 0

	if *userItem.Count > 1 {
		userItemCount := *userItem.Count - 1
		countUseLeft := userItem.Item.CountUse
		user.UpdateUserItem(UserItem{
			ID:           userItem.ID,
			CountUseLeft: countUseLeft,
			Count:        &userItemCount,
		})
	} else {
		user.UpdateUserItem(UserItem{
			ID:           userItem.ID,
			CountUseLeft: &zeroValue,
			Count:        &zeroValue,
		})

		if user.LeftHandId != nil && *user.LeftHandId == int(userItem.Item.ID) {
			SetNullUserField(user, "left_hand_id")
		}
		if user.RightHandId != nil && *user.RightHandId == int(userItem.Item.ID) {
			SetNullUserField(user, "right_hand_id")
		}
	}

	return fmt.Sprintf("\n💥 Инструмент «%s %s» был сломан! 💥\n_Осталось: %d шт_.", userItem.Item.View, userItem.Item.Name, *userItem.Count), errors.New("item is broken")
}
