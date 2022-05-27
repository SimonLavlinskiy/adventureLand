package newsController

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/src/repositories"
)

func GetNewsMsg() (msg string) {
	news := repositories.GetNews()

	if len(news) == 0 {
		return "Новостей нет ¯ \\ _ (ツ) _ / ¯ "
	}

	for _, n := range news {
		date := n.CreatedAt.Format("02.01.2006")
		msg += fmt.Sprintf("_%s_ - *%s*\n_%s_%s", date, n.Title, n.Text, v.GetString("msg_separator"))
	}

	msg = fmt.Sprintf("📰 *Новости* 📰%s%s", v.GetString("msg_separator"), msg)

	return msg
}
