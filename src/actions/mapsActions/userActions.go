package mapsActions

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/actions/mapsActions/movingActions"
	"project0/src/actions/mapsActions/userGetBoxAction"
	"project0/src/controllers/instrumentController"
	"project0/src/controllers/itemController"
	"project0/src/controllers/userItemController"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/services/helpers"
)

func CheckUserActions(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup, err error) {
	_, _, ItemHead := UsersHandItemsView(user)

	switch charData[0] {
	// Действия/кнопки  на карте
	case "move":
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		msg, buttons = movingActions.UserMoving(user, cell)
	case "chooseInstrument":
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		msg, buttons = userTouchItem(user, cell)
	case user.Avatar:
		msg, buttons = mapWithUserInfo(user)
	case "box":
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		msg, buttons = userGetBoxAction.UserGetBox(user, cell)

	// Использование надетых итемов
	case "hand", "fist", "item":
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		msg, buttons = useHandOrInstrument(user, charData, cell)
	case "step":
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		msg, buttons = movingActions.UserMoving(user, cell)
	case "head":
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		msg, buttons = userHeadItem(user, cell, ItemHead)

	default:
		err = errors.New("not user actions")
	}

	return msg, buttons, err
}

func userTouchItem(user models.User, cell models.Cell) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg, _ = userMapController.GetMyMap(user)
	buttons = instrumentController.ChooseInstrumentKeyboard(cell, user)

	return msg, buttons
}

func userHeadItem(user models.User, cell models.Cell, ItemHead models.Item) (msg string, buttons tg.InlineKeyboardMarkup) {
	text, err := userItemController.UpdateUserInstrument(user, ItemHead)
	msg = itemController.ViewItemInfo(cell)
	if err != nil {
		msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), text)
	}
	buttons = helpers.CancelButton()
	return msg, buttons
}

func useHandOrInstrument(user models.User, charData []string, cell models.Cell) (msg string, buttons tg.InlineKeyboardMarkup) {
	resultOfGetItem := itemController.UserGetItemUpdateModels(user, cell, charData)

	msgMap, mapButton := userMapController.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msgMap, v.GetString("msg_separator"), resultOfGetItem)

	newCell := models.Cell{ID: cell.ID}.GetCell()

	if newCell.ItemCount == nil || *newCell.ItemCount < 1 {
		return msg, mapButton
	}

	buttons = instrumentController.ChooseInstrumentKeyboard(newCell, user)

	return msg, buttons
}

func mapWithUserInfo(user models.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg, buttons = userMapController.GetMyMap(user)
	msg = fmt.Sprintf("%s\n\n%s", user.GetUserInfo(), msg)
	return msg, buttons
}

func UsersHandItemsView(user models.User) (models.Item, models.Item, models.Item) {
	ItemLeftHand := models.Item{View: v.GetString("message.emoji.hand")}
	ItemRightHand := models.Item{View: v.GetString("message.emoji.hand")}
	var ItemHead models.Item

	if user.LeftHand != nil {
		ItemLeftHand = *user.LeftHand
	}
	if user.RightHand != nil {
		ItemRightHand = *user.RightHand
	}
	if user.Head != nil {
		ItemHead = *user.Head
	}

	return ItemLeftHand, ItemRightHand, ItemHead
}
