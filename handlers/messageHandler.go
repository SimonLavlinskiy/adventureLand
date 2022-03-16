package handlers

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/helpers"
	r "project0/repository"
	"strings"
	"time"
)

var msg tg.MessageConfig

func messageResolver(update tg.Update) tg.MessageConfig {
	user := r.GetOrCreateUser(update)

	switch user.MenuLocation {
	case v.GetString("user_location.menu"):
		msg = userMenuLocation(update, user)
	case v.GetString("user_location.maps"):
		msg = userMapLocation(update, user)
	case v.GetString("user_location.profile"):
		msg = userProfileLocation(update, user)
	default:
		msg = userMenuLocation(update, user)
	}

	msg.ParseMode = "markdown"

	return msg
}

func callBackResolver(update tg.Update) (tg.MessageConfig, bool) {
	msg.ChatID = update.CallbackQuery.Message.Chat.ID
	buttons := tg.ReplyKeyboardMarkup{}
	charData := strings.Fields(update.CallbackQuery.Data)
	deletePrevMessage := true

	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	ItemLeftHand, ItemRightHand, ItemHead := usersHandItemsView(user)

	if len(charData) == 1 {
		switch charData[0] {
		case v.GetString("callback_char.cancel"):
			msg.Text, buttons = r.GetMyMap(update)
			msg = tg.NewMessage(update.CallbackQuery.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = buttons
		}
	}

	switch charData[0] {
	case v.GetString("callback_char.backpack_moving"):
		msg = BackPackMoving(charData, update)
	case v.GetString("callback_char.goods_moving"):
		msg = GoodsMoving(charData, update)
	case v.GetString("callback_char.eat_food"):
		UserEatItem(update, charData)
	case v.GetString("callback_char.delete_item"):
		UserDeleteItem(update, charData)
	case v.GetString("callback_char.dress_good"):
		msg = dressUserItem(update, charData)
	case v.GetString("callback_char.take_off_good"):
		userTakeOffGood(update, charData)
	case v.GetString("callback_char.change_left_hand"), v.GetString("callback_char.change_right_hand"):
		userItem := r.UserItem{ID: r.ToInt(charData[1])}.UserGetUserItem()
		updateUserHand(update, charData, userItem)
		charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		msg = GoodsMoving(charDataForOpenGoods, update)
		msg.Text = fmt.Sprintf("%s%sВы надели %s", msg.Text, v.GetString("msg_separator"), userItem.Item.View)
	case v.GetString("callback_char.change_avatar"):
		res := r.User{Avatar: charData[1]}.UpdateUser(update)
		msg.Text = r.GetUserInfo(update)
		msg.ReplyMarkup = helpers.ProfileKeyboard(res)
	case v.GetString("callback_char.description"):
		msg.Text = r.UserItem{ID: r.ToInt(charData[1])}.GetFullDescriptionOfUserItem()
		deletePrevMessage = false
	case v.GetString("callback_char.workbench"):
		msg = Workbench(nil, charData)
	case v.GetString("callback_char.receipt"):
		msg.Text = fmt.Sprintf("📖 *Рецепты*: 📖%s%s", v.GetString("msg_separator"), AllReceiptsMsg())
		msg.ReplyMarkup = nil
		deletePrevMessage = false
	case v.GetString("callback_char.put_item"):
		userItem := r.GetUserItems(user.ID)
		msg.ReplyMarkup = helpers.ChooseUserItemButton(userItem, charData)
		msg = OpenWorkbenchMessage(charData)
		msg.Text = fmt.Sprintf("%s%sВыбери предмет:", msg.Text, v.GetString("msg_separator"))
	case v.GetString("callback_char.put_count_item"):
		msg = OpenWorkbenchMessage(charData)
		msg = PutCountComponent(charData)
		msg.Text = fmt.Sprintf("%s%s⚠️ Сколько выкладываешь?", msg.Text, v.GetString("msg_separator"))
	case v.GetString("callback_char.make_new_item"):
		resp := GetReceiptFromData(charData)
		receipt := r.FindReceiptForUser(resp)
		msg, deletePrevMessage = UserCraftItem(user, receipt)
	case v.GetString("message.emoji.hand"), ItemLeftHand.View, ItemRightHand.View:
		msg = UserUseHandOrInstrument(update, charData)
	case v.GetString("message.emoji.foot"):
		msg = UserMoving(update, user, charData[1])
	case ItemHead.View:
		res := directionCell(update, charData[1])
		status, text := r.UpdateUserInstrument(update, user, ItemHead)
		if status != "Ok" {
			msg.Text = fmt.Sprintf("%s%s%s", r.ViewItemInfo(res), v.GetString("msg_separator"), text)
		} else {
			msg.Text = r.ViewItemInfo(res)
		}
	case v.GetString("callback_char.throw_out_item"):
		userWantsToThrowOutItem(update, charData)
	case v.GetString("callback_char.count_of_delete"):
		msg = userThrowOutItem(update, user, charData)
	case "quests":
		msg.Text = v.GetString("user_location.tasks_menu_message")
		msg.ReplyMarkup = helpers.AllQuestsMessageKeyboard(user)
	case "quest":
		msg = OpenQuest(uint(r.ToInt(charData[1])), user)
	case "user_get_quest":
		r.UserQuest{
			UserId:  user.ID,
			QuestId: uint(r.ToInt(charData[1])),
		}.GetOrCreateUserQuest()
		msg = OpenQuest(uint(r.ToInt(charData[1])), user)
	case "user_done_quest":
		msg = UserDoneQuest(uint(r.ToInt(charData[1])), user)
	}

	msg.ParseMode = "markdown"
	return msg, deletePrevMessage
}

func useSpecialCell(update tg.Update, char []string, user r.User) tg.MessageConfig {
	buttons := tg.ReplyKeyboardMarkup{}
	ItemLeftHand, ItemRightHand, ItemHead := usersHandItemsView(user)
	msg.ChatID = update.Message.Chat.ID

	switch char[0] {
	case v.GetString("message.doing.up"), v.GetString("message.doing.down"), v.GetString("message.doing.left"), v.GetString("message.doing.right"):
		msg = UserMoving(update, user, char[0])
	case v.GetString("message.emoji.foot"):
		msg = UserMoving(update, user, char[1])
	case v.GetString("message.emoji.hand"), ItemLeftHand.View, ItemRightHand.View:
		msg = UserUseHandOrInstrument(update, char)
	case v.GetString("message.emoji.exclamation_mark"):
		cellLocation := directionCell(update, char[3])
		cell := r.Cell{MapsId: *cellLocation.MapsId, AxisX: *cellLocation.AxisX, AxisY: *cellLocation.AxisY}
		cell = cell.GetCell()
		msg.Text = "В зависимости от предмета в твоих руках ты можешь получить разный результат. Выбирай..."
		msg.ReplyMarkup = helpers.ChooseInstrument(char, cell, user)
	case v.GetString("message.emoji.stop_use"):
		msg = tg.NewMessage(update.Message.Chat.ID, v.GetString("errors.user_not_has_item_in_hand"))
	case "Рюкзак":
		resUserItems := r.GetBackpackItems(user.ID)
		msg.Text = MessageBackpackUserItems(resUserItems, 0)
		msg.ReplyMarkup = helpers.BackpackInlineKeyboard(resUserItems, 0)
	case "Вещи":
		userItems := r.GetInventoryItems(user.ID)
		msg.Text = MessageGoodsUserItems(user, userItems, 0)
		msg.ReplyMarkup = helpers.GoodsInlineKeyboard(user, userItems, 0)
	case v.GetString("message.emoji.online"):
		userOnline := true
		user = r.User{OnlineMap: &userOnline}.UpdateUser(update)
		msg.Text, buttons = r.GetMyMap(update)
		msg.Text = fmt.Sprintf("%s%sОнлайн включен!", msg.Text, v.GetString("msg_separator"))
		msg.ReplyMarkup = buttons
	case v.GetString("message.emoji.offline"):
		userOnline := false
		user = r.User{OnlineMap: &userOnline}.UpdateUser(update)
		msg.Text, buttons = r.GetMyMap(update)
		msg.Text = fmt.Sprintf("%s%sОнлайн выключен!", msg.Text, v.GetString("msg_separator"))
		msg.ReplyMarkup = buttons
	case ItemHead.View:
		res := directionCell(update, char[1])
		status, text := r.UpdateUserInstrument(update, user, ItemHead)
		if status != "Ok" {
			msg.Text = fmt.Sprintf("%s%s%s", r.ViewItemInfo(res), v.GetString("msg_separator"), text)
		} else {
			msg.Text = r.ViewItemInfo(res)
		}
	case v.GetString("message.emoji.wrench"):
		loc := directionCell(update, char[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		charWorkbench := strings.Fields("workbench usPoint 0 1stComp null 0 2ndComp null 0 3rdComp null 0")
		msg = Workbench(&cell, charWorkbench)
	case v.GetString("message.emoji.quest"):
		loc := directionCell(update, char[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		msg = Quest(&cell, user)
	default:
		msg.Text, buttons = r.GetMyMap(update)
		msg = tg.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n\nНет инструмента в руке!", msg.Text))
		msg.ReplyMarkup = buttons
	}

	return msg
}

func userMenuLocation(update tg.Update, user r.User) tg.MessageConfig {
	newMessage := update.Message.Text
	msg.ChatID = update.Message.Chat.ID

	switch newMessage {
	case "🗺 Карта 🗺":
		msg.Text, msg.ReplyMarkup = r.GetMyMap(update)
		r.User{MenuLocation: "Карта"}.UpdateUser(update)
	case fmt.Sprintf("%s Профиль 👔", user.Avatar):
		msg.Text = r.GetUserInfo(update)
		msg.ReplyMarkup = helpers.ProfileKeyboard(user)
		r.User{MenuLocation: "Профиль"}.UpdateUser(update)
	default:
		msg.Text = "Меню"
		msg.ReplyMarkup = helpers.MainKeyboard(user)
		r.User{MenuLocation: "Меню"}.UpdateUser(update)
	}

	return msg
}

func userMapLocation(update tg.Update, user r.User) tg.MessageConfig {
	newMessage := update.Message.Text
	char := strings.Fields(newMessage)

	if len(char) != 1 {
		msg = useSpecialCell(update, char, user)
	} else {
		msg = useDefaultCell(update, user)
	}

	return msg
}

func userProfileLocation(update tg.Update, user r.User) tg.MessageConfig {
	newMessage := update.Message.Text

	if user.Username == "waiting" {
		r.User{Username: newMessage}.UpdateUser(update)
		msg = tg.NewMessage(update.Message.Chat.ID, r.GetUserInfo(update))
		msg.ReplyMarkup = helpers.ProfileKeyboard(user)
	} else {
		switch newMessage {
		case "📝 Изменить имя? 📝":
			r.User{Username: "waiting"}.UpdateUser(update)
			msg = tg.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
		case fmt.Sprintf("%s Изменить аватар? %s", user.Avatar, user.Avatar):
			msg = tg.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар...")
			msg.ReplyMarkup = helpers.EmojiInlineKeyboard()
		case "/menu", v.GetString("user_location.menu"):
			msg = tg.NewMessage(update.Message.Chat.ID, "Меню")
			msg.ReplyMarkup = helpers.MainKeyboard(user)
			r.User{MenuLocation: "Меню"}.UpdateUser(update)
		default:
			msg = tg.NewMessage(update.Message.Chat.ID, r.GetUserInfo(update))
			msg.ReplyMarkup = helpers.ProfileKeyboard(user)
		}
	}

	return msg
}

func useDefaultCell(update tg.Update, user r.User) tg.MessageConfig {
	newMessage := update.Message.Text
	msg.ChatID = update.Message.Chat.ID
	buttons := tg.ReplyKeyboardMarkup{}
	currentTime := time.Now()

	switch newMessage {
	case v.GetString("message.doing.up"), v.GetString("message.doing.down"), v.GetString("message.doing.left"), v.GetString("message.doing.right"):
		msg = UserMoving(update, user, newMessage)
	case v.GetString("message.emoji.water"):
		msg.Text = "Ты не похож на Jesus! 👮‍♂️"
	case v.GetString("message.emoji.clock"):
		msg.Text = fmt.Sprintf("%s\nЧасики тикают...", currentTime.Format("15:04:05"))
	case user.Avatar:
		msg.Text, buttons = r.GetMyMap(update)
		msg.Text = fmt.Sprintf("%s\n\n%s", r.GetUserInfo(update), msg.Text)
		msg.ReplyMarkup = buttons
	case "/menu", v.GetString("user_location.menu"):
		msg.Text = "Меню"
		msg.ReplyMarkup = helpers.MainKeyboard(user)
		r.User{MenuLocation: "Меню"}.UpdateUser(update)
	case v.GetString("message.emoji.casino"):
		msg.Text = "💰💵🤑 Ставки на JOY CASINO дот COM! 🤑💵💰"
	case v.GetString("message.emoji.forbidden"):
		msg.Text = "🚫 Сюда нельзя! 🚫"
	default:
		msg.Text, buttons = r.GetMyMap(update)
		msg.Text = fmt.Sprintf("%s%sХммм....🤔", msg.Text, v.GetString("msg_separator"))
		msg.ReplyMarkup = buttons
	}

	return msg
}

func directionCell(update tg.Update, direction string) r.Location {
	res := r.GetOrCreateMyLocation(update)

	switch direction {
	case v.GetString("message.doing.up"):
		y := *res.AxisY + 1
		return r.Location{MapsId: res.MapsId, AxisX: res.AxisX, AxisY: &y}
	case v.GetString("message.doing.down"):
		y := *res.AxisY - 1
		return r.Location{MapsId: res.MapsId, AxisX: res.AxisX, AxisY: &y}
	case v.GetString("message.doing.left"):
		x := *res.AxisX - 1
		return r.Location{MapsId: res.MapsId, AxisX: &x, AxisY: res.AxisY}
	case v.GetString("message.doing.right"):
		x := *res.AxisX + 1
		return r.Location{MapsId: res.MapsId, AxisX: &x, AxisY: res.AxisY}
	}
	return res
}

func MessageBackpackUserItems(userItems []r.UserItem, rowUser int) string {
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
		userItemMsg += fmt.Sprintf("%s   %d%s     *HP*:  _%d_ ♥️️     *ST*:  _%d_ \U0001F9C3 ️\n", firstCell, *item.Count, item.Item.View, *item.Item.Healing, *item.Item.Satiety)

	}

	return userItemMsg
}

func MessageGoodsUserItems(user r.User, userItems []r.UserItem, rowUser int) string {
	var userItemMsg = "🧥 *Вещички* 🎒\n\n"
	userItemMsg = messageUserDressedGoods(user) + userItemMsg

	if len(userItems) == 0 {
		return userItemMsg + "👻 \U0001F9B4  Пусто .... 🕸 🕷"
	}

	for i, item := range userItems {
		_, res := user.IsDressedItem(userItems[i])

		if res == v.GetString("callback_char.take_off_good") {
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
		userItemMsg += fmt.Sprintf("%s  %s %dшт.   %s %s   (%d/%d)\n", firstCell, item.Item.View, *item.Count, res, item.Item.Name, *item.CountUseLeft, *item.Item.CountUse)

	}

	return userItemMsg
}

func BackPackMoving(charData []string, update tg.Update) tg.MessageConfig {
	i := r.ToInt(charData[1])
	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	userItems := r.GetBackpackItems(user.ID)

	switch i {
	case len(userItems):
		i = i - 1
	}

	msg.Text = MessageBackpackUserItems(userItems, i)
	msg.ReplyMarkup = helpers.BackpackInlineKeyboard(userItems, i)

	return msg
}

func GoodsMoving(charData []string, update tg.Update) tg.MessageConfig {
	i := r.ToInt(charData[1])

	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	userItems := r.GetInventoryItems(user.ID)

	switch i {
	case len(userItems):
		i = i - 1
	}

	msg.Text = MessageGoodsUserItems(user, userItems, i)
	msg.ReplyMarkup = helpers.GoodsInlineKeyboard(user, userItems, i)

	return msg
}

func UserEatItem(update tg.Update, charData []string) tg.MessageConfig {
	userTgId := r.GetUserTgId(update)
	userItemId := r.ToInt(charData[1])

	user := r.GetUser(r.User{TgId: userTgId})
	userItem := r.UserItem{ID: userItemId}.UserGetUserItem()

	res := userItem.EatItem(update, user)
	charDataForOpenBackPack := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.backpack_moving"), charData[2]))
	msg = BackPackMoving(charDataForOpenBackPack, update)
	msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), res)

	return msg
}

func UserDeleteItem(update tg.Update, charData []string) tg.MessageConfig {
	userItemId := r.ToInt(charData[1])
	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	userItem := r.UserItem{ID: userItemId}.UserGetUserItem()

	countAfterUserThrowOutItem := 0
	var updateUserItemStruct = r.UserItem{
		ID:    userItemId,
		Count: &countAfterUserThrowOutItem,
	}

	user.UpdateUserItem(updateUserItemStruct)

	var charDataForOpenList []string
	switch charData[3] {
	case "good":
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		userTakeOffGood(update, charData)
		msg = GoodsMoving(charDataForOpenList, update)
	case "backpack":
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.backpack_moving"), charData[2]))
		msg = BackPackMoving(charDataForOpenList, update)
	}

	msg.Text = fmt.Sprintf("%s%s🗑 Вы выкинули %s%dшт.", msg.Text, v.GetString("msg_separator"), userItem.Item.View, *userItem.Count)

	return msg
}

func usersHandItemsView(user r.User) (r.Item, r.Item, r.Item) {
	ItemLeftHand := r.Item{View: v.GetString("message.emoji.hand")}
	ItemRightHand := r.Item{View: v.GetString("message.emoji.hand")}
	var ItemHead r.Item

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

func messageUserDressedGoods(user r.User) string {
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

	var messageUserGoods = fmt.Sprintf("〰️☁️〰️〰️☀️\n"+
		"〰️〰️%s〰️〰️\n"+
		"〰️〰️%s〰️〰️\n"+
		"〰️%s%s%s〰️\n"+
		"〰️〰️%s〰️〰️\n"+
		"〰️〰️%s〰️️🍺️\n"+
		"\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\n\n",
		head, user.Avatar, leftHand, body, rightHand, foot, shoes)

	return messageUserGoods
}

func userTakeOffGood(update tg.Update, charData []string) {
	userItemId := r.ToInt(charData[1])
	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	userItem := r.UserItem{ID: userItemId}.UserGetUserItem()

	if user.HeadId != nil && userItem.ItemId == *user.HeadId {
		r.SetNullUserField(update, "head_id")
	} else if user.LeftHandId != nil && userItem.ItemId == *user.LeftHandId {
		r.SetNullUserField(update, "left_hand_id")
	} else if user.RightHandId != nil && userItem.ItemId == *user.RightHandId {
		r.SetNullUserField(update, "right_hand_id")
	} else if user.BodyId != nil && userItem.ItemId == *user.BodyId {
		r.SetNullUserField(update, "body_id")
	} else if user.FootId != nil && userItem.ItemId == *user.FootId {
		r.SetNullUserField(update, "foot_id")
	} else if user.ShoesId != nil && userItem.ItemId == *user.ShoesId {
		r.SetNullUserField(update, "shoes_id")
	}

	charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
	msg = GoodsMoving(charDataForOpenGoods, update)
	msg.Text = fmt.Sprintf("%s%sВещь снята!", msg.Text, v.GetString("msg_separator"))
}

func dressUserItem(update tg.Update, charData []string) tg.MessageConfig {
	userItemId := r.ToInt(charData[1])
	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	userItem := r.UserItem{ID: userItemId}.UserGetUserItem()
	changeHandItem := false

	var result = fmt.Sprintf("Вы надели %s", userItem.Item.View)

	switch *userItem.Item.DressType {
	case "hand":
		if user.LeftHandId == nil {
			r.User{LeftHandId: &userItem.ItemId}.UpdateUser(update)
		} else if user.RightHandId == nil {
			r.User{RightHandId: &userItem.ItemId}.UpdateUser(update)
		} else {
			result = "У вас заняты все руки! Что хочешь снять?"
			changeHandItem = true
		}
	case "head":
		if user.HeadId == nil {
			r.User{HeadId: &userItem.ItemId}.UpdateUser(update)
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "body":
		if user.BodyId == nil {
			r.User{BodyId: &userItem.ItemId}.UpdateUser(update)
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "foot":
		if user.FootId == nil {
			r.User{FootId: &userItem.ItemId}.UpdateUser(update)
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "shoes":
		if user.ShoesId == nil {
			r.User{ShoesId: &userItem.ItemId}.UpdateUser(update)
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	}

	if changeHandItem {
		msg.ReplyMarkup = helpers.ChangeItemInHand(user, userItemId, charData[2])
	} else {
		charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		msg = GoodsMoving(charDataForOpenGoods, update)
	}

	msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), result)

	return msg
}

func userThrowOutItem(update tg.Update, user r.User, charData []string) tg.MessageConfig {
	userItem := r.UserItem{ID: r.ToInt(charData[1])}.UserGetUserItem()

	*userItem.Count = *userItem.Count - r.ToInt(charData[3])

	res := r.UpdateCellUnderUser(update, userItem, r.ToInt(charData[3]))
	var msgtext string
	if res != "Ok" {
		msgtext = fmt.Sprintf("%s%s", v.GetString("msg_separator"), res)
	} else {
		msgtext = fmt.Sprintf("%sВы сбросили %s %sшт. на карту!", v.GetString("msg_separator"), userItem.Item.View, charData[3])
		user.UpdateUserItem(r.UserItem{ID: userItem.ID, Count: userItem.Count})
	}

	var charDataForOpenList []string
	switch charData[4] {
	case "good":
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		if *userItem.Count == 0 {
			userTakeOffGood(update, charData)
		}
		msg = GoodsMoving(charDataForOpenList, update)
	case "backpack":
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.backpack_moving"), charData[2]))
		msg = BackPackMoving(charDataForOpenList, update)
	}

	msg.Text = fmt.Sprintf("%s%s", msg.Text, msgtext)

	return msg
}

func userWantsToThrowOutItem(update tg.Update, charData []string) tg.MessageConfig {
	userItem := r.UserItem{ID: r.ToInt(charData[1])}.UserGetUserItem()

	if userItem.CountUseLeft != nil && *userItem.CountUseLeft != *userItem.Item.CountUse {
		*userItem.Count = *userItem.Count - 1
	}

	if *userItem.Count == 0 {
		var charDataForOpenList []string
		switch charData[3] {
		case "good":
			charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
			if *userItem.CountUseLeft == *userItem.Item.CountUse {
				userTakeOffGood(update, charData)
			}
			msg = GoodsMoving(charDataForOpenList, update)
		case "backpack":
			charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.backpack_moving"), charData[2]))
			msg = BackPackMoving(charDataForOpenList, update)
		}
		msg.Text = fmt.Sprintf("%s%sНельзя выкинуть на карту предмет, который уже был использован!", msg.Text, v.GetString("msg_separator"))
	} else {
		msg.ReplyMarkup = helpers.CountItemUserWantsToThrow(charData, userItem)
		msg.Text = fmt.Sprintf("%sСколько %s ты хочешь скинуть?", v.GetString("msg_separator"), userItem.Item.View)
	}

	return msg
}

func Quest(cell *r.Cell, user r.User) tg.MessageConfig {
	if !cell.IsQuest() {
		msg.Text = v.GetString("error.no_quest_item")
		return msg
	}

	msg.Text = v.GetString("user_location.tasks_menu_message")
	msg.ReplyMarkup = helpers.AllQuestsMessageKeyboard(user)

	return msg
}

func Workbench(cell *r.Cell, char []string) tg.MessageConfig {
	var charData []string
	if !cell.IsWorkbench() {
		msg.Text = "Здесь нет верстака!"
		return msg
	}

	if cell != nil {
		charData = strings.Fields("workbench usPoint 0 1stComp nil 0 2ndComp nil 0 3rdComp nil 0")
	} else {
		charData = strings.Fields(fmt.Sprintf("workbench usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", char[2], char[4], char[5], char[7], char[8], char[10], char[11]))
	}

	msg = OpenWorkbenchMessage(charData)
	msg.ReplyMarkup = helpers.WorkbenchButton(charData)

	return msg
}

func OpenWorkbenchMessage(char []string) tg.MessageConfig {
	// char = workbench usPoint 0 1stComp: id 0 2ndComp id 0 3rdComp id 0

	fstCnt := getViewEmojiForMsg(char, 0)
	secCnt := getViewEmojiForMsg(char, 1)
	trdCnt := getViewEmojiForMsg(char, 2)

	fstComponentView := viewComponent(char[4])
	secComponentView := viewComponent(char[7])
	trdComponentView := viewComponent(char[10])

	cellUser := r.ToInt(char[2])
	userPointer := strings.Fields("〰️ 〰️ 〰️")
	userPointer[cellUser] = "👇"

	msg.Text = fmt.Sprintf(
		"〰️〰️〰️☁️〰️〰️〰️☀️〰️\n"+
			"〰️〰️%s〰️%s〰️%s〰️〰️\n"+
			"🔬〰️%s〰️%s〰️%s〰️📡\n"+
			"\U0001F7EB\U0001F7EB%s\U0001F7EB%s\U0001F7EB%s\U0001F7EB\U0001F7EB\n"+
			"〰️\U0001F7EB〰️〰️〰️〰️〰️\U0001F7EB〰️\n"+
			"〰️\U0001F7EB〰️〰️〰️🍺〰️\U0001F7EB〰️\n"+
			"\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9",
		userPointer[0], userPointer[1], userPointer[2],
		fstComponentView, secComponentView, trdComponentView,
		fstCnt, secCnt, trdCnt)

	return msg
}

func viewComponent(id string) string {
	if id != "nil" {
		component := r.UserItem{ID: r.ToInt(id)}.UserGetUserItem()
		return component.Item.View
	}
	return "⚪"
}

func AllReceiptsMsg() string {
	receipts := r.GetReceipts()
	var msgText string
	for _, receipt := range receipts {
		var fstEl string
		var secEl string
		var trdEl string

		if receipt.Component1ID != nil {
			fstEl = fmt.Sprintf("%d⃣%s", *receipt.Component1Count, receipt.Component1.View)
		}
		if receipt.Component2ID != nil {
			secEl = fmt.Sprintf("➕%d⃣%s", *receipt.Component2Count, receipt.Component2.View)
		}
		if receipt.Component3ID != nil {
			trdEl = fmt.Sprintf("➕%d⃣%s", *receipt.Component3Count, receipt.Component3.View)
		}
		msgText = msgText + fmt.Sprintf("%s 🔚 %s%s%s\n", receipt.ItemResult.View, fstEl, secEl, trdEl)
	}
	return msgText
}

func GetReceiptFromData(char []string) r.Receipt {
	var result r.Receipt

	if char[4] != "nil" && char[5] != "0" {
		fstItem := r.UserItem{ID: r.ToInt(char[4])}.UserGetUserItem()
		id := int(fstItem.Item.ID)
		c := r.ToInt(char[5])

		result.Component1ID = &id
		result.Component1Count = &c
	}

	if char[7] != "nil" && char[8] != "0" {
		fstItem := r.UserItem{ID: r.ToInt(char[7])}.UserGetUserItem()
		id := int(fstItem.Item.ID)
		c := r.ToInt(char[8])

		result.Component2ID = &id
		result.Component2Count = &c
	}

	if char[10] != "nil" && char[11] != "0" {
		fstItem := r.UserItem{ID: r.ToInt(char[10])}.UserGetUserItem()
		id := int(fstItem.Item.ID)
		c := r.ToInt(char[11])

		result.Component3ID = &id
		result.Component3Count = &c
	}

	return result
}

func PutCountComponent(char []string) tg.MessageConfig {
	userItemId := char[r.ToInt(char[2])+(4+r.ToInt(char[2])*2)] // char[x + (4+x*2 )] = char[4]
	userItem := r.UserItem{ID: r.ToInt(userItemId)}.UserGetUserItem()

	msg.ReplyMarkup = helpers.ChangeCountUserItem(char, userItem)

	return msg
}

func UserCraftItem(user r.User, receipt *r.Receipt) (tg.MessageConfig, bool) {
	deletePrevMessage := true
	if receipt == nil {
		msg.Text = "Такого рецепта не существует!"
		msg.ReplyMarkup = nil
		deletePrevMessage = false
		return msg, deletePrevMessage
	}

	msg.ReplyMarkup = nil
	resultItem := r.UserItem{UserId: int(user.ID), ItemId: receipt.ItemResultID}.UserGetUserItem()

	if *receipt.ItemResultCount+*resultItem.Count > *resultItem.Item.MaxCountUserHas {
		msg.Text = fmt.Sprintf("Вы не можете иметь больше, чем %d %s!\nСори... такие правила(", *resultItem.Item.MaxCountUserHas, resultItem.Item.View)
		msg.ReplyMarkup = nil
		deletePrevMessage = false
		return msg, deletePrevMessage
	}

	if receipt.Component1ID != nil && receipt.Component1Count != nil {
		userItem := r.UserItem{UserId: int(user.ID), ItemId: *receipt.Component1ID}.UserGetUserItem()
		countItem1 := *userItem.Count - *receipt.Component1Count
		user.UpdateUserItem(r.UserItem{ID: userItem.ID, ItemId: *receipt.Component1ID, Count: &countItem1}) // CountUseLeft: resultItem.CountUseLeft
	}
	if receipt.Component2ID != nil && receipt.Component2Count != nil {
		userItem := r.UserItem{UserId: int(user.ID), ItemId: *receipt.Component2ID}.UserGetUserItem()
		countItem2 := *userItem.Count - *receipt.Component2Count
		user.UpdateUserItem(r.UserItem{ID: userItem.ID, ItemId: *receipt.Component2ID, Count: &countItem2}) // CountUseLeft: resultItem.CountUseLeft
	}
	if receipt.Component3ID != nil && receipt.Component3Count != nil {
		userItem := r.UserItem{UserId: int(user.ID), ItemId: *receipt.Component3ID}.UserGetUserItem()
		countItem3 := *userItem.Count - *receipt.Component3Count
		user.UpdateUserItem(r.UserItem{ID: userItem.ID, ItemId: *receipt.Component3ID, Count: &countItem3}) // CountUseLeft: resultItem.CountUseLeft
	}

	if *resultItem.Count == 0 {
		resultItem.CountUseLeft = resultItem.Item.CountUse
	}
	*resultItem.Count = *resultItem.Count + *receipt.ItemResultCount
	user.UpdateUserItem(r.UserItem{ID: resultItem.ID, Count: resultItem.Count, CountUseLeft: resultItem.CountUseLeft})

	charData := strings.Fields("workbench usPoint 0 1stComp nil 0 2ndComp nil 0 3rdComp nil 0")
	msg = Workbench(nil, charData)
	msg.Text = fmt.Sprintf("%s%sСупер! Ты получил %s %d шт. %s!", msg.Text, v.GetString("msg_separator"), resultItem.Item.View, *receipt.ItemResultCount, receipt.ItemResult.Name)
	return msg, deletePrevMessage
}

func getViewEmojiForMsg(char []string, i int) string {
	count := i + 5 + i*2

	if char[count] == "0" {
		return "\U0001F7EB"
	}

	return fmt.Sprintf("%s⃣", char[count])
}

func updateUserHand(update tg.Update, char []string, userItem r.UserItem) {
	switch char[0] {
	case v.GetString("callback_char.change_left_hand"):
		r.User{LeftHandId: &userItem.ItemId}.UpdateUser(update)
	case v.GetString("callback_char.change_right_hand"):
		r.User{RightHandId: &userItem.ItemId}.UpdateUser(update)
	}
}

func UserMoving(update tg.Update, user r.User, char string) tg.MessageConfig {
	var text string
	res := directionCell(update, char)

	r.UpdateLocation(update, res)
	lighterMsg := user.CheckUserHasLighter(update)
	if lighterMsg != "Ok" {
		text = fmt.Sprintf("%s%s", v.GetString("msg_separator"), lighterMsg)
	}
	msg.Text, msg.ReplyMarkup = r.GetMyMap(update)
	msg.Text = msg.Text + text

	return msg
}

func UserUseHandOrInstrument(update tg.Update, char []string) tg.MessageConfig {
	res := directionCell(update, char[1])
	resultOfGetItem := r.UserGetItem(update, res, char)
	resText, buttons := r.GetMyMap(update)
	msg.Text = fmt.Sprintf("%s%s%s", resText, v.GetString("msg_separator"), resultOfGetItem)
	msg.ReplyMarkup = buttons

	return msg
}

func OpenQuest(questId uint, user r.User) tg.MessageConfig {
	quest := r.Quest{ID: questId}.GetQuest()
	userQuest := r.UserQuest{UserId: user.ID, QuestId: questId}.GetUserQuest()

	msg.Text = quest.QuestInfo(userQuest)
	msg.ReplyMarkup = helpers.OpenQuestKeyboard(quest, userQuest)

	return msg
}

func UserDoneQuest(questId uint, user r.User) tg.MessageConfig {
	userQuest := r.UserQuest{UserId: user.ID, QuestId: questId}.GetUserQuest()
	if !userQuest.Quest.Task.HasUserDoneTask(user) {
		msg.Text = v.GetString("errors.user_did_not_task")
		return msg
	}

	userQuest.UserDoneQuest(user)
	user.UserGetResult(userQuest.Quest.Result)

	msg = OpenQuest(questId, user)
	msg.Text = fmt.Sprintf("*Задание выполнено!*\n%s", msg.Text)

	return msg
}
