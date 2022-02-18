package handlers

import (
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

func CallbackResolver(update tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	buttons := tgbotapi.ReplyKeyboardMarkup{}
	charData := strings.Fields(update.CallbackQuery.Data)
	deletePrevMessage := true

	user := repository.GetUser(repository.User{TgId: uint(update.CallbackQuery.From.ID)})
	viewItemLeftHand, viewItemRightHand := usersHandItemsView(user)

	if len(charData) != 1 {
		switch charData[0] {
		case "backpackMoving":
			msg = BackPackMoving(charData, update)
		case "goodsMoving":
			msg = GoodsMoving(charData, update)
		case "eatFood":
			UserEatItem(update, charData)
		case "throwOutFood", "throwOutGood":
			UserThrowOutItem(update, charData)
		case "dressGood":
			dressUserItem(update, charData)
		case "takeOffGood":
			userTakeOffGood(update, charData)
		case "changeAvatar":
			res := repository.UpdateUser(update, repository.User{Avatar: charData[1]})
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, repository.GetUserInfo(update))
			msg.ReplyMarkup = helpers.ProfileKeyboard(res)
		case "description":
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, repository.GetFullDescriptionOfUserItem(repository.UserItem{ID: repository.ToInt(charData[1])}))
			deletePrevMessage = false
		case "👋", viewItemLeftHand, viewItemRightHand:
			res := directionMovement(update, charData[1])
			resultOfGetItem := repository.UserGetItem(update, res, charData)
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msg.Text+"\n\n"+resultOfGetItem)
			msg.ReplyMarkup = buttons
		}
	} else {
		switch charData[0] {
		case "cancel":
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = buttons
		}
	}

	msg.ParseMode = "markdown"
	return msg, deletePrevMessage
}

func useSpecialCell(update tgbotapi.Update, char []string, user repository.User) tgbotapi.MessageConfig {
	buttons := tgbotapi.ReplyKeyboardMarkup{}

	user = repository.GetUser(repository.User{TgId: user.TgId})
	viewItemLeftHand, viewItemRightHand := usersHandItemsView(user)

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
	case "❗":
		cellLocation := directionMovement(update, char[3])
		cell := repository.GetCellule(repository.Cellule{MapsId: *cellLocation.MapsId, AxisX: *cellLocation.AxisX, AxisY: *cellLocation.AxisY})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "В зависимости от предмета в твоих руках ты можешь получить разный результат. Выбирай...")
		msg.ReplyMarkup = helpers.ChooseInstrument(char, cell, user)
	case "🚷":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нельзя взять без инструмента в руке")
	default:
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text+"\n\nНет инструмента в руке!")
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
		case "Офлайн (📴♻️📳)":
			userOnline := true
			user = repository.UpdateUser(update, repository.User{OnlineMap: &userOnline})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
			msg.ReplyMarkup = helpers.ProfileKeyboard(user)
		case "Онлайн (📳♻️📴)":
			userOnline := false
			user = repository.UpdateUser(update, repository.User{OnlineMap: &userOnline})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
			msg.ReplyMarkup = helpers.ProfileKeyboard(user)
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
		return repository.Location{MapsId: res.MapsId, AxisX: res.AxisX, AxisY: &y}
	case "🔽":
		y := *res.AxisY - 1
		return repository.Location{MapsId: res.MapsId, AxisX: res.AxisX, AxisY: &y}
	case "◀️️":
		x := *res.AxisX - 1
		return repository.Location{MapsId: res.MapsId, AxisX: &x, AxisY: res.AxisY}
	case "▶️":
		x := *res.AxisX + 1
		return repository.Location{MapsId: res.MapsId, AxisX: &x, AxisY: res.AxisY}
	}
	return res
}

func useDefaultCell(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	newMessage := update.Message.Text
	buttons := tgbotapi.ReplyKeyboardMarkup{}
	currentTime := time.Now()
	//userTgId := repository.GetUserTgId(update)

	switch newMessage {
	case "🔼", "🔽", "◀️️", "▶️":
		res := directionMovement(update, newMessage)
		repository.UpdateLocation(update, res)
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
	case "🎒":
		resUserItems := repository.GetBackpackItems(user.ID)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, MessageBackpackUserItems(resUserItems, 0))
		msg.ReplyMarkup = helpers.BackpackInlineKeyboard(resUserItems, 0)
	case "🧥🎒":
		userItems := repository.GetInventoryItems(user.ID)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, MessageGoodsUserItems(user, userItems, 0))
		msg.ReplyMarkup = helpers.GoodsInlineKeyboard(user, userItems, 0)
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
	case "🚫":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Сюда нельзя! 🚫")
	default:
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
	}

	return msg
}

func MessageBackpackUserItems(userItems []repository.UserItem, rowUser int) string {
	var userItemMsg = "🎒 *Рюкзачок*\n \n"

	if len(userItems) == 0 {
		return userItemMsg + "👻 \U0001F9B4  Пусто .... 🕸 🕷"
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

func MessageGoodsUserItems(user repository.User, userItems []repository.UserItem, rowUser int) string {
	var userItemMsg = "🧥 *Вещички* 🎒  (✅ - Надето)\n\n"
	userItemMsg = messageUserDressedGoods(user) + userItemMsg

	if len(userItems) == 0 {
		return userItemMsg + "👻 \U0001F9B4  Пусто .... 🕸 🕷"
	}

	for i, item := range userItems {
		_, res := repository.IsDressedItem(user, userItems[i])

		if res == "takeOffGood" {
			res = "✅"
		} else {
			res = ""
		}

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
		userItemMsg += firstCell + "   " + item.Item.View + " " + repository.ToString(*item.Count) +
			"шт.     " + res + " " + item.Item.Name + "    " + "\n"

	}

	return userItemMsg
}

func BackPackMoving(charData []string, update tgbotapi.Update) tgbotapi.MessageConfig {
	i := repository.ToInt(charData[1])

	user := repository.GetUser(repository.User{TgId: uint(update.CallbackQuery.From.ID)})
	userItems := repository.GetBackpackItems(user.ID)

	switch i {
	case len(userItems):
		i = i - 1
	}

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, MessageBackpackUserItems(userItems, i))
	msg.ReplyMarkup = helpers.BackpackInlineKeyboard(userItems, i)

	return msg
}

func GoodsMoving(charData []string, update tgbotapi.Update) tgbotapi.MessageConfig {
	i := repository.ToInt(charData[1])

	user := repository.GetUser(repository.User{TgId: uint(update.CallbackQuery.From.ID)})
	userItems := repository.GetInventoryItems(user.ID)

	switch i {
	case len(userItems):
		i = i - 1
	}

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, MessageGoodsUserItems(user, userItems, i))
	msg.ReplyMarkup = helpers.GoodsInlineKeyboard(user, userItems, i)

	return msg
}

func UserEatItem(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
	userItemId := repository.ToInt(charData[1])
	userTgId := repository.GetUserTgId(update)
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

func UserThrowOutItem(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
	userItemId := repository.ToInt(charData[1])
	userTgId := repository.GetUserTgId(update)
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

	var charDataForOpenList []string
	switch charData[0] {
	case "throwOutGood":
		charDataForOpenList = strings.Fields("goodsMoving " + charData[2])
		userTakeOffGood(update, charData)
		msg = GoodsMoving(charDataForOpenList, update)
	case "throwOutFood":
		charDataForOpenList = strings.Fields("backpackMoving " + charData[2])
		msg = BackPackMoving(charDataForOpenList, update)
	}

	msg.Text = msg.Text + "\n\n" + "🗑 Вы выкинули " + repository.ToString(*userItem.Count) + "шт. " + userItem.Item.View

	return msg
}

func usersHandItemsView(user repository.User) (string, string) {
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

func messageUserDressedGoods(user repository.User) string {
	var head string
	var body string
	var leftHand string
	var rightHand string
	var foot string
	var shoes string

	if user.Head != nil {
		head = user.Head.View
	} else {
		head = "🟦"
	}
	if user.LeftHand != nil {
		leftHand = user.LeftHand.View
	} else {
		leftHand = "✋"
	}
	if user.RightHand != nil {
		rightHand = user.RightHand.View
	} else {
		rightHand = "🤚"
	}
	if user.Body != nil {
		body = user.Body.View
	} else {
		body = "👔"
	}
	if user.Foot != nil {
		foot = user.Foot.View
	} else {
		foot = "\U0001FA72"
	}
	if user.Shoes != nil {
		shoes = user.Shoes.View
	} else {
		shoes = "👣"
	}

	var messageUserGoods = "\U0001F7E6☁️\U0001F7E6\U0001F7E6☀️\n" +
		"\U0001F7E6\U0001F7E6" + head + "\U0001F7E6\U0001F7E6\n" +
		"\U0001F7E6\U0001F7E6" + user.Avatar + "\U0001F7E6\U0001F7E6\n" +
		"\U0001F7E6" + leftHand + body + rightHand + "\U0001F7E6\n" +
		"\U0001F7E6\U0001F7E6" + foot + "\U0001F7E6\U0001F7E6\n" +
		"\U0001F7E9\U0001F7E9" + shoes + "\U0001F7E9\U0001F7E9\n" +
		"\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\n\n"

	return messageUserGoods
}

func userTakeOffGood(update tgbotapi.Update, charData []string) {
	userItemId := repository.ToInt(charData[1])
	userTgId := repository.GetUserTgId(update)
	user := repository.GetUser(repository.User{TgId: userTgId})
	userItem, _ := repository.GetUserItem(repository.UserItem{ID: userItemId})

	if user.HeadId != nil && userItem.ItemId == *user.HeadId {
		repository.SetNullUserField(update, "head_id")
	} else if user.LeftHandId != nil && userItem.ItemId == *user.LeftHandId {
		repository.SetNullUserField(update, "left_hand_id")
	} else if user.RightHandId != nil && userItem.ItemId == *user.RightHandId {
		repository.SetNullUserField(update, "right_hand_id")
	} else if user.BodyId != nil && userItem.ItemId == *user.BodyId {
		repository.SetNullUserField(update, "body_id")
	} else if user.FootId != nil && userItem.ItemId == *user.FootId {
		repository.SetNullUserField(update, "foot_id")
	} else if user.ShoesId != nil && userItem.ItemId == *user.ShoesId {
		repository.SetNullUserField(update, "shoes_id")
	}

	charDataForOpenGoods := strings.Fields("goodMoving " + charData[2])
	msg = GoodsMoving(charDataForOpenGoods, update)
	msg.Text = msg.Text + "\n\n" + "Вещь снята"
}

func dressUserItem(update tgbotapi.Update, charData []string) {
	userItemId := repository.ToInt(charData[1])
	userTgId := repository.GetUserTgId(update)
	user := repository.GetUser(repository.User{TgId: userTgId})
	userItem, _ := repository.GetUserItem(repository.UserItem{ID: userItemId})

	var result = "Вы надели " + userItem.Item.View

	switch *userItem.Item.DressType {
	case "hand":
		if user.LeftHandId == nil {
			repository.UpdateUser(update, repository.User{LeftHandId: &userItem.ItemId})
		} else if user.RightHandId == nil {
			repository.UpdateUser(update, repository.User{RightHandId: &userItem.ItemId})
		} else {
			result = "У вас заняты все руки! Сначала снимите предмет, чтоб надеть другой"
		}
	case "head":
		if user.HeadId == nil {
			repository.UpdateUser(update, repository.User{HeadId: &userItem.ItemId})
		} else {
			result = "Сначала снимите предмет, чтоб надеть другой"
		}
	case "body":
		if user.BodyId == nil {
			repository.UpdateUser(update, repository.User{BodyId: &userItem.ItemId})
		} else {
			result = "Сначала снимите предмет, чтоб надеть другой"
		}
	case "foot":
		if user.FootId == nil {
			repository.UpdateUser(update, repository.User{FootId: &userItem.ItemId})
		} else {
			result = "Сначала снимите предмет, чтоб надеть другой"
		}
	case "shoes":
		if user.ShoesId == nil {
			repository.UpdateUser(update, repository.User{ShoesId: &userItem.ItemId})
		} else {
			result = "Сначала снимите предмет, чтоб надеть другой"
		}
	}

	charDataForOpenGoods := strings.Fields("goodMoving " + charData[2])
	msg = GoodsMoving(charDataForOpenGoods, update)
	msg.Text = msg.Text + "\n\n" + result
}
