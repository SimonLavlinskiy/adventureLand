package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/repository"
	"strings"
	"time"
)

var msg tgbotapi.MessageConfig

//var updateMsg tgbotapi.EditMessageTextConfig

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
		msg.ReplyMarkup = profileKeyboard(user)
		repository.UpdateUser(update, repository.User{MenuLocation: "Профиль"})
	case "👜 Инвентарь 👜":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, newMessage)
		msg.ReplyMarkup = backpackKeyboard
		repository.UpdateUser(update, repository.User{MenuLocation: "Инвентарь"})
	default:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
		msg.ReplyMarkup = mainKeyboard(user)
		repository.UpdateUser(update, repository.User{MenuLocation: "Меню"})
	}

	return msg
}

func userMapLocation(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	newMessage := update.Message.Text
	char := strings.Fields(newMessage)

	if len(char) != 1 {
		msg = useItems(update, char)
	} else {
		msg = useDefaultItems(update, user)
	}

	return msg
}

func userProfileLocation(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	newMessage := update.Message.Text

	if user.Username == "waiting" {
		repository.UpdateUser(update, repository.User{Username: newMessage})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
		msg.ReplyMarkup = profileKeyboard(user)
	} else {
		switch newMessage {
		case "📝 Изменить имя? 📝":
			repository.UpdateUser(update, repository.User{Username: "waiting"})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case user.Avatar + " Изменить аватар? " + user.Avatar:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар...")
			msg.ReplyMarkup = EmodjiInlineKeyboard()
		case "/menu", "Меню":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
			msg.ReplyMarkup = mainKeyboard(user)
			repository.UpdateUser(update, repository.User{MenuLocation: "Меню"})
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
			msg.ReplyMarkup = profileKeyboard(user)
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

func useDefaultItems(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
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
		resUserItems := repository.GetUserItems(resUser.ID)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, MessageBackpackUserItems(resUserItems, 0))
		msg.ReplyMarkup = backpackInlineKeyboard(resUserItems, 0)
	case "\U0001F7E6": // Вода
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ты не похож на Jesus! 👮‍♂️")
	case "🕦":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, currentTime.Format("15:04:05")+"\nЧасики тикают...")
	case user.Avatar:
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update)+"\n \n"+msg.Text)
	case "/menu", "Меню":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
		msg.ReplyMarkup = mainKeyboard(user)
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

func useItems(update tgbotapi.Update, char []string) tgbotapi.MessageConfig {
	buttons := tgbotapi.ReplyKeyboardMarkup{}

	switch char[0] {
	case "🔼", "🔽", "◀️️", "▶️":
		res := directionMovement(update, char[0])
		repository.UpdateLocation(update, res)
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
	case "👋":
		res := directionMovement(update, char[1])
		countItem := repository.UserGetItem(update, res)
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text+"\n\nТы взял: "+repository.ToString(countItem)+" шт "+char[2])
		msg.ReplyMarkup = buttons
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
		case "changeAvatar":
			res := repository.UpdateUser(update, repository.User{Avatar: charData[1]})
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, repository.GetUserInfo(update))
			msg.ReplyMarkup = profileKeyboard(res)
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
		userItemMsg += firstCell + "   " + repository.ToString(*item.Count) + item.Item.View +
			"     *HP*:  _+" + repository.ToString(*item.Item.Healing) + "_ ♥️️" +
			"     *ST*:  _+" + repository.ToString(*item.Item.Satiety) + "_\U0001F9C3 ️\n"

	}

	return userItemMsg
}

func backpackInlineKeyboard(items []repository.UserItem, i int) tgbotapi.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(items[i].Item.View+" "+
				repository.ToString(*items[i].Count)+"шт."+
				"   +"+repository.ToString(*items[i].Item.Healing)+" ♥️️"+
				"   +"+repository.ToString(*items[i].Item.Satiety)+"\U0001F9C3", "callbackAnswerAlert"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍽 1шт", "eatFood "+repository.ToString(items[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", "backpackMoving "+repository.ToString(i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑 все!", "throwOutFood "+repository.ToString(items[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔻", "backpackMoving "+repository.ToString(i+1)),
		),
	)
}

func BackPackMoving(charData []string, update tgbotapi.Update) tgbotapi.MessageConfig {
	i := repository.ToInt(charData[1])

	user := repository.GetUser(repository.User{TgId: uint(update.CallbackQuery.From.ID)})
	userItems := repository.GetUserItems(user.ID)

	switch i {
	case len(userItems):
		i = i - 1
	}

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, MessageBackpackUserItems(userItems, i))
	msg.ReplyMarkup = backpackInlineKeyboard(userItems, i)

	return msg
}

func UserEatItem(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
	userItemId := repository.ToInt(charData[1])
	userTgId := uint(update.CallbackQuery.From.ID)

	user := repository.GetUser(repository.User{TgId: userTgId})
	item, err := repository.GetUserItem(repository.UserItem{ID: userItemId})
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Еда магически исчезла из твоих рук! и ты ее больше не нашел)")
	}

	res := repository.EatItem(update, user, item)
	charDataForOpenBackPack := strings.Fields("backpackMoving " + charData[2])
	msg = BackPackMoving(charDataForOpenBackPack, update)
	msg.Text = res + "\n\n" + msg.Text

	return msg
}

func EmodjiInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	var listOfAvatar []string
	listOfAvatar = strings.Fields("🐶 🐱 🐭 🐹 🐰 🦊 🐻 🐼 ‍️🐨 🐯 🦁 🐮 🐷 🐸 🐵 🐦 🐧 🐔 🐤 🐥 🦆 🐴 🦄 🐺 🐗 🐝 🦋 🐛 🐌 🐞 🪲 🪰 🐜 🕷 🪳 🦖 🦕 🐙 🦀 🐟 🐠 🐡 🦭")

	for x := 0; x < len(listOfAvatar); x = x + 8 {
		var row []tgbotapi.InlineKeyboardButton
		for i := 0; i < 8; i++ {
			sum := x + i
			if len(listOfAvatar) > sum {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(listOfAvatar[sum], "changeAvatar "+listOfAvatar[sum]))
			}
		}
		buttons = append(buttons, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func profileKeyboard(user repository.User) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📝 Изменить имя? 📝"),
			tgbotapi.NewKeyboardButton(user.Avatar+" Изменить аватар? "+user.Avatar),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Меню"),
		),
	)
}

func mainKeyboard(user repository.User) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗺 Карта 🗺"),
			tgbotapi.NewKeyboardButton(user.Avatar+" Профиль 👔"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👜 Инвентарь 👜"),
		),
	)
}
