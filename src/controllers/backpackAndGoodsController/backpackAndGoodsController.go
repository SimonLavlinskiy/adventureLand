package backpackAndGoodsController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/backpackAndGoodsController/backpackController"
	"project0/src/controllers/backpackAndGoodsController/goodsController"
	"project0/src/controllers/cellController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"strings"
)

func UserWantsToThrowOutItem(user models.User, charData []string) (msgText string, buttons tg.InlineKeyboardMarkup) {
	userItem := models.UserItem{ID: helpers.ToInt(charData[1])}.UserGetUserItem()

	if userItem.CountUseLeft != nil && *userItem.CountUseLeft != *userItem.Item.CountUse {
		*userItem.Count = *userItem.Count - 1
	}

	if *userItem.Count == 0 {
		var charDataForOpenList []string
		if charData[3] == "good" {
			charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
			if *userItem.CountUseLeft == *userItem.Item.CountUse {
				goodsController.UserTakeOffGood(user, charData)
			}
			msgText, buttons = goodsController.GoodsMoving(charDataForOpenList, user)
		} else {
			charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), charData[2], charData[3]))
			msgText, buttons = backpackController.BackPackMoving(charDataForOpenList, user)
		}
		msgText = fmt.Sprintf("%s%sНельзя выкинуть на карту предмет, который уже был использован!", msgText, v.GetString("msg_separator"))
	} else {
		buttons = CountItemUserWantsToThrowKeyboard(charData, userItem)
		msgText = fmt.Sprintf("%sСколько %s ты хочешь скинуть на карту?", v.GetString("msg_separator"), userItem.Item.View)
	}

	return msgText, buttons
}

func CountItemUserWantsToThrowKeyboard(buttonData []string, userItem models.UserItem) tg.InlineKeyboardMarkup {
	maxCountItem := *userItem.Count
	var buttons [][]tg.InlineKeyboardButton

	for x := 1; x < 10; x = x + 5 {
		var row []tg.InlineKeyboardButton
		if x > maxCountItem {
			break
		}
		for y := 0; y < 5; y++ {
			if x+y > maxCountItem {
				break
			}
			row = append(row, tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", x+y),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.throw_out_item"), buttonData[1], buttonData[2], x+y, buttonData[3])),
			)
		}
		buttons = append(buttons, row)
	}

	for y := 20; y <= maxCountItem; y = y + 20 {
		var row []tg.InlineKeyboardButton
		if y <= maxCountItem {
			row = append(row, tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", y),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.throw_out_item"), buttonData[1], buttonData[2], y, buttonData[3])),
			)
		}
		x := y + 10
		if x <= maxCountItem {
			row = append(row, tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", x),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.throw_out_item"), buttonData[1], buttonData[2], x, buttonData[3])),
			)
		}
		buttons = append(buttons, row)
	}

	// Кнопка Всё
	var row []tg.InlineKeyboardButton
	row = append(row, tg.NewInlineKeyboardButtonData("Все!",
		fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.throw_out_item"), buttonData[1], buttonData[2], maxCountItem, buttonData[3])),
	)
	buttons = append(buttons, row)

	// Кнопка Отмена
	buttons = append(buttons, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData("Отмена", fmt.Sprintf("goodsMoving %s", buttonData[2])),
	))

	return tg.NewInlineKeyboardMarkup(buttons...)
}

func UserDeleteItem(user models.User, charData []string) (msgText string, buttons tg.InlineKeyboardMarkup) {
	userItemId := helpers.ToInt(charData[1])
	userItem := models.UserItem{ID: userItemId}.UserGetUserItem()

	if charData[4] == "false" {
		buttons = DeleteItem(charData)
		msgText = fmt.Sprintf("Вы точно хотите уничтожить 🚮 %s %s _(%d шт.)_?", userItem.Item.View, userItem.Item.Name, *userItem.Count)
		return msgText, buttons
	}

	countAfterUserThrowOutItem := 0
	var updateUserItemStruct = models.UserItem{
		ID:    userItemId,
		Count: &countAfterUserThrowOutItem,
	}

	user.UpdateUserItem(updateUserItemStruct)

	var charDataForOpenList []string
	if charData[3] == "good" {
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		goodsController.UserTakeOffGood(user, charData)
		user = repositories.GetUser(models.User{TgId: user.TgId})
		msgText, buttons = goodsController.GoodsMoving(charDataForOpenList, user)
	} else {
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), charData[2], charData[3]))
		msgText, buttons = backpackController.BackPackMoving(charDataForOpenList, user)
	}

	msgText = fmt.Sprintf("%s%s🗑 Вы уничтожили %s%dшт.", msgText, v.GetString("msg_separator"), userItem.Item.View, *userItem.Count)

	return msgText, buttons
}

func UserThrowOutItem(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	cellType := "item"
	userItem := models.UserItem{ID: helpers.ToInt(charData[1])}.UserGetUserItem()

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
			goodsController.UserTakeOffGood(user, charData)
		}
		msg, buttons = goodsController.GoodsMoving(charDataForOpenList, user)
	} else {
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), charData[2], charData[4]))
		msg, buttons = backpackController.BackPackMoving(charDataForOpenList, user)
	}

	msg = fmt.Sprintf("%s%s", msg, msgText)

	return msg, buttons
}

func DescriptionInlineButton(char []string) tg.InlineKeyboardMarkup {
	switch char[3] {
	case "food":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("🍽 1шт", fmt.Sprintf("%s %s %s", v.GetString("callback_char.eat_food"), char[1], char[2])),
				tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %s %s food", v.GetString("callback_char.count_of_throw_out"), char[1], char[2])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %s %s food false", v.GetString("callback_char.delete_item"), char[1], char[2])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s food", v.GetString("callback_char.backpack_moving"), char[2])),
			),
		)
	case "resource":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %s %s %s", v.GetString("callback_char.count_of_throw_out"), char[1], char[2], char[3])),
				tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %s %s %s false", v.GetString("callback_char.delete_item"), char[1], char[2], char[3])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), char[2], char[3])),
			),
		)
	case "sprout":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("👋\U0001F9A0🗺", fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.throw_out_item"), char[1], char[2], 1, char[3])),
				tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %s %s %s false", v.GetString("callback_char.delete_item"), char[1], char[2], char[3])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), char[2], char[3])),
			),
		)
	case "furniture":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("👋\U0001F9A0🗺", fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.throw_out_item"), char[1], char[2], 1, char[3])),
				tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %s %s %s false", v.GetString("callback_char.delete_item"), char[1], char[2], char[3])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), char[2], char[3])),
			),
		)
	case "good":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s good", v.GetString("callback_char.goods_moving"), char[2])),
			),
		)
	default:
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
			),
		)
	}
}

func DeleteItem(char []string) tg.InlineKeyboardMarkup {
	button := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("🚮 Уничтожить! 💥", fmt.Sprintf("%s %s %s %s true", v.GetString("callback_char.delete_item"), char[1], char[2], char[3])),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Отмена", fmt.Sprintf("goodsMoving %s", char[2])),
		),
	)
	return button
}
