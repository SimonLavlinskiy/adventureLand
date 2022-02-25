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
var messageSeparator = "\n\n〰️〰️〰️〰️〰️〰️〰️〰️〰️\n"

func messageResolver(update tgbotapi.Update) tgbotapi.MessageConfig {
	user := repository.GetOrCreateUser(update)

	switch user.MenuLocation {
	case "Меню":
		msg = userMenuLocation(update, user)
	case "Карта":
		msg = userMapLocation(update, user)
	case "Профиль":
		msg = userProfileLocation(update, user)
	default:
		msg = userMenuLocation(update, user)
	}

	msg.ParseMode = "markdown"

	return msg
}

func CallbackResolver(update tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	msg.BaseChat.ChatID = update.CallbackQuery.Message.Chat.ID
	buttons := tgbotapi.ReplyKeyboardMarkup{}
	charData := strings.Fields(update.CallbackQuery.Data)
	deletePrevMessage := true

	userTgId := repository.GetUserTgId(update)
	user := repository.GetUser(repository.User{TgId: userTgId})
	ItemLeftHand, ItemRightHand, ItemHead := usersHandItemsView(user)

	if len(charData) != 1 {
		switch charData[0] {
		case "backpackMoving":
			msg = BackPackMoving(charData, update)
		case "goodsMoving":
			msg = GoodsMoving(charData, update)
		case "eatFood":
			UserEatItem(update, charData)
		case "deleteItem":
			UserDeleteItem(update, charData)
		case "dressGood":
			msg = dressUserItem(update, charData)
		case "takeOffGood":
			userTakeOffGood(update, charData)
		case "changeLeftHand":
			userItem, _ := repository.GetUserItem(repository.UserItem{ID: repository.ToInt(charData[1])})
			repository.UpdateUser(update, repository.User{LeftHandId: &userItem.ItemId})
			charDataForOpenGoods := strings.Fields("goodMoving " + charData[2])
			msg = GoodsMoving(charDataForOpenGoods, update)
			msg.Text = fmt.Sprintf("%s%sВы надели %s", msg.Text, messageSeparator, userItem.Item.View)
		case "changeRightHand":
			userItem, _ := repository.GetUserItem(repository.UserItem{ID: repository.ToInt(charData[1])})
			repository.UpdateUser(update, repository.User{RightHandId: &userItem.ItemId})
			charDataForOpenGoods := strings.Fields("goodMoving " + charData[2])
			msg = GoodsMoving(charDataForOpenGoods, update)
			msg.Text = fmt.Sprintf("%s%sВы надели %s", msg.Text, messageSeparator, userItem.Item.View)
		case "changeAvatar":
			res := repository.UpdateUser(update, repository.User{Avatar: charData[1]})
			msg.Text = repository.GetUserInfo(update)
			msg.ReplyMarkup = helpers.ProfileKeyboard(res)
		case "description":
			msg.Text = repository.GetFullDescriptionOfUserItem(repository.UserItem{ID: repository.ToInt(charData[1])})
			deletePrevMessage = false
		case "👋", ItemLeftHand.View, ItemRightHand.View:
			res := directionMovement(update, charData[1])
			resultOfGetItem := repository.UserGetItem(update, res, charData)
			resText, buttons := repository.GetMyMap(update)
			msg.Text = resText + messageSeparator + resultOfGetItem
			msg.ReplyMarkup = buttons
		case ItemHead.View:
			res := directionMovement(update, charData[1])
			status, text := repository.UpdateUserInstrument(update, user, ItemHead)
			if status != "Ok" {
				msg.Text = repository.ViewItemInfo(res) + messageSeparator + text
			} else {
				msg.Text = repository.ViewItemInfo(res)
			}
		case "throwOutItem":
			userWantsToThrowOutItem(update, charData)
		case "countOfDelete":
			msg = userThrowOutItem(update, user, charData)
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
	ItemLeftHand, ItemRightHand, ItemHead := usersHandItemsView(user)
	msg.ChatID = update.Message.Chat.ID

	switch char[0] {
	case "🔼", "🔽", "◀️️", "▶️":
		res := directionMovement(update, char[0])
		repository.UpdateLocation(update, res)
		text := repository.CheckUserHasLighter(update, user)
		msg.Text, buttons = repository.GetMyMap(update)
		msg.Text = msg.Text + messageSeparator + text
		msg.ReplyMarkup = buttons
	case "👋", ItemLeftHand.View, ItemRightHand.View:
		res := directionMovement(update, char[1])
		resultOfGetItem := repository.UserGetItem(update, res, char)
		msg.Text, buttons = repository.GetMyMap(update)
		msg.Text = msg.Text + messageSeparator + resultOfGetItem
		msg.ReplyMarkup = buttons
	case "❗":
		cellLocation := directionMovement(update, char[3])
		cell := repository.GetCellule(repository.Cellule{MapsId: *cellLocation.MapsId, AxisX: *cellLocation.AxisX, AxisY: *cellLocation.AxisY})
		msg.Text = "В зависимости от предмета в твоих руках ты можешь получить разный результат. Выбирай..."
		msg.ReplyMarkup = helpers.ChooseInstrument(char, cell, user)
	case "🚷":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нельзя взять без инструмента в руке")
	case "Рюкзак":
		resUserItems := repository.GetBackpackItems(user.ID)
		msg.Text = MessageBackpackUserItems(resUserItems, 0)
		msg.ReplyMarkup = helpers.BackpackInlineKeyboard(resUserItems, 0)
	case "Вещи":
		userItems := repository.GetInventoryItems(user.ID)
		msg.Text = MessageGoodsUserItems(user, userItems, 0)
		msg.ReplyMarkup = helpers.GoodsInlineKeyboard(user, userItems, 0)
	case "📴":
		userOnline := true
		user = repository.UpdateUser(update, repository.User{OnlineMap: &userOnline})
		msg.Text, buttons = repository.GetMyMap(update)
		msg.Text = msg.Text + messageSeparator + "Онлайн включен!"
		msg.ReplyMarkup = buttons
	case "📳":
		userOnline := false
		user = repository.UpdateUser(update, repository.User{OnlineMap: &userOnline})
		msg.Text, buttons = repository.GetMyMap(update)
		msg.Text = msg.Text + messageSeparator + "Онлайн выключен!"
		msg.ReplyMarkup = buttons
	case ItemHead.View:
		res := directionMovement(update, char[1])
		status, text := repository.UpdateUserInstrument(update, user, ItemHead)
		if status != "Ok" {
			msg.Text = repository.ViewItemInfo(res) + messageSeparator + text
		} else {
			msg.Text = repository.ViewItemInfo(res)
		}
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

	switch newMessage {
	case "🔼", "🔽", "◀️️", "▶️":
		res := directionMovement(update, newMessage)
		repository.UpdateLocation(update, res)
		text := repository.CheckUserHasLighter(update, user)
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text+text)
		msg.ReplyMarkup = buttons
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
		userItemMsg += fmt.Sprintf("%s   %s%s     *HP*:  _%s_ ♥️️     *ST*:  _%s_ \U0001F9C3 ️\n", firstCell, repository.ToString(*item.Count), item.Item.View, repository.ToString(*item.Item.Healing), repository.ToString(*item.Item.Satiety))

	}

	return userItemMsg
}

func MessageGoodsUserItems(user repository.User, userItems []repository.UserItem, rowUser int) string {
	var userItemMsg = "🧥 *Вещички* 🎒\n\n"
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
		userItemMsg += fmt.Sprintf("%s  %s %sшт.   %s %s   (%s/%s)\n", firstCell, item.Item.View, repository.ToString(*item.Count), res, item.Item.Name, repository.ToString(*item.CountUseLeft), repository.ToString(*item.Item.CountUse))

	}

	return userItemMsg
}

func BackPackMoving(charData []string, update tgbotapi.Update) tgbotapi.MessageConfig {
	i := repository.ToInt(charData[1])
	userTgId := repository.GetUserTgId(update)
	user := repository.GetUser(repository.User{TgId: userTgId})
	userItems := repository.GetBackpackItems(user.ID)

	switch i {
	case len(userItems):
		i = i - 1
	}

	msg.Text = MessageBackpackUserItems(userItems, i)
	msg.ReplyMarkup = helpers.BackpackInlineKeyboard(userItems, i)

	return msg
}

func GoodsMoving(charData []string, update tgbotapi.Update) tgbotapi.MessageConfig {
	i := repository.ToInt(charData[1])

	userTgId := repository.GetUserTgId(update)
	user := repository.GetUser(repository.User{TgId: userTgId})
	userItems := repository.GetInventoryItems(user.ID)

	switch i {
	case len(userItems):
		i = i - 1
	}

	msg.Text = MessageGoodsUserItems(user, userItems, i)
	msg.ReplyMarkup = helpers.GoodsInlineKeyboard(user, userItems, i)

	return msg
}

func UserEatItem(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
	userTgId := repository.GetUserTgId(update)
	userItemId := repository.ToInt(charData[1])

	user := repository.GetUser(repository.User{TgId: userTgId})
	userItem, err := repository.GetUserItem(repository.UserItem{ID: userItemId})
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Еда магически исчезла из твоих рук! и ты ее больше не нашел)")
	}

	res := repository.EatItem(update, user, userItem)
	charDataForOpenBackPack := strings.Fields("backpackMoving " + charData[2])
	msg = BackPackMoving(charDataForOpenBackPack, update)
	msg.Text = msg.Text + messageSeparator + res

	return msg
}

func UserDeleteItem(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
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
		msg.Text = "Еда магически исчезла из твоих рук! и ты ее больше не нашел)"
	}

	var charDataForOpenList []string
	switch charData[3] {
	case "good":
		charDataForOpenList = strings.Fields("goodsMoving " + charData[2])
		userTakeOffGood(update, charData)
		msg = GoodsMoving(charDataForOpenList, update)
	case "backpack":
		charDataForOpenList = strings.Fields("backpackMoving " + charData[2])
		msg = BackPackMoving(charDataForOpenList, update)
	}

	msg.Text = fmt.Sprintf("%s%s🗑 Вы выкинули %s%sшт.", msg.Text, messageSeparator, userItem.Item.View, repository.ToString(*userItem.Count))

	return msg
}

func usersHandItemsView(user repository.User) (repository.Item, repository.Item, repository.Item) {
	ItemLeftHand := repository.Item{View: "👋"}
	ItemRightHand := repository.Item{View: "👋"}
	var ItemHead repository.Item

	if user.LeftHand != nil {
		ItemLeftHand = *user.LeftHand
	}
	if user.RightHand != nil {
		ItemRightHand = *user.RightHand
	}
	if user.Head != nil {
		ItemHead = *user.Head
	}

	return ItemLeftHand, ItemRightHand, ItemHead
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
	msg.Text = fmt.Sprintf("%s%sВещь снята!", msg.Text, messageSeparator)
}

func dressUserItem(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
	userItemId := repository.ToInt(charData[1])
	userTgId := repository.GetUserTgId(update)
	user := repository.GetUser(repository.User{TgId: userTgId})
	userItem, _ := repository.GetUserItem(repository.UserItem{ID: userItemId})
	changeHandItem := false

	var result = "Вы надели " + userItem.Item.View

	switch *userItem.Item.DressType {
	case "hand":
		if user.LeftHandId == nil {
			repository.UpdateUser(update, repository.User{LeftHandId: &userItem.ItemId})
		} else if user.RightHandId == nil {
			repository.UpdateUser(update, repository.User{RightHandId: &userItem.ItemId})
		} else {
			result = "У вас заняты все руки! Что хочешь снять?"
			changeHandItem = true
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

	if changeHandItem {
		msg.ReplyMarkup = helpers.ChangeItemInHand(user, userItemId, charData[2])
	} else {
		charDataForOpenGoods := strings.Fields("goodMoving " + charData[2])
		msg = GoodsMoving(charDataForOpenGoods, update)
	}

	msg.Text = fmt.Sprintf("%s%s%s", msg.Text, messageSeparator, result)

	return msg
}

func userThrowOutItem(update tgbotapi.Update, user repository.User, charData []string) tgbotapi.MessageConfig {
	userItem, _ := repository.GetUserItem(repository.UserItem{ID: repository.ToInt(charData[1])})

	*userItem.Count = *userItem.Count - repository.ToInt(charData[3])

	res := repository.UpdateCellUnderUser(update, userItem, repository.ToInt(charData[3]))
	var msgtext string
	if res != "Ok" {
		msgtext = fmt.Sprintf("%s%s", messageSeparator, res)
	} else {
		msgtext = fmt.Sprintf("%sВы сбросили %s %sшт. на карту!", messageSeparator, userItem.Item.View, charData[3])
		repository.UpdateUserItem(user, repository.UserItem{ID: userItem.ID, Count: userItem.Count})
	}

	var charDataForOpenList []string
	switch charData[4] {
	case "good":
		charDataForOpenList = strings.Fields("goodsMoving " + charData[2])
		if *userItem.Count == 0 {
			userTakeOffGood(update, charData)
		}
		msg = GoodsMoving(charDataForOpenList, update)
	case "backpack":
		charDataForOpenList = strings.Fields("backpackMoving " + charData[2])
		msg = BackPackMoving(charDataForOpenList, update)
	}

	msg.Text = msg.Text + msgtext

	return msg
}

func userWantsToThrowOutItem(update tgbotapi.Update, charData []string) tgbotapi.MessageConfig {
	userItem, _ := repository.GetUserItem(repository.UserItem{ID: repository.ToInt(charData[1])})

	if userItem.CountUseLeft != nil && *userItem.CountUseLeft != *userItem.Item.CountUse {
		*userItem.Count = *userItem.Count - 1
	}

	if *userItem.Count == 0 {
		var charDataForOpenList []string
		switch charData[3] {
		case "good":
			charDataForOpenList = strings.Fields("goodsMoving " + charData[2])
			if *userItem.CountUseLeft == *userItem.Item.CountUse {
				userTakeOffGood(update, charData)
			}
			msg = GoodsMoving(charDataForOpenList, update)
		case "backpack":
			charDataForOpenList = strings.Fields("backpackMoving " + charData[2])
			msg = BackPackMoving(charDataForOpenList, update)
		}
		msg.Text = fmt.Sprintf("%s%sНельзя выкинуть на карту предмет, который уже был использован!", msg.Text, messageSeparator)
	} else {
		msg.ReplyMarkup = helpers.CountItemUserWantsToThrow(charData, userItem)
		msg.Text = fmt.Sprintf("%sСколько %s ты хочешь скинуть?", messageSeparator, userItem.Item.View)
	}

	return msg
}
