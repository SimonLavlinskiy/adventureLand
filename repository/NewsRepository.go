package repository

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	"time"
)

type News struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"embedded"`
	Text      string    `gorm:"embedded"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func GetNews() []News {
	var results []News

	config.Db.Find(&results).Order("id desc")

	return results
}

func GetNewsMsg() (msg string) {
	news := GetNews()
	if len(news) == 0 {
		return "Новостей нет ¯ \\ _ (ツ) _ / ¯ "
	}
	for _, n := range news {
		date := n.CreatedAt.Format("02.01.2006")
		msg += fmt.Sprintf("_%s_ - *%s*\n_%s_%s", date, n.Title, n.Text, v.GetString("msg_separator"))
	}

	return fmt.Sprintf("📰 *Новости* 📰%s%s", v.GetString("msg_separator"), msg)
}
