package questController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
	str "strings"
)

func OpenQuestKeyboard(q models.Quest, uq models.UserQuest) tg.InlineKeyboardMarkup {
	switch uq.Status {
	case "":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Взять в работу", fmt.Sprintf("user_get_quest %d", q.ID)),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", "quests"),
			),
		)
	case "processed":
		var buttonStatus tg.InlineKeyboardButton
		if q.Task.HasUserDoneTask(uq.User) {
			buttonStatus = tg.NewInlineKeyboardButtonData("Готово! Я всё сделаль!", fmt.Sprintf("user_done_quest %d", uq.QuestId))
		} else {
			buttonStatus = tg.NewInlineKeyboardButtonData("Еще в работе... Прийду потом", "quests")
		}

		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				buttonStatus,
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", "quests"),
			),
		)
	default:
		return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Назад", "quests"),
		))
	}
}

func AllQuestsMessageKeyboard(u models.User) tg.InlineKeyboardMarkup {
	listOfQuests := models.Quest{}.GetQuests()
	userQuests := models.User{ID: u.ID}.GetUserQuests()

	var quests []models.Quest

	if str.Contains(u.MenuLocation, "learning") {
		for _, quest := range listOfQuests {
			if quest.Type == "learning" {
				quests = append(quests, quest)
			}
		}
	} else {
		quests = listOfQuests
	}

	type statusQuest struct {
		status string
		quest  models.Quest
	}

	m := map[uint]statusQuest{}
	for _, quest := range quests {
		m[quest.ID] = statusQuest{status: "new", quest: quest}
	}

	for _, uq := range userQuests {
		m[uq.QuestId] = statusQuest{status: uq.Status, quest: uq.Quest}
	}

	if len(quests) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "cancel"),
			),
		)
	}

	var result [][]tg.InlineKeyboardButton

	for _, i := range m {
		status := v.GetString(fmt.Sprintf("quest_statuses.%s_emoji", i.status))
		result = append(result,
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s - Задание: «%s»", status, i.quest.Name),
					fmt.Sprintf("quest %d", i.quest.ID),
				),
			),
		)
	}

	result = append(result, tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData("Выйти", "cancel")))

	return tg.NewInlineKeyboardMarkup(result...)
}
