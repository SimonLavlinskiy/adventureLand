package instrumentServices

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
)

func ChooseInstrumentKeyboard(cell models.Cell, user models.User) (mapButton tg.InlineKeyboardMarkup) {
	instruments := GetInstrumentsUserCanUse(user, cell)

	if len(instruments) != 0 {
		var row []tg.InlineKeyboardButton

		for view, instrument := range instruments {
			button := tg.NewInlineKeyboardButtonData(
				getButtonTextAndData(cell, instrument, view, user),
			)
			row = append(row, button)
		}

		mapButton = tg.NewInlineKeyboardMarkup(
			row,
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Отмена", v.GetString("callback_char.cancel")),
			),
		)
	}

	return mapButton
}

func getButtonTextAndData(cell models.Cell, instrument models.Item, instrumentView string, user models.User) (text string, data string) {
	if cell.NeedPay && cell.Item.Cost != nil && *cell.Item.Cost > 0 {
		if (instrument.Type == "hand" || instrument.Type == "swap") && (*cell.Type == "swap" || *cell.Type == "item" && *cell.ItemCount != 0) {
			text = fmt.Sprintf("%s ( %d💰 )", instrumentView, *cell.Item.Cost)
		} else {
			text = instrumentView
		}
	} else {
		text = instrumentView
	}

	if instrument.DressType != nil && *instrument.DressType == "head" {
		data = fmt.Sprintf("head %d %d", cell.ID, instrument.ID)
	} else if instrument.DressType != nil {
		data = fmt.Sprintf("item %d %d", cell.ID, instrument.ID)
	} else if instrument.Type == "fist" {
		data = fmt.Sprintf("fist %d %d", cell.ID, instrument.ID)
	} else if cell.IsBox(user) && instrument.Type == "hand" {
		data = fmt.Sprintf("box %d", cell.ID)
	} else {
		data = fmt.Sprintf("%s %d", instrument.Type, cell.ID)
	}

	return text, data
}
