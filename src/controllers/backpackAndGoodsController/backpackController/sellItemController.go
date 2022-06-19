package backpackController

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/backpackAndGoodsController/goodsController"
	"project0/src/models"
	"project0/src/repositories"
	"strconv"
	"strings"
)

func SellItem(charData []string) (msg string, buttons tgbotapi.InlineKeyboardMarkup) {
	fmt.Println("char - ", charData)

	userItemId, _ := strconv.Atoi(charData[1])
	userItem := models.UserItem{ID: userItemId}
	userItem = userItem.GetOrCreateUserItem()

	sellCount, _ := strconv.Atoi(charData[3])

	*userItem.Count = *userItem.Count - sellCount
	userItem.User.UpdateUserItem(models.UserItem{ID: userItem.ID, Count: userItem.Count})

	money := sellCount * *userItem.Item.Cost / 2
	*userItem.User.Money = *userItem.User.Money + money
	repositories.UpdateUser(userItem.User)

	msgText := fmt.Sprintf("%s*Продано* %d %s\n*Получено*: _%d_ 💰", v.GetString("msg_separator"), sellCount, userItem.Item.View, money)

	var charDataForOpenList []string

	if charData[3] == "good" {
		goodsController.UserTakeOffClothes(userItem.User, charData)

		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		msg, buttons = goodsController.GoodsMoving(charDataForOpenList, userItem.User)
	} else {
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), charData[2], charData[4]))

		msg, buttons = BackPackMoving(charDataForOpenList, userItem.User)
	}

	msg = fmt.Sprintf("%s%s", msg, msgText)

	return msg, buttons
}
