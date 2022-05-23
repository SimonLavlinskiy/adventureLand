package instrumentController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
)

func GetInstrumentsUserCanUse(user models.User, cell models.Cell) map[string]models.Item {
	instrumentsUserCanUse := map[string]models.Item{}

	if cell.Item == nil {
		return instrumentsUserCanUse
	}
	instruments := cell.Item.Instruments

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

	if cell.Item.CanTake {
		instrumentsUserCanUse["👋"] = models.Item{Type: "hand"}
	}
	if cell.Item.CanStep && *cell.Type != "swap" {
		instrumentsUserCanUse["👣"] = models.Item{Type: "step"}
	}

	return instrumentsUserCanUse
}

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
