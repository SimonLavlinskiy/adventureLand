package boxController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/mapController"
	"project0/src/controllers/resultController"
	"project0/src/models"
)

func UserGetBox(user models.User, cell models.Cell) (msg string, buttons tg.InlineKeyboardMarkup) {
	var resultMsg string

	for _, instrument := range cell.ItemCell.Item.Instruments {
		if instrument.Type == "get" {
			resultController.UserGetResult(user, *instrument.Result)
			resultMsg = resultController.UserGetResultMsg(*instrument.Result)
		}
	}

	models.UserBox{BoxId: cell.ItemCell.Item.ID, UserId: user.ID}.CreateUserBox()

	msg, buttons = mapController.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), resultMsg)
	return msg, buttons
}
