package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/helpers"
)

var deleteBotMsg tg.DeleteMessageConfig

func GetMessage(telegramApiToken string) {
	bot, err := tg.NewBotAPI(telegramApiToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	updateConfig := tg.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		helpers.CheckEventsForUpdate()

		if update.CallbackQuery != nil {
			deleteBotMsg = tg.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			msgs, delMes := callBackResolver(update)
			for i := range msgs {
				helpers.SendMessage(msgs[i], bot)
				if delMes {
					go helpers.DeleteMessage(deleteBotMsg, bot)
				}
			}
		}

		if update.Message != nil {
			msgs = messageResolver(update)
			for _, msg := range msgs {
				helpers.SendMessage(msg, bot)
			}
		}
	}
}

func GetMessageFromChat(tgApiToken string) {
	bot, err := tg.NewBotAPI(tgApiToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	updateConfig := tg.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		helpers.NotifyUsers(SendUserMessageAllChatUsers(update))
	}
}
