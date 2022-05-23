package mapsActions

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	wrenchController2 "project0/src/controllers/wrenchController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"strings"
)

func CheckWrenchActions(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup, err error) {
	// Крафтинг
	switch charData[0] {
	case "wrench":
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		msg, buttons = wrenchController2.Workbench(&cell, charData)
	case v.GetString("callback_char.workbench"):
		msg, buttons = wrenchController2.Workbench(nil, charData)
	case v.GetString("callback_char.receipt"):
		msg, buttons = listOfReceipt(charData)
	case v.GetString("callback_char.put_item"):
		msg, buttons = changeComponent(user, charData)
	case v.GetString("callback_char.put_count_item"):
		msg, buttons = changeCountComponent(charData)
	case v.GetString("callback_char.make_new_item"):
		msg, buttons = craftItem(user, charData)
	default:
		err = errors.New("not wrench actions")
	}

	return msg, buttons, err
}

func craftItem(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	resp := GetReceiptFromData(charData)
	receipt := repositories.FindReceiptForUser(resp)
	msg, buttons = wrenchController2.UserCraftItem(user, receipt, charData)
	return msg, buttons
}

func GetReceiptFromData(char []string) models.Receipt {
	var result models.Receipt

	if char[4] != "nil" && char[5] != "0" {
		fstItem := models.UserItem{ID: helpers.ToInt(char[4])}.UserGetUserItem()
		id := int(fstItem.Item.ID)
		c := helpers.ToInt(char[5])

		result.Component1ID = &id
		result.Component1Count = &c
	}

	if char[7] != "nil" && char[8] != "0" {
		fstItem := models.UserItem{ID: helpers.ToInt(char[7])}.UserGetUserItem()
		id := int(fstItem.Item.ID)
		c := helpers.ToInt(char[8])

		result.Component2ID = &id
		result.Component2Count = &c
	}

	if char[10] != "nil" && char[11] != "0" {
		fstItem := models.UserItem{ID: helpers.ToInt(char[10])}.UserGetUserItem()
		id := int(fstItem.Item.ID)
		c := helpers.ToInt(char[11])

		result.Component3ID = &id
		result.Component3Count = &c
	}

	return result
}

func changeCountComponent(charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	buttons = wrenchController2.PutCountComponent(charData)
	msg = fmt.Sprintf("%s%s⚠️ Сколько выкладываешь?", wrenchController2.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
	return msg, buttons
}

func changeComponent(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	userItem := repositories.GetUserItemsByType(user.ID, strings.Fields("food resource"))
	buttons = ChooseUserItemKeyboard(userItem, charData)
	msg = fmt.Sprintf("%s%sВыбери предмет:", wrenchController2.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
	return msg, buttons
}

func listOfReceipt(charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg = fmt.Sprintf("📖 *Рецепты*: 📖%s%s", v.GetString("msg_separator"), wrenchController2.AllReceiptsMsg())
	buttons = wrenchController2.ReturnToWorkbench(charData)
	return msg, buttons
}

func ChooseUserItemKeyboard(userItem []models.UserItem, char []string) tg.InlineKeyboardMarkup {
	var buttons [][]tg.InlineKeyboardButton

	var itemData string

	for x := 0; x < len(userItem); x = x + 5 {

		var row []tg.InlineKeyboardButton

		for i := 0; i < 5; i++ {
			if i+x < len(userItem) {
				switch char[2] {
				case "0":
					itemData = fmt.Sprintf("%s usPoint %s 1stComp %d %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), char[2], userItem[x+i].ID, char[5], char[7], char[8], char[10], char[11])
				case "1":
					itemData = fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %d %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), char[2], char[4], char[5], userItem[x+i].ID, char[8], char[10], char[11])
				case "2":
					itemData = fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %d %s", v.GetString("callback_char.put_count_item"), char[2], char[4], char[5], char[7], char[8], userItem[x+i].ID, char[11])
				}
				row = append(row, tg.NewInlineKeyboardButtonData(userItem[x+i].Item.View, itemData))
			}
		}
		buttons = append(buttons, row)
	}

	return tg.NewInlineKeyboardMarkup(buttons...)
}
