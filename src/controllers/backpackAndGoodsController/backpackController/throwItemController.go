package backpackController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/backpackAndGoodsController/goodsController"
	"project0/src/controllers/cellController"
	"project0/src/models"
	"project0/src/services/helpers"
	"strings"
)

func UserThrowOutItem(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	cellType := "item"
	userItem := models.UserItem{ID: helpers.ToInt(charData[1])}.GetOrCreateUserItem()

	*userItem.Count = *userItem.Count - helpers.ToInt(charData[3])

	var msgText string

	if charData[4] == "other" && userItem.Item.Type == "chat" {
		cellType = "chat"
	}

	err := cellController.UpdateCellUnderUserWhenUserThrowItem(user, userItem, helpers.ToInt(charData[3]), cellType)
	if err != nil {
		msgText = fmt.Sprintf("%s%s", v.GetString("msg_separator"), err)
	} else {
		msgText = fmt.Sprintf("%sВы сбросили %s %sшт. на карту!", v.GetString("msg_separator"), userItem.Item.View, charData[3])
		user.UpdateUserItem(models.UserItem{ID: userItem.ID, Count: userItem.Count})
	}

	var charDataForOpenList []string

	if charData[4] == "good" {

		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))

		if *userItem.Count == 0 {
			goodsController.UserTakeOffClothes(user, charData)
		}

		msg, buttons = goodsController.GoodsMoving(charDataForOpenList, user)
	} else {

		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), charData[2], charData[4]))

		msg, buttons = BackPackMoving(charDataForOpenList, user)
	}

	msg = fmt.Sprintf("%s%s", msg, msgText)

	return msg, buttons
}
