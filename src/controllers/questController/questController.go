package questController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/actionsCounterController"
	"project0/src/controllers/resultController"
	"project0/src/models"
	"project0/src/services/helpers"
	"time"
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

	DoneQuest(user, userQuest)
	resultController.UserGetResult(user, userQuest.Quest.Result)

	questResult := resultController.UserGetResultMsg(userQuest.Quest.Result)

	msgText, buttons = OpenQuest(questId, user)
	msgText = fmt.Sprintf("*Задание выполнено!*\n%s%s%s", msgText, v.GetString("msg_separator"), questResult)

	return msgText, buttons
}

func DoneQuest(user models.User, uq models.UserQuest) {
	task := uq.Quest.Task

	switch task.Type {
	case "haveItem":
		ui := models.UserItem{ItemId: *task.ItemId, UserId: int(user.ID)}.GetOrCreateUserItem()
		countItemResult := *ui.Count - *task.CountItem
		user.UpdateUserItem(models.UserItem{ID: ui.ID, Count: &countItemResult})
	case "userLocation":
		uLoc := user.GetUserLocation()
		if uLoc.MapsId == task.MapId && uLoc.AxisX == task.UserAxisX && uLoc.AxisY == task.UserAxisY {
			return
		}
	}

	t := time.Now()
	uq.Status = "completed"
	uq.DoneAt = &t
	uq.UpdateUserQuest()

	actionsCounterController.UserDo(user, "done_quest")

	return
}
