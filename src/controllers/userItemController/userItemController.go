package userItemController

import (
	"errors"
	"fmt"
	"project0/src/models"
	"project0/src/repositories"
)

func GetBackpackItems(userId uint, itemType ...string) []models.UserItem {
	userItems := repositories.GetUserItemsByType(userId, itemType)

	var backpackUserItem []models.UserItem

	for _, userItem := range userItems {
		if userItem.Item.IsBackpack == true {
			backpackUserItem = append(backpackUserItem, userItem)
		}
	}

	return backpackUserItem
}

func GetInventoryItems(userId uint) []models.UserItem {
	userItems := repositories.GetUserItems(userId)

	var inventoryUserItem []models.UserItem

	for _, userItem := range userItems {
		if userItem.Item.IsInventory == true {
			inventoryUserItem = append(inventoryUserItem, userItem)
		}
	}

	return inventoryUserItem
}

func UpdateUserInstrument(user models.User, instrument models.Item) (result string, err error) {
	userItem := models.UserItem{ItemId: int(instrument.ID), UserId: int(user.ID)}.UserGetUserItem()

	newCountUseLeft := *userItem.Item.CountUse - 1

	if userItem.CountUseLeft != nil {
		newCountUseLeft = *userItem.CountUseLeft - 1
	}

	if newCountUseLeft > 0 {
		user.UpdateUserItem(models.UserItem{ID: userItem.ID, CountUseLeft: &newCountUseLeft})
		return "Ok", nil
	}

	userItem = CrushUserItem(user, userItem)
	result = fmt.Sprintf("\n💥 Инструмент «%s %s» был сломан! 💥\n_Осталось: %d шт_.",
		userItem.Item.View, userItem.Item.Name, *userItem.Count)

	return result, errors.New("item is broken")
}

func EatItem(userItem models.UserItem, user models.User) string {
	userItemCount := userItem.Count

	if *userItemCount > 0 {
		itemHeal := uint(*userItem.Item.Healing)
		itemSatiety := uint(*userItem.Item.Satiety)
		itemDamage := uint(*userItem.Item.Damage)

		*userItem.Count = *userItem.Count - 1
		user.Health = user.Health + itemHeal - itemDamage
		user.Satiety = user.Satiety + itemSatiety

		if user.Satiety > 100 {
			user.Satiety = 100
		}

		if user.Health > 100 {
			user.Health = 100
		}

		repositories.UpdateUser(user)
		user.UpdateUserItem(userItem)
	}

	message := fmt.Sprintf("🍽 Ты съел 1 %s", userItem.Item.View)

	return message
}

func CrushUserItem(user models.User, userItem models.UserItem) models.UserItem {
	if *userItem.Count > 1 {
		userItemCount := *userItem.Count - 1
		newUserItem := models.UserItem{
			ID:           userItem.ID,
			Count:        &userItemCount,
			CountUseLeft: userItem.Item.CountUse,
		}
		user.UpdateUserItem(newUserItem)
	} else {
		zero := 0
		newUserItem := models.UserItem{
			ID:           userItem.ID,
			Count:        &zero,
			CountUseLeft: &zero,
		}
		user.UpdateUserItem(newUserItem)

		if user.LeftHandId != nil && *user.LeftHandId == int(userItem.Item.ID) {
			repositories.SetNullUserField(user, "left_hand_id")
		}
		if user.RightHandId != nil && *user.RightHandId == int(userItem.Item.ID) {
			repositories.SetNullUserField(user, "right_hand_id")
		}
	}

	return userItem
}
