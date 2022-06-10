package eventChecker

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	"project0/src/controllers/itemCellController"
	"project0/src/models"
	"project0/src/services/helpers"
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
		helpers.NotifyUsers(chatUser, v.GetString("main_info.message_chat_is_closed"))
		go chat.DeleteChatUser()
		go chat.DeleteChat()
	}
}

func UpdateCellWithNextStateTime() {
	var results []models.ItemCell
	err := config.Db.
		Preload("Item").
		Preload("Item.ItemAfterBreaking").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.Result").
		Preload("Item.Instruments.NextStageItem").
		Preload("Item.Instruments.GrowingItem").
		Preload("Item.ItemAfterBreaking").
		Preload("ContainedItem").
		Preload("ContainedItem.ItemAfterBreaking").
		Preload("ContainedItem.Instruments").
		Preload("ContainedItem.Instruments.Good").
		Preload("ContainedItem.Instruments.Result").
		Preload("ContainedItem.Instruments.NextStageItem").
		Preload("ContainedItem.Instruments.GrowingItem").
		Where("contained_item_broken_time <= ?", t.Now()).
		Or("growing_time <= ?", t.Now()).
		Or("broken_time <= ?", t.Now()).
		Find(&results).
		Error
	if err != nil || len(results) == 0 {
		fmt.Println("Ничего не найдено для обновления!")
	}

	for _, result := range results {

		growTime := result.GrowingTime
		brokenTime := result.BrokenTime
		containedNextTime := result.ContainedItemBrokenTime

		if growTime != nil && growTime.Before(t.Now()) {
			for _, instrument := range result.Item.Instruments {
				if instrument.Type == "growing" {
					itemCellController.UpdateItemCellAfterGrowing(result, instrument)
				}
			}
		}
		if brokenTime != nil && brokenTime.Before(t.Now()) {
			itemCellController.UpdateItemCellAfterBreaking(result)
		}
		if containedNextTime != nil && containedNextTime.Before(t.Now()) {
			result.UpdateContainedItemCellAfterBreaking()
		}
	}
}

func UpdateCellWithFiredChat(chat models.Chat) {
	chatId := int(chat.ID)
	var results []models.Cell

	config.Db.
		Preload("ItemCell").
		Where(models.Cell{ChatId: &chatId}).
		Find(&results)

	for _, cell := range results {
		cell.UpdateCellIfChatIsTimeout()
	}
}
