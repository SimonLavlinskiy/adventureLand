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

func userMenuLocation(update tg.Update, user r.User) tg.MessageConfig {
	newMessage := update.Message.Text

	switch newMessage {
	case "🗺 Карта 🗺":
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		r.User{TgId: user.TgId, MenuLocation: "Карта"}.UpdateUser()
	case fmt.Sprintf("%s Профиль 👔", user.Avatar):
		msg.Text = user.GetUserInfo()
		msg.ReplyMarkup = helpers.ProfileKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Профиль"}.UpdateUser()
	default:
		msg.Text = "Меню"
		msg.ReplyMarkup = helpers.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	}

	msg.ChatID = update.Message.Chat.ID
	return msg
}

func userProfileLocation(update tg.Update, user r.User) tg.MessageConfig {
	newMessage := update.Message.Text

	if user.Username == "waiting" {
		r.User{TgId: user.TgId, Username: newMessage}.UpdateUser()
		msg.Text = user.GetUserInfo()
		msg.ReplyMarkup = helpers.ProfileKeyboard(user)
	} else {
		switch newMessage {
		case "📝 Изменить имя? 📝":
			r.User{TgId: user.TgId, Username: "waiting"}.UpdateUser()
			msg.Text = "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️"
			msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
		case fmt.Sprintf("%s Изменить аватар? %s", user.Avatar, user.Avatar):
			msg.Text = "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар..."
			msg.ReplyMarkup = helpers.EmojiInlineKeyboard()
		case "/menu", v.GetString("user_location.menu"):
			msg.Text = "Меню"
			msg.ReplyMarkup = helpers.MainKeyboard(user)
			r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
		default:
			msg.Text = user.GetUserInfo()
			msg.ReplyMarkup = helpers.ProfileKeyboard(user)
		}
	}

	msg.ChatID = update.Message.Chat.ID
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

func useSpecialCell(update tg.Update, char []string, user r.User) tg.MessageConfig {
	ItemLeftHand, ItemRightHand, ItemHead := helpers.UsersHandItemsView(user)

	// При нажатии кнопок
	switch char[0] {
	case v.GetString("message.doing.up"), v.GetString("message.doing.down"), v.GetString("message.doing.left"), v.GetString("message.doing.right"):
		msg = helpers.UserMoving(user, char, char[0])
	case v.GetString("message.emoji.foot"):
		msg = helpers.UserMoving(user, char, char[1])
	case v.GetString("message.emoji.hand"), ItemLeftHand.View, ItemRightHand.View:
		msg = helpers.UserUseHandOrInstrument(user, char)
	case v.GetString("message.emoji.exclamation_mark"):
		cellLocation := helpers.DirectionCell(user, char[3])
		cell := r.Cell{MapsId: *cellLocation.MapsId, AxisX: *cellLocation.AxisX, AxisY: *cellLocation.AxisY}
		cell = cell.GetCell()
		msg.Text = "В зависимости от предмета в твоих руках ты можешь получить разный результат. Выбирай..."
		msg.ReplyMarkup = helpers.ChooseInstrumentKeyboard(char, cell, user)
	case v.GetString("message.emoji.stop_use"):
		msg.Text = v.GetString("errors.user_not_has_item_in_hand")
	case "Рюкзак":
		msg.ReplyMarkup, msg.Text = helpers.BackpackCategoryKeyboard()
	case "Вещи":
		userItems := r.GetInventoryItems(user.ID)
		msg.Text = helpers.MessageGoodsUserItems(user, userItems, 0)
		msg.ReplyMarkup = helpers.GoodsInlineKeyboard(user, userItems, 0)
	case v.GetString("message.emoji.online"):
		userOnline := true
		user = r.User{TgId: user.TgId, OnlineMap: &userOnline}.UpdateUser()
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%sОнлайн включен!", msg.Text, v.GetString("msg_separator"))
	case v.GetString("message.emoji.offline"):
		userOnline := false
		user = r.User{TgId: user.TgId, OnlineMap: &userOnline}.UpdateUser()
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%sОнлайн выключен!", msg.Text, v.GetString("msg_separator"))
	case ItemHead.View:
		res := helpers.DirectionCell(user, char[1])
		text, err := r.UpdateUserInstrument(user, ItemHead)
		if err != nil {
			msg.Text = fmt.Sprintf("%s%s%s", r.ViewItemInfo(res), v.GetString("msg_separator"), text)
		} else {
			msg.Text = r.ViewItemInfo(res)
		}
	case v.GetString("message.emoji.wrench"):
		loc := helpers.DirectionCell(user, char[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		charWorkbench := strings.Fields("workbench usPoint 0 1stComp null 0 2ndComp null 0 3rdComp null 0")
		msg = helpers.Workbench(&cell, charWorkbench)
	case v.GetString("message.emoji.quest"):
		loc := helpers.DirectionCell(user, char[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		msg = helpers.Quest(&cell, user)
	default:
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s\n\nНет инструмента в руке!", msg.Text)
	}

	msg.ChatID = update.Message.Chat.ID
	return msg
}

func useDefaultCell(update tg.Update, user r.User) tg.MessageConfig {
	newMessage := strings.Fields(update.Message.Text)
	currentTime := time.Now()

	// Взаимодействие с предметами на карте, у которых нет действий
	switch newMessage[0] {
	case v.GetString("message.doing.up"), v.GetString("message.doing.down"), v.GetString("message.doing.left"), v.GetString("message.doing.right"):
		msg = helpers.UserMoving(user, newMessage, newMessage[0])
	case v.GetString("message.emoji.water"):
		msg.Text = "Ты не похож на Jesus! 👮‍♂️"
	case v.GetString("message.emoji.clock"):
		msg.Text = fmt.Sprintf("%s\nЧасики тикают...", currentTime.Format("15:04:05"))
	case user.Avatar:
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s\n\n%s", user.GetUserInfo(), msg.Text)
	case "/menu", v.GetString("user_location.menu"):
		msg.Text = "Меню"
		msg.ReplyMarkup = helpers.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	case v.GetString("message.emoji.casino"):
		msg.Text = "💰💵🤑 Ставки на JOY CASINO дот COM! 🤑💵💰"
	case v.GetString("message.emoji.forbidden"):
		msg.Text = "🚫 Сюда нельзя! 🚫"
	default:
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%sХммм....🤔", msg.Text, v.GetString("msg_separator"))
	}

	msg.ChatID = update.Message.Chat.ID
	return msg
}

func callBackResolver(update tg.Update) (tg.MessageConfig, bool) {
	buttons := tg.ReplyKeyboardMarkup{}
	charData := strings.Fields(update.CallbackQuery.Data)
	deletePrevMessage := true

	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	ItemLeftHand, ItemRightHand, ItemHead := helpers.UsersHandItemsView(user)

	if len(charData) == 1 && charData[0] == v.GetString("callback_char.cancel") {
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
	}

	fmt.Println(charData)

	switch charData[0] {

	// Действия в рюкзаке
	case v.GetString("callback_char.category"):
		resUserItems := r.GetBackpackItems(user.ID, charData[1])
		msg.Text = helpers.MessageBackpackUserItems(resUserItems, 0, charData[1])
		msg.ReplyMarkup = helpers.BackpackInlineKeyboard(resUserItems, 0, charData[1])
	case v.GetString("callback_char.backpack_moving"):
		msg = helpers.BackPackMoving(charData, user)
	case v.GetString("callback_char.eat_food"):
		msg = helpers.UserEatItem(user, charData)

	// Действия в инвентаре
	case v.GetString("callback_char.goods_moving"):
		msg = helpers.GoodsMoving(charData, user)
	case v.GetString("callback_char.dress_good"):
		msg = helpers.DressUserItem(user, charData)
	case v.GetString("callback_char.change_left_hand"), v.GetString("callback_char.change_right_hand"):
		user, userItem := r.UpdateUserHand(user, charData)
		charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		msg = helpers.GoodsMoving(charDataForOpenGoods, user)
		msg.Text = fmt.Sprintf("%s%sВы надели %s", msg.Text, v.GetString("msg_separator"), userItem.Item.View)
	case v.GetString("callback_char.take_off_good"):
		msg = helpers.UserTakeOffGood(user, charData)

	// Удаление, выкидывание, описание итема
	case v.GetString("callback_char.delete_item"):
		msg = helpers.UserDeleteItem(user, charData)
	case v.GetString("callback_char.throw_out_item"):
		msg = helpers.UserWantsToThrowOutItem(user, charData)
	case v.GetString("callback_char.count_of_delete"):
		msg = helpers.UserThrowOutItem(user, charData)
	case v.GetString("callback_char.description"):
		msg.Text = r.UserItem{ID: r.ToInt(charData[1])}.GetFullDescriptionOfUserItem()
		msg.ReplyMarkup = helpers.DescriptionInlineButton(charData)

	// Профиль
	case v.GetString("callback_char.change_avatar"):
		res := r.User{TgId: user.TgId, Avatar: charData[1]}.UpdateUser()
		msg.Text = user.GetUserInfo()
		msg.ReplyMarkup = helpers.ProfileKeyboard(res)

	// Крафтинг
	case v.GetString("callback_char.workbench"):
		msg = helpers.Workbench(nil, charData)
	case v.GetString("callback_char.receipt"):
		msg.Text = fmt.Sprintf("📖 *Рецепты*: 📖%s%s", v.GetString("msg_separator"), helpers.AllReceiptsMsg())
		msg.ReplyMarkup = nil
		deletePrevMessage = false
	case v.GetString("callback_char.put_item"):
		userItem := r.GetUserItemsByType(user.ID, strings.Fields("food resource"))
		msg.ReplyMarkup = helpers.ChooseUserItemKeyboard(userItem, charData)
		msg.Text = fmt.Sprintf("%s%sВыбери предмет:", helpers.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
	case v.GetString("callback_char.put_count_item"):
		msg = helpers.PutCountComponent(charData)
		msg.Text = fmt.Sprintf("%s%s⚠️ Сколько выкладываешь?", helpers.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
	case v.GetString("callback_char.make_new_item"):
		resp := helpers.GetReceiptFromData(charData)
		receipt := r.FindReceiptForUser(resp)
		msg, deletePrevMessage = helpers.UserCraftItem(user, receipt)

	// Использование надетых итемов
	case v.GetString("message.emoji.hand"), ItemLeftHand.View, ItemRightHand.View:
		msg = helpers.UserUseHandOrInstrument(user, charData)
	case v.GetString("message.emoji.foot"):
		msg = helpers.UserMoving(user, charData, charData[1])
	case ItemHead.View:
		res := helpers.DirectionCell(user, charData[1])
		text, err := r.UpdateUserInstrument(user, ItemHead)
		if err != nil {
			msg.Text = fmt.Sprintf("%s%s%s", r.ViewItemInfo(res), v.GetString("msg_separator"), text)
		} else {
			msg.Text = r.ViewItemInfo(res)
		}

	// Квесты
	case v.GetString("callback_char.quests"):
		msg.Text = v.GetString("user_location.tasks_menu_message")
		msg.ReplyMarkup = helpers.AllQuestsMessageKeyboard(user)
	case v.GetString("callback_char.quest"):
		msg = helpers.OpenQuest(uint(r.ToInt(charData[1])), user)
	case v.GetString("callback_char.user_get_quest"):
		r.UserQuest{
			UserId:  user.ID,
			QuestId: uint(r.ToInt(charData[1])),
		}.GetOrCreateUserQuest()
		msg = helpers.OpenQuest(uint(r.ToInt(charData[1])), user)
	case v.GetString("callback_char.user_done_quest"):
		msg = helpers.UserDoneQuest(uint(r.ToInt(charData[1])), user)

	// Дом юзера
	case v.GetString("callback_char.buy_home"):
		err := user.CreateUserHouse()
		text := "Поздравляю с покупкой дома!"

		switch err.Error() {
		case "user doesn't have money enough":
			text = "Не хватает деняк! Прийдется еще поднакопить :( "
		default:
			text = "Не получилось :("
		}

		msg.Text, buttons = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), text)
		msg.ReplyMarkup = buttons
	}

	msg.ParseMode = "markdown"
	msg.ChatID = update.CallbackQuery.Message.Chat.ID
	return msg, deletePrevMessage
}
