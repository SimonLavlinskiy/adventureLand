package learningController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/mapController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/menu"
	"strings"
)

func Learning(update tg.Update, user models.User) (text string, buttons tg.InlineKeyboardMarkup) {
	if update.CallbackQuery == nil && user.MenuLocation == "learning" {
		return greetingUser(user)
	}

	if update.CallbackQuery == nil {
		return
	}
	data := update.CallbackQuery.Data

	switch true {
	case strings.Contains(user.MenuLocation, "step1"):
		text, buttons = learningStep1(data, user)
	case strings.Contains(user.MenuLocation, "step2"):
		text, buttons = learningStep2(data, user)
	case strings.Contains(user.MenuLocation, "step3"):
		text, buttons = learningStep3(data, user)
	case strings.Contains(user.MenuLocation, "step4"):
		text, buttons = learningStep4(data, user)
	case strings.Contains(user.MenuLocation, "step5"):
		text, buttons = learningStep5(data, user)
	case strings.Contains(user.MenuLocation, "step6"):
		text, buttons = learningStep6(data, user)
	default:
		if len(data) != 0 {
			text, buttons = startUserAction(data, user)
		} else {
			text, buttons = greetingUser(user)
		}
	}

	return text, buttons
}

func greetingUser(user models.User) (text string, buttons tg.InlineKeyboardMarkup) {
	text = fmt.Sprintf("Здравствуй, *%s %s*! 🤝\n"+
		"Здесь я научу тебя основам, которые помогут не потеряться в кнопках игры!\n\n"+
		"%s - это ты!\n\n"+
		"Ты можешь выбрать себе *новый аватар* или выбрать потом в *Меню > Профиль*", user.FirstName, user.LastName, user.Avatar)
	buttons = startLearningButton()

	return text, buttons
}

func startLearningButton() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Выбрать аватар", "chooseAvatar"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Продолжить", "step1"),
		),
	)
}

func startUserAction(data string, user models.User) (text string, buttons tg.InlineKeyboardMarkup) {
	charData := strings.Fields(data)

	switch true {
	case strings.Contains(data, "chooseAvatar"):
		text = "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар..."
		buttons = menu.EmojiInlineKeyboard()
	case strings.Contains(data, v.GetString("callback_char.change_avatar")):
		res := repositories.UpdateUser(models.User{TgId: user.TgId, Avatar: charData[1]})
		text, buttons = greetingUser(res)
	case strings.Contains(data, "step1"):
		user.MenuLocation = "learning step1"
		repositories.UpdateUser(user)

		text, buttons = mapController.GetMyMap(user)
		text = fmt.Sprintf("Это первая карта, которую я создал, когда начинал писать игру!\n\n"+
			"%s - это ты!\n\n"+
			"*Шаг 1:*\nВидишь снизу кнопки-стрелочки (◀️ 🔼 ▶️ 🔽)? Они позволяют тебе ходить!\n"+
			"Попробуй пройтись по карте, а как освоишься, бери квест \U0001FAA7 на обучение и заходи в дверь 🚪%s%s", user.Avatar, v.GetString("msg_separator"), text)
	}

	return text, buttons
}
