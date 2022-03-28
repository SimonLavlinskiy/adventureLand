package helpers

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	r "project0/repository"
	t "time"
)

func CheckEventsForUpdate() {
	go UpdateFiredChats()
	go UpdateCellWithNextStateTime()
}

func UpdateFiredChats() {
	var results []r.Chat
	config.Db.
		Where("expired_at <= ?", t.Now()).
		Where("deleted", false).
		Find(&results)

	if len(results) != 0 {
		for _, chat := range results {
			go r.UpdateCellWithFiredChat(chat)
			chatUser := chat.GetChatUsers()
			NotifyUsers(chatUser, v.GetString("main_info.message_chat_is_closed"))
			go chat.DeleteChatUser()
			go chat.DeleteChat()
		}
	}
}

func UpdateCellWithNextStateTime() {
	var results []r.Cell
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
