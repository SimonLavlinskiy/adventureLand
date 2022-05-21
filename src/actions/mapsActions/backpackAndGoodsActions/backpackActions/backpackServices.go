package backpackActions

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/userItemController"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"strings"
)

func ListOfBackpackItems(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	if len(charData) == 1 {
		msg, buttons = BackpackCategoryKeyboard()
	} else {
		resUserItems := userItemController.GetBackpackItems(user.ID, charData[1])
		msg = MessageBackpackUserItems(user, resUserItems, 0, charData[1])
		buttons = BackpackInlineKeyboard(resUserItems, 0, charData[1])
	}
	return msg, buttons
}

func BackPackMoving(charData []string, user models.User) (msgText string, buttons tg.InlineKeyboardMarkup) {
	updatedUser := repositories.GetUser(models.User{ID: user.ID})

	category := charData[2]
	userItems := userItemController.GetBackpackItems(updatedUser.ID, category)

	var i int
	switch charData[1] {
	case "-1":
		i = len(userItems) - 1
	case fmt.Sprintf("%d", len(userItems)):
		i = 0
	default:
		i = helpers.ToInt(charData[1])
	}

	msgText = MessageBackpackUserItems(updatedUser, userItems, i, category)
	buttons = BackpackInlineKeyboard(userItems, i, category)

	return msgText, buttons
}

func MessageBackpackUserItems(user models.User, userItems []models.UserItem, rowUser int, itemType string) string {
	var userItemMsg = fmt.Sprintf("%s\n🎒*Рюкзачок* ➡️ *%s* \n\n", userMapController.GetStatsLine(user), v.GetString(fmt.Sprintf("user_location.item_categories.%s", itemType)))

	if len(userItems) == 0 {
		return userItemMsg + "👻 \U0001F9B4  Пусто .... 🕸 🕷"
	}

	for i, item := range userItems {
		var firstCell string
		switch rowUser {
		case i:
			firstCell += item.User.Avatar
		case i + 1, i - 1:
			firstCell += "◻️"
		case i + 2, i - 2:
			firstCell += "◽️️"
		default:
			firstCell += "▫️"
		}
		switch itemType {
		case "food":
			userItemMsg += fmt.Sprintf("%s   %d%s     *HP*:  _%d_ ♥️️     *ST*:  _%d_ \U0001F9C3 ️\n", firstCell, *item.Count, item.Item.View, *item.Item.Healing, *item.Item.Satiety)
		case "resource", "sprout", "furniture":
			userItemMsg += fmt.Sprintf("%s   %s %d шт. - _%s_\n", firstCell, item.Item.View, *item.Count, item.Item.Name)
		default:
			userItemMsg += fmt.Sprintf("%s   %s %d шт.\n", firstCell, item.Item.View, *item.Count)
		}
	}

	return userItemMsg
}

func UserEatItem(user models.User, charData []string) (msgText string, buttons tg.InlineKeyboardMarkup) {
	userItemId := helpers.ToInt(charData[1])

	userItem := models.UserItem{ID: userItemId}.UserGetUserItem()

	res := userItemController.EatItem(userItem, user)
	charDataForOpenBackPack := strings.Fields(fmt.Sprintf("%s %s food", v.GetString("callback_char.backpack_moving"), charData[2]))
	msgText, buttons = BackPackMoving(charDataForOpenBackPack, user)
	msgText = fmt.Sprintf("%s%s%s", msgText, v.GetString("msg_separator"), res)

	return msgText, buttons
}
