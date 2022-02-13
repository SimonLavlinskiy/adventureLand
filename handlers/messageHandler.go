package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/helpers"
	"project0/repository"
	"strings"
	"time"
)

var msg tgbotapi.MessageConfig

func messageResolver(update tgbotapi.Update) tgbotapi.MessageConfig {
	resUser := repository.GetOrCreateUser(update)

	switch resUser.MenuLocation {
	case "Меню":
		msg = userMenuLocation(update, resUser)
	case "Карта":
		msg = userMapLocation(update, resUser)
	case "Профиль":
		msg = userProfileLocation(update, resUser)
	default:
		msg = userMenuLocation(update, resUser)
	}

	msg.ParseMode = "markdown"

	return msg
}

func useSpecialCell(update tgbotapi.Update, char []string, user repository.User) tgbotapi.MessageConfig {
	buttons := tgbotapi.ReplyKeyboardMarkup{}

	viewItemLeftHand, viewItemRightHand := usersHandsItemsView(user)

	switch char[0] {
	case "🔼", "🔽", "◀️️", "▶️":
		res := directionMovement(update, char[0])
		repository.UpdateLocation(update, res)
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
	case "👋", viewItemLeftHand, viewItemRightHand:
		res := directionMovement(update, char[1])
		resultOfGetItem := repository.UserGetItem(update, res, char)
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text+"\n\n"+resultOfGetItem)
		msg.ReplyMarkup = buttons
	case "🚷":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нельзя взять без инструмента в руке")
	default:
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
	}

	return msg
}

func userMenuLocation(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	buttons := tgbotapi.ReplyKeyboardMarkup{}
	newMessage := update.Message.Text

	switch newMessage {
	case "🗺 Карта 🗺":
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
		repository.UpdateUser(update, repository.User{MenuLocation: "Карта"})
	case user.Avatar + " Профиль 👔":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
		msg.ReplyMarkup = helpers.ProfileKeyboard(user)
		repository.UpdateUser(update, repository.User{MenuLocation: "Профиль"})
	default:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
		msg.ReplyMarkup = helpers.MainKeyboard(user)
		repository.UpdateUser(update, repository.User{MenuLocation: "Меню"})
	}

	return msg
}

func userMapLocation(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	newMessage := update.Message.Text
	char := strings.Fields(newMessage)

	if len(char) != 1 {
		msg = useSpecialCell(update, char, user)
	} else {
		msg = useDefaultCell(update, user)
	}

	return msg
}

func userProfileLocation(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	newMessage := update.Message.Text

	if user.Username == "waiting" {
		repository.UpdateUser(update, repository.User{Username: newMessage})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
		msg.ReplyMarkup = helpers.ProfileKeyboard(user)
	} else {
		switch newMessage {
		case "📝 Изменить имя? 📝":
			repository.UpdateUser(update, repository.User{Username: "waiting"})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case user.Avatar + " Изменить аватар? " + user.Avatar:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар...")
			msg.ReplyMarkup = helpers.EmodjiInlineKeyboard()
		case "/menu", "Меню":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
			msg.ReplyMarkup = helpers.MainKeyboard(user)
			repository.UpdateUser(update, repository.User{MenuLocation: "Меню"})
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
			msg.ReplyMarkup = helpers.ProfileKeyboard(user)
		}
	}
	return msg
}

func directionMovement(update tgbotapi.Update, direction string) repository.Location {
	res := repository.GetOrCreateMyLocation(update)

	switch direction {
	case "🔼":
		y := *res.AxisY + 1
		return repository.Location{Map: res.Map, AxisX: res.AxisX, AxisY: &y}
	case "🔽":
		y := *res.AxisY - 1
		return repository.Location{Map: res.Map, AxisX: res.AxisX, AxisY: &y}
	case "◀️️":
		x := *res.AxisX - 1
		return repository.Location{Map: res.Map, AxisX: &x, AxisY: res.AxisY}
	case "▶️":
		x := *res.AxisX + 1
		return repository.Location{Map: res.Map, AxisX: &x, AxisY: res.AxisY}
	}
	return res
}

func useDefaultCell(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	newMessage := update.Message.Text
	buttons := tgbotapi.ReplyKeyboardMarkup{}
	currentTime := time.Now()

	switch newMessage {
	case "🔼", "🔽", "◀️️", "▶️":
		res := directionMovement(update, newMessage)
		repository.UpdateLocation(update, res)
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
	case "🎒":
		resUser := repository.GetOrCreateUser(update)
		resUserItems := repository.GetUserItems(resUser.ID, "food")
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, MessageBackpackUserItems(resUserItems, 0))
		msg.ReplyMarkup = helpers.BackpackInlineKeyboard(resUserItems, 0)
	case "🧥🎒":

	case "\U0001F7E6": // Вода
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ты не похож на Jesus! 👮‍♂️")
	case "🕦":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, currentTime.Format("15:04:05")+"\nЧасики тикают...")
	case user.Avatar:
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update)+"\n \n"+msg.Text)
		msg.ReplyMarkup = buttons
	case "/menu", "Меню":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
		msg.ReplyMarkup = helpers.MainKeyboard(user)
		repository.UpdateUser(update, repository.User{MenuLocation: "Меню"})
	case "🎰":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "💰💵🤑 Ставки на JOY CASINO дот COM! 🤑💵💰 ")
	default:
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
	}

	return msg
}

func CallbackResolver(update tgbotapi.Update) tgbotapi.MessageConfig {
	charData := strings.Fields(update.CallbackQuery.Data)

	if len(charData) != 1 {
		switch charData[0] {
		case "backpackMoving":
			msg = BackPackMoving(charData, update)
		case "eatFood":
			UserEatItem(update, charData)
		case "throwOutFood":
			UserThrowOutFood(update, charData)
		case "changeAvatar":
			res := repository.UpdateUser(update, repository.User{Avatar: charData[1]})
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, repository.GetUserInfo(update))
			msg.ReplyMarkup = helpers.ProfileKeyboard(res)
		}
	} else {
		fmt.Println("callbackQuery содержит 1 элемент")
	}

	msg.ParseMode = "markdown"
	return msg
}

func MessageBackpackUserItems(userItems []repository.UserItem, rowUser int) string {
	var userItemMsg = "🎒 *Рюкзачок*\n \n"

	if len(userItems) == 0 {
		return "👻 \U0001F9B4  Пусто .... 🕸 🕷"
	}

	for i, item := range userItems {
		var firstCell string
		switch rowUser {
		case i:
			firstCell += item.User.Avatar
		case i + 1, i - 1:
			firstCell += "◻️"
		case i + 2, i - 2:
			firstCell += "◽️️"
		default:
			firstCell += "▫️"
		}
		userItemMsg += firstCell + "   " + helpers.ToString(*item.Count) + item.Item.View +
			"     *HP*:  _+" + helpers.ToString(*item.Item.Healing) + "_ ♥️️" +
			"     *ST*:  _+" + helpers.ToString(*item.Item.Satiety) + "_\U0001F9C3 ️\n"

	}

	return userItemMsg
}

func BackPackMoving(charData []string, update tgbotapi.Update) tgbotapi.MessageConfig {
	i := helpers.ToInt(charData[1])

	user := repository.GetUser(repository.User{TgId: uint(update.CallbackQuery.From.ID)})
	userItems := repository.GetUserItems(user.ID, "food")

	switch i {
	case len(userItems):
		i = i - 1
	}

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, MessageBackpackUserItems(userItems, i))
	msg.ReplyMarkup = helpers.BackpackInlineKeyboard(userItems, i)

	return msg
}

func UserEatItem(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
	userItemId := helpers.ToInt(charData[1])
	userTgId := uint(update.CallbackQuery.From.ID)

	user := repository.GetUser(repository.User{TgId: userTgId})
	userItem, err := repository.GetUserItem(repository.UserItem{ID: userItemId})
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Еда магически исчезла из твоих рук! и ты ее больше не нашел)")
	}

	res := repository.EatItem(update, user, userItem)
	charDataForOpenBackPack := strings.Fields("backpackMoving " + charData[2])
	msg = BackPackMoving(charDataForOpenBackPack, update)
	msg.Text = msg.Text + "\n\n" + res

	return msg
}

func UserThrowOutFood(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
	userItemId := helpers.ToInt(charData[1])
	userTgId := uint(update.CallbackQuery.From.ID)

	user := repository.GetUser(repository.User{TgId: userTgId})
	userItem, err := repository.GetUserItem(repository.UserItem{ID: userItemId})

	countAfterUserThrowOutItem := 0
	var updateUserItemStruct = repository.UserItem{
		ID:    userItemId,
		Count: &countAfterUserThrowOutItem,
	}

	repository.UpdateUserItem(user, updateUserItemStruct)
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Еда магически исчезла из твоих рук! и ты ее больше не нашел)")
	}

	charDataForOpenBackPack := strings.Fields("backpackMoving " + charData[2])
	msg = BackPackMoving(charDataForOpenBackPack, update)
	msg.Text = msg.Text + "\n\n" + "🗑 Вы выкинули все " + helpers.ToString(*userItem.Count) + " " + userItem.Item.View

	return msg
}

func usersHandsItemsView(user repository.User) (string, string) {
	viewItemLeftHand := "👋"
	viewItemRightHand := "👋"
	if user.LeftHand != nil {
		viewItemLeftHand = user.LeftHand.View
	}
	if user.RightHand != nil {
		viewItemRightHand = user.RightHand.View
	}

	return viewItemLeftHand, viewItemRightHand
}
