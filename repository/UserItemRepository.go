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
	resUser := GetOrCreateUser(update)
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

func GetUserItems(userId uint, itemTypes ...string) []UserItem {
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

	var resultItemType []UserItem

	for _, userItem := range result {
		for _, itemType := range itemTypes {
			if userItem.Item.Type == itemType {
				resultItemType = append(resultItemType, userItem)
			}
		}
	}

	return resultItemType
}

func UpdateUserItem(user User, userItem UserItem) {
	err := config.Db.Where(UserItem{UserId: int(user.ID), ID: userItem.ID}).Updates(userItem).Error
	if err != nil {
		panic(err)
	}
}

func AddUserItemCount(update tgbotapi.Update, userItem UserItem, cellule Cellule, updateUserMoney int) {
	resUser := GetOrCreateUser(update)
	userId := int(resUser.ID)

	sumCount := *userItem.Count + 1

	err := config.Db.
		//Preload("Item").
		Where(UserItem{UserId: userId, ItemId: *cellule.ItemID}).
		Updates(UserItem{Count: &sumCount}).
		Error
	if err != nil {
		panic(err)
	}

	UpdateUser(update, User{Money: &updateUserMoney})

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
