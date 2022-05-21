package repositories

import (
	v "github.com/spf13/viper"
	"project0/config"
	"project0/src/models"
	"strings"
)

func GetOrCreateUserItem(user models.User, item models.Item) models.UserItem {
	userId := int(user.ID)
	countItem := 0
	result := models.UserItem{
		Count:        &countItem,
		UserId:       userId,
		ItemId:       int(item.ID),
		CountUseLeft: item.CountUse,
	}
	err := config.Db.
		Preload("Item").
		Preload("Item.Instruments").
		Where(models.UserItem{UserId: userId, ItemId: int(item.ID)}).
		FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func GetUserItems(userId uint) []models.UserItem {
	var result []models.UserItem

	err := config.Db.
		Preload("Item").
		Preload("User").
		Where(models.UserItem{UserId: int(userId)}).
		Where("count > 0").
		Find(&result).
		Error
	if err != nil {
		panic(err)
	}

	return result
}

func GetUserItemsByType(userId uint, itemType []string) []models.UserItem {
	var userItems []models.UserItem

	err := config.Db.
		Preload("Item").
		Preload("User").
		Where(models.UserItem{UserId: int(userId)}).
		Where("count > 0").
		Find(&userItems).
		Error
	if err != nil {
		panic(err)
	}

	var result []models.UserItem

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
