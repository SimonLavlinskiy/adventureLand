package questController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/resultController"
	"project0/src/models"
	"project0/src/services/helpers"
)

func OpenQuest(questId uint, user models.User) (msgText string, buttons tg.InlineKeyboardMarkup) {
	quest := models.Quest{ID: questId}.GetQuest()
	checkOrUpdateUserQuest(user)
	userQuest := models.UserQuest{UserId: user.ID, QuestId: questId}.GetUserQuest()

	msgText = quest.QuestInfo(userQuest)
	buttons = OpenQuestKeyboard(quest, userQuest)

	return msgText, buttons
}

func UserDoneQuest(questId uint, user models.User) (msgText string, buttons tg.InlineKeyboardMarkup) {
	userQuest := models.UserQuest{UserId: user.ID, QuestId: questId}.GetUserQuest()
	if !userQuest.Quest.Task.HasUserDoneTask(user) {
		msgText = v.GetString("errors.user_did_not_task")
		return msgText, helpers.CancelButton()
	}

	userQuest.UserDoneQuest(user)
	resultController.UserGetResult(user, userQuest.Quest.Result)

	questResult := resultController.UserGetResultMsg(userQuest.Quest.Result)

	msgText, buttons = OpenQuest(questId, user)
	msgText = fmt.Sprintf("*Задание выполнено!*\n%s%s%s", msgText, v.GetString("msg_separator"), questResult)

	return msgText, buttons
}
