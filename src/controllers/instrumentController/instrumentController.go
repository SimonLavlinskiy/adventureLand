package instrumentController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
)

func GetInstrumentsUserCanUse(user models.User, cell models.Cell) map[string]models.Item {
	instrumentsUserCanUse := map[string]models.Item{}

	if cell.ItemCell == nil || cell.ItemCell.Item == nil {
		return instrumentsUserCanUse
	}
	instruments := cell.ItemCell.Item.Instruments

	for _, instrument := range instruments {
		if instrument.GoodId != nil && user.LeftHandId != nil && user.LeftHand.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.LeftHand.View] = *user.LeftHand
		}
		if instrument.GoodId != nil && user.RightHandId != nil && user.RightHand.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.RightHand.View] = *user.RightHand
		}
		if instrument.GoodId != nil && user.HeadId != nil && user.Head.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.Head.View] = *user.Head
		}
		if instrument.GoodId != nil && instrument.Good.Type == "fist" {
			instrumentsUserCanUse["🤜"] = *instrument.Good
		}
		if *cell.Type == "item" && instrument.Type == "goToSleep" {
			instrumentsUserCanUse["💤"] = models.Item{Type: "sleep"}
		}
	}

	if cell.ItemCell.Item.CanTake {
		instrumentsUserCanUse["👋"] = models.Item{Type: "hand"}
	}
	if cell.ItemCell.Item.CanStep && *cell.Type != "swap" {
		instrumentsUserCanUse["👣"] = models.Item{Type: "step"}
	}
	if cell.ItemCell.ContainedItemId != nil && cell.ItemCell.ContainedItemCount != nil && *cell.ItemCell.ContainedItemCount > 0 {
		instrumentsUserCanUse[fmt.Sprintf("👋%s", cell.ItemCell.ContainedItem.View)] = models.Item{Type: "takeContain"}
	}

	return instrumentsUserCanUse
}

func ChooseInstrumentKeyboard(cell models.Cell, user models.User) (mapButton tg.InlineKeyboardMarkup) {
	instruments := GetInstrumentsUserCanUse(user, cell)

	if len(instruments) != 0 {
		var rows [][]tg.InlineKeyboardButton
		var row []tg.InlineKeyboardButton

		i := 0

		for view, instrument := range instruments { // todo tut proishodit kakaya to magic, кнопки все время в разном порядке
			i = i + 1

			button := tg.NewInlineKeyboardButtonData(
				getButtonTextAndData(cell, instrument, view, user),
			)

			if i > 5 {
				rows = append(rows, row)
				row = []tg.InlineKeyboardButton{}
			}

			row = append(row, button)
			if len(instruments) == i {
				rows = append(rows, row)
			}
		}

		rows = append(rows, tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData("Отмена", v.GetString("callback_char.cancel"))))

		mapButton = tg.NewInlineKeyboardMarkup(
			rows...,
		)
	}

	return mapButton
}

func getButtonTextAndData(cell models.Cell, instrument models.Item, instrumentView string, user models.User) (text string, data string) {
	if cell.NeedPay && cell.ItemCell.Item.Cost != nil && *cell.ItemCell.Item.Cost > 0 {
		if (instrument.Type == "hand" || instrument.Type == "swap") && (*cell.Type == "swap" || *cell.Type == "item" && *cell.ItemCell.ItemCount != 0) {
			text = fmt.Sprintf("%s ( %d💰 )", instrumentView, *cell.ItemCell.Item.Cost)
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
