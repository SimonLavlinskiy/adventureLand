package boxController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/resultController"
	"project0/src/controllers/userMapController"
	"project0/src/models"
)

func UserGetBox(user models.User, cell models.Cell) (msg string, buttons tg.InlineKeyboardMarkup) {
	var resultMsg string

	for _, instrument := range cell.Item.Instruments {
		if instrument.Type == "get" {
			resultController.UserGetResult(user, *instrument.Result)
			resultMsg = UserGetBoxResultMsg(*instrument.Result)
		}
	}

	models.UserBox{BoxId: cell.Item.ID, UserId: user.ID}.CreateUserBox()

	msg, buttons = userMapController.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), resultMsg)
	return msg, buttons
}

func UserGetBoxResultMsg(result models.Result) string {
	result = result.GetResult()

	msg := "🏆 *Ты получил*:"
	if result.Item != nil {
		msg = fmt.Sprintf("%s\n_%s %s - %d шт._", msg, result.Item.View, result.Item.Name, *result.CountItem)
	}
	if result.SpecialItem != nil {
		msg = fmt.Sprintf("%s\n_%s %s - %d шт._", msg, result.SpecialItem.View, result.SpecialItem.Name, *result.SpecialItemCount)
	}

	return msg
}
