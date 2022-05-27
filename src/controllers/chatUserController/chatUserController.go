package chatUserController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/src/models"
)

func OpenChatKeyboard(cell models.Cell, user models.User) (msgText string, buttons tg.InlineKeyboardMarkup) {
	var button tg.InlineKeyboardButton
	msgText = "Присоединяйся и общайся!"

	if !cell.IsChat() {
		msgText = "Здесь нет чата! Поищи в другом месте..."
		button = tg.NewInlineKeyboardButtonData("Жаль...", "cancel")
	} else {
		userChat := cell.Chat.GetChatUser(user)

		if userChat == nil {
			button = tg.NewInlineKeyboardButtonData("Присоединиться к беседе", fmt.Sprintf("joinToChat %d cell %d", *cell.ChatId, cell.ID))
		} else {
			button = tg.NewInlineKeyboardButtonURL("Перейти в беседу", "https://t.me/AdventureChatBot")
		}
	}

	buttons = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			button,
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Назад", "cancel"),
		),
	)
	return msgText, buttons
}
