package questController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
	str "strings"
	"time"
)

func OpenQuestKeyboard(q models.Quest, uq models.UserQuest) (buttons tg.InlineKeyboardMarkup) {
	switch uq.Status {
	case "new":
		buttons = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData("Взять в работу", fmt.Sprintf("user_get_quest %d", q.ID))),
			quitButton(q),
		)
	case "processed":
		buttonStatus := tg.NewInlineKeyboardButtonData("Еще в работе... Прийду потом", "quests")

		if q.Task.HasUserDoneTask(uq.User) {
			buttonStatus = tg.NewInlineKeyboardButtonData("Готово! Я всё сделаль!", fmt.Sprintf("user_done_quest %d", uq.QuestId))
		}

		buttons = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(buttonStatus), quitButton(q))
	default:
		return tg.NewInlineKeyboardMarkup(quitButton(q))
	}

	return buttons
}

func AllQuestsMessageKeyboard(user models.User, daily bool) tg.InlineKeyboardMarkup {
	list := listOfQuests(daily)
	userQuests := checkOrUpdateUserQuest(user)

	var quests []models.Quest

	if str.Contains(user.MenuLocation, "learning") {
		for _, quest := range list {
			if quest.Type == "learning" {
				quests = append(quests, quest)
			}
		}
	} else {
		quests = list
	}

	if len(quests) == 0 {
		return tg.NewInlineKeyboardMarkup(
			questMenuLine(),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "cancel"),
			),
		)
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
		if daily && uq.Quest.Timeout != nil {
			m[uq.QuestId] = statusQuest{status: uq.Status, quest: uq.Quest}
		} else if !daily && uq.Quest.Timeout == nil {
			m[uq.QuestId] = statusQuest{status: uq.Status, quest: uq.Quest}
		}
	}

	var result [][]tg.InlineKeyboardButton
	result = append(result, questMenuLine())

	for _, i := range m {
		status := v.GetString(fmt.Sprintf("quest_statuses.%s_emoji", i.status))
		result = append(result, questButton(status, i.quest.Name, i.quest.ID))
	}

	result = append(result, tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData("Выйти", "cancel")))

	return tg.NewInlineKeyboardMarkup(result...)
}

func questMenuLine() []tg.InlineKeyboardButton {
	return tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData("📜 Квесты", "quests"),
		tg.NewInlineKeyboardButtonData("📆 Ежедневки", "dailyQuests"),
	)
}

func questButton(status string, name string, questId uint) []tg.InlineKeyboardButton {
	return tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s - Задание: «%s»", status, name),
			fmt.Sprintf("quest %d", questId),
		),
	)
}

func listOfQuests(daily bool) []models.Quest {
	if daily {
		return models.Quest{}.GetDailyQuests()
	} else {
		return models.Quest{}.GetQuests()
	}
}

func quitButton(quest models.Quest) []tg.InlineKeyboardButton {
	if quest.Timeout != nil {
		return tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Назад", "dailyQuests"),
		)
	} else {
		return tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Назад", "quests"),
		)
	}
}

func checkOrUpdateUserQuest(user models.User) []models.UserQuest {
	quests := models.User{ID: user.ID}.GetUserQuests()

	_, week := time.Now().ISOWeek()
	day := time.Now().Day()
	for _, quest := range quests {
		if quest.Status != "completed" || quest.Quest.Timeout == nil {
			continue
		}

		userQuestDay := quest.DoneAt.Day()
		_, userQuestWeek := quest.DoneAt.ISOWeek()
		if *quest.Quest.Timeout == "weekly" && userQuestWeek != week {
			quest.Status = "new"
			quest.UpdateUserQuestStatus()
		}
		if *quest.Quest.Timeout == "daily" && userQuestDay != day {
			quest.Status = "new"
			quest.UpdateUserQuestStatus()
		}
	}

	return models.User{ID: user.ID}.GetUserQuests()
}
