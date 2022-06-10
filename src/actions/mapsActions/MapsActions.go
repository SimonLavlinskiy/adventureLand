package mapsActions

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/src/models"
	"strings"
)

func MapsActions(user models.User, char string) (string, tg.InlineKeyboardMarkup) {
	charData := strings.Fields(char)

	if msg, buttons, err := CheckUserActions(user, charData); err == nil {
		return msg, buttons
	}
	if msg, buttons, err := CheckCellEmojiAction(user, charData); err == nil {
		return msg, buttons
	}
	if msg, buttons, err := CheckBackpackAndGoodsAction(user, charData); err == nil {
		return msg, buttons
	}
	if msg, buttons, err := CheckWrenchActions(user, charData); err == nil {
		return msg, buttons
	}
	if msg, buttons, err := CheckQuestActions(user, charData); err == nil {
		return msg, buttons
	}
	if msg, buttons, err := CheckEventUserActions(user, charData); err == nil {
		return msg, buttons
	}

	msg, buttons := CheckDefaultActions(user, charData)
	return msg, buttons
}
