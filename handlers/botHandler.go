package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/services"
)

var deleteBotMsg tg.DeleteMessageConfig

func GetMessage(telegramApiToken string) {
	//services.UpdateMap()

	var msgs []tg.MessageConfig
	var delMes bool
	bot, err := tg.NewBotAPI(telegramApiToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	updateConfig := tg.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		services.CheckEventsForUpdate()

		if update.CallbackQuery != nil {
			deleteBotMsg = tg.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			msgs, delMes = callBackResolver(update)
			for i := range msgs {
				services.SendMessage(msgs[i], bot)
				if delMes {
					go services.DeleteMessage(deleteBotMsg, bot)
				}
			}
		}

		if update.Message != nil {
			msgs = messageResolver(update)
			for _, msg := range msgs {
				services.SendMessage(msg, bot)
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
		services.NotifyUsers(SendUserMessageAllChatUsers(update))
	}
}
