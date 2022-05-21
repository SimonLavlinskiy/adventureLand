package sleepUserController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"math"
	"project0/config"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
	"time"
)

func UpdateUserSleepTime(user models.User) {
	t := time.Now()
	sleepTime := models.UserSleep{UserId: user.ID, SleptAt: t}
	repositories.GetOrCreateUserSleep(user, sleepTime)
}

func UserGetHpAfterSleep(user models.User) (result string) {
	sleepTime := models.UserSleep{UserId: user.ID}
	err := config.Db.First(&sleepTime).Error

	if err != nil {
		panic(err)
	}

	dur := time.Since(sleepTime.SleptAt).Minutes()
	countMin := viper.GetInt("main_info.count_minute_sleep_to_get_hp")
	hp := Round(dur / float64(countMin))

	repositories.DeleteUserSleepTime(sleepTime)

	if hp == 0 {
		return "Слишком мало спал..."
	}

	userHp := user.Health + uint(hp)

	if userHp > 100 {
		user.Health = 100
		result = "Хорошо поспал! 💪\nЖизни ♥️ полностью восстановлены! "
	} else {
		user.Health = userHp
		result = fmt.Sprintf("Вы получили %v ♥️ хп!", hp)
	}

	repositories.UpdateUser(user)

	return result
}

func MsgSleepUser() string {
	return "🌙⬛️✨⬛️✨\n⬛️✨💤️✨⬛️\n📚⬛️🛌⬛️🔭\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\nЗа каждые 10 мин сна - вы получите 1 хп!\n\U0001F971Добрых снов!💤"
}

func SleepButton() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("⏰ Проснуться! 🌅", "wakeUp"),
		),
	)
}

// Round возвращает ближайшее целочисленное значение.
func Round(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}

func UserSleep(user models.User, char string) (msg string, buttons tg.InlineKeyboardMarkup) {
	switch char {
	case "wakeUp":
		msg, buttons = userMapController.GetMyMap(user)
		msg = fmt.Sprintf("%s%s%s", msg, viper.GetString("msg_separator"), UserGetHpAfterSleep(user))
		user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	}

	return msg, buttons
}
