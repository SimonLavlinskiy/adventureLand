package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/repository"
	"project0/services"
	"time"
)

var deleteBotMsg tg.DeleteMessageConfig

func GetMessage(telegramApiToken string) {
	//services.UpdateMap()

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
	var msgs []tg.MessageConfig
	var delMes bool

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
		repository.UserMsgCreate(update)
		msgs = messageResolver(update)
		time.Sleep(500 * time.Millisecond)
		if repository.IsUserMsgLatest(update) {
			for _, msg := range msgs {
				services.SendMessage(msg, bot)
			}
		}
	}
}
