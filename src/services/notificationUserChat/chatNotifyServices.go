package notificationUserChat

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
	"project0/src/repositories"
	"strings"
)

func SendUserMessageAllChatUsers(update tg.Update) (chUsWithoutSender []models.ChatUser, message string) {
	user := repositories.GetOrCreateUser(update)
	chUser := repositories.GetChatOfUser(user)
	chatUsers := models.Chat{ID: chUser.ChatID}.GetChatUsers()

	for _, chatUser := range chatUsers {
		if chatUser.User.TgId != uint(update.Message.From.ID) {
			chUsWithoutSender = append(chUsWithoutSender, chatUser)
		}
	}

	replacer := strings.NewReplacer(
		"/start", fmt.Sprintf("<i>%s</i> %s <code>присоединился к чатику<code>", user.Avatar, user.Username),
	)
	userMsg := replacer.Replace(update.Message.Text)

	message = fmt.Sprintf("<code>От %s %s %s</code>%s%s", user.Avatar, user.Username, user.Avatar, v.GetString("msg_separator"), userMsg)

	return chUsWithoutSender, message
}
