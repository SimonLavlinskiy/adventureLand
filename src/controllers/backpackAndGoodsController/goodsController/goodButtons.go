package goodsController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
)

func GoodsInlineKeyboard(user models.User, userItems []models.UserItem, i int) tg.InlineKeyboardMarkup {
	if len(userItems) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyGoods"),
			),
		)
	}

	text, data := user.IsDressedItem(userItems[i])
	itemDescription := "Описания нет("
	if userItems[i].Item.Description != nil {
		itemDescription = "Инфо"
	}

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s - %s %s", itemDescription, userItems[i].Item.View, userItems[i].Item.Name),
				fmt.Sprintf("%s %d %d good", v.GetString("callback_char.description"), userItems[i].ID, i),
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(text, fmt.Sprintf("%s %d %d", data, userItems[i].ID, i)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d", v.GetString("callback_char.goods_moving"), i-1)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d good", v.GetString("callback_char.count_of_throw_out"), userItems[i].ID, i)),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d", v.GetString("callback_char.goods_moving"), i+1)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %d %d good false", v.GetString("callback_char.delete_item"), userItems[i].ID, i)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
		),
	)
}
