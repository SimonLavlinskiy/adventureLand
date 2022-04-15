package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/services"
)

//var deleteBotMsg tg.DeleteMessageConfig

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
		go messageHandler(bot, update)
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

func messageHandler(bot *tg.BotAPI, update tg.Update) {

	services.CheckEventsForUpdate()

	if update.CallbackQuery != nil {
		msg, buttons := callBackResolver(update)

		services.UpdateMessage(msg, buttons, bot)
	}

	if update.Message != nil {
		msg := messageResolver(update)
		services.SendMessage(msg, bot)
	}
}
