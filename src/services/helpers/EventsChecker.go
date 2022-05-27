package helpers

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	"project0/src/models"
	t "time"
)

func CheckEventsForUpdate() {
	go UpdateFiredChats()
	go UpdateCellWithNextStateTime()
}

func UpdateFiredChats() {
	var results []models.Chat
	config.Db.
		Where("expired_at <= ?", t.Now()).
		Where("deleted", false).
		Find(&results)

	for _, chat := range results {
		go UpdateCellWithFiredChat(chat)
		chatUser := chat.GetChatUsers()
		NotifyUsers(chatUser, v.GetString("main_info.message_chat_is_closed"))
		go chat.DeleteChatUser()
		go chat.DeleteChat()
	}
}

func UpdateCellWithNextStateTime() {
	var results []models.Cell
	err := config.Db.
		Preload("Item").
		Preload("PrevItem").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.Result").
		Preload("Item.Instruments.NextStageItem").
		Where("next_state_time <= ?", t.Now()).
		Find(&results).
		Error
	if err != nil || len(results) == 0 {
		fmt.Println("Ничего не найдено для обновления!")
	}

	for _, result := range results {
		for _, instrument := range result.Item.Instruments {
			if instrument.Type == "growing" {
				result.UpdateCellAfterGrowing(instrument)
			}
		}
	}
}

func UpdateCellWithFiredChat(chat models.Chat) {
	chatId := int(chat.ID)
	var results []models.Cell

	config.Db.Where(models.Cell{ChatId: &chatId}).Find(&results)

	for _, cell := range results {
		cell.UpdateCellIfChatIsTimeout()
	}
}
