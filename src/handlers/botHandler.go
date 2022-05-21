package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	helpers2 "project0/src/services/helpers"
	"project0/src/services/notificationUserChat"
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
		helpers2.NotifyUsers(notificationUserChat.SendUserMessageAllChatUsers(update))
	}
}

func messageHandler(bot *tg.BotAPI, update tg.Update) {

	helpers2.CheckEventsForUpdate()

	if update.CallbackQuery != nil {
		msg, buttons, sendLikeNewMsg := callBackResolver(update)
		if !sendLikeNewMsg {
			helpers2.UpdateMessage(msg, buttons, bot)
		} else {
			newMsg := tg.NewMessage(update.CallbackQuery.From.ID, msg.Text)
			newMsg.ReplyMarkup = buttons.ReplyMarkup
			newMsg.ParseMode = "markdown"
			helpers2.SendMessage(newMsg, bot)

			deletedMsg := tg.DeleteMessageConfig{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.MessageID}
			helpers2.DeleteMessage(deletedMsg, bot)
		}
	}

	if update.Message != nil {
		msg := messageResolver(update)
		helpers2.SendMessage(msg, bot)
	}
}
