package wordleController

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
	"strings"
)

func GameWordle(update tg.Update, user models.User) (msgText string, btns tg.InlineKeyboardMarkup) {

	if update.CallbackQuery == nil {
		msgText, btns = UserSendNextWord(user, update.Message.Text)
		return msgText, btns
	}

	charData := strings.Fields(update.CallbackQuery.Data)

	switch charData[0] {
	case "/map":
		msgText, btns = userMapController.GetMyMap(user)
		user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	case v.GetString("callback_char.wordle_regulations"):
		msgText = v.GetString("wordle.regulations")
		btns = wordleExitButton()
	case "wordleUserStatistic":
		msgText = GetWordleUserStatistic(user)
		btns = wordleExitButton()
	case "wordleMenu":
		msgText, btns = WordleMap(user)
	}

	return msgText, btns
}

func wordleExitButton() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("В меню...", "wordleMenu"),
		),
	)
}
