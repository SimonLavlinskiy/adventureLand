package handlers

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	r "project0/repository"
	s "project0/services"
	"strings"
	"time"
)

func messageResolver(update tg.Update) (msg tg.MessageConfig) {
	var btns tg.InlineKeyboardMarkup
	user := r.GetOrCreateUser(update)

	fmt.Println(user.Username + " делает действие!")

	switch user.MenuLocation {
	case v.GetString("user_location.menu"):
		msg, btns = userMenuLocation(update, user)
	case v.GetString("user_location.maps"):
		msg, btns = userMapLocation(update, user)
	case v.GetString("user_location.profile"):
		msg.Text, btns = userProfileLocation(update, user)
	case v.GetString("user_location.wordle"):
		msg.Text, btns = gameWordle(update, user)
	default:
		msg, btns = userMenuLocation(update, user)
	}

	msg = tg.NewMessage(update.Message.From.ID, msg.Text)
	msg.ReplyMarkup = btns
	msg.ParseMode = "markdown"

	return msg
}

func userMenuLocation(update tg.Update, user r.User) (msg tg.MessageConfig, buttons tg.InlineKeyboardMarkup) {
	newMessage := update.Message.Text

	switch newMessage {
	case "/start":
		msg.Text = v.GetString("main_info.start_msg")
	default:
		msg.Text = "Меню"
		buttons = s.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	}

	return msg, buttons
}

func userProfileLocation(update tg.Update, user r.User) (msgText string, buttons tg.InlineKeyboardMarkup) {
	var newMessage string
	if update.Message != nil {
		newMessage = update.Message.Text
	} else {
		newMessage = update.CallbackQuery.Data
	}

	if user.Username == "waiting" {
		user = r.User{TgId: user.TgId, Username: newMessage}.UpdateUser()
		msgText = user.GetUserInfo()
		buttons = s.ProfileKeyboard(user)
	} else {
		switch newMessage {
		case "/menu", v.GetString("user_location.menu"):
			msgText = "Меню"
			buttons = s.MainKeyboard(user)
			r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
		default:
			msgText = user.GetUserInfo()
			buttons = s.ProfileKeyboard(user)
		}
	}

	return msgText, buttons
}

func userMapLocation(update tg.Update, user r.User) (msg tg.MessageConfig, buttons tg.InlineKeyboardMarkup) {
	newMessage := update.Message.Text

	fmt.Println(newMessage)

	if newMessage == "/menu" || newMessage == v.GetString("user_location.menu") {
		msg.Text = "Меню:"
		buttons = s.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	} else {
		msg.Text, buttons = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%s🤨 Не пойму... 🧐", msg.Text, v.GetString("msg_separator"))
	}

	return msg, buttons
}

func callBackResolver(update tg.Update) (msg tg.EditMessageTextConfig, buttons tg.EditMessageReplyMarkupConfig) {
	var btns tg.InlineKeyboardMarkup
	var msg1, msg2 string
	t := time.Now()

	charData := strings.Fields(update.CallbackQuery.Data)

	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	ItemLeftHand, ItemRightHand, ItemHead := s.UsersHandItemsView(user)

	if len(charData) == 1 && charData[0] == v.GetString("callback_char.cancel") {
		msg.Text, btns = r.GetMyMap(user)
		user = r.User{TgId: user.TgId, MenuLocation: "Карта"}.UpdateUser()
	}

	fmt.Println(charData)

	switch user.MenuLocation {
	case "wordle":
		msg.Text, btns = gameWordle(update, user)
		msg = tg.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, msg.Text)
		buttons = tg.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, btns)
		msg.ParseMode = "markdown"
		return msg, buttons
	case "Меню":
		switch update.CallbackQuery.Data {
		case "/menu", v.GetString("user_location.menu"):
			msg.Text = "Меню:"
			btns = s.MainKeyboard(user)
			r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
		case "🗺 Карта 🗺":
			msg.Text, btns = r.GetMyMap(user)
			r.User{TgId: user.TgId, MenuLocation: "Карта"}.UpdateUser()
		case fmt.Sprintf("%s Профиль 👔", user.Avatar):
			msg.Text = user.GetUserInfo()
			btns = s.ProfileKeyboard(user)
			r.User{TgId: user.TgId, MenuLocation: "Профиль"}.UpdateUser()
		}
		msg = tg.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, msg.Text)
		buttons = tg.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, btns)
		msg.ParseMode = "markdown"

		return msg, buttons
	case "Профиль":
		if strings.Contains(update.CallbackQuery.Data, v.GetString("callback_char.change_avatar")) {
			res := r.User{TgId: user.TgId, Avatar: charData[1]}.UpdateUser()
			msg.Text, btns = userProfileLocation(update, res)
		}

		switch update.CallbackQuery.Data {
		case "cancelChangeName":
			user = r.User{TgId: user.TgId, Username: update.CallbackQuery.From.UserName}.UpdateUser()
			msg.Text, btns = userProfileLocation(update, user)
		case "📝 Изменить имя? 📝":
			r.User{TgId: user.TgId, Username: "waiting"}.UpdateUser()
			msg.Text = "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️"
			btns = s.CancelChangeNameButton(update.CallbackQuery.From.UserName)
		case "avatarList":
			msg.Text = "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар..."
			btns = s.EmojiInlineKeyboard()
		case "/menu", v.GetString("user_location.menu"):
			msg.Text = "Меню:"
			btns = s.MainKeyboard(user)
			r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
		}

		msg = tg.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, msg.Text)
		buttons = tg.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, btns)
		msg.ParseMode = "markdown"

		return msg, buttons
	}

	switch charData[0] {

	//Действия на карте
	case v.GetString("message.doing.up"), v.GetString("message.doing.down"), v.GetString("message.doing.left"), v.GetString("message.doing.right"):
		msg, btns = s.UserMoving(user, charData, charData[0])

	case v.GetString("message.emoji.exclamation_mark"):
		msgMap, _ := r.GetMyMap(user)
		cell := s.DirectionCell(user, charData[3])
		msg.Text, btns = s.ChoseInstrumentMessage(user, charData, cell)
		msg.Text = fmt.Sprintf("%s%s%s", msgMap, v.GetString("msg_separator"), msg.Text)

	case v.GetString("message.emoji.stop_use"):
		msg.Text = v.GetString("errors.user_not_has_item_in_hand")

	case user.Avatar:
		msg.Text, btns = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s\n\n%s", user.GetUserInfo(), msg.Text)

	// Действия в рюкзаке
	case v.GetString("callback_char.category"):
		if len(charData) == 1 {
			msg.Text, btns = s.BackpackCategoryKeyboard()
		} else {
			resUserItems := r.GetBackpackItems(user.ID, charData[1])
			msg.Text = s.MessageBackpackUserItems(user, resUserItems, 0, charData[1])
			btns = s.BackpackInlineKeyboard(resUserItems, 0, charData[1])
		}
	case v.GetString("callback_char.backpack_moving"):
		msg.Text, btns = s.BackPackMoving(charData, user)
	case v.GetString("callback_char.eat_food"):
		msg.Text, btns = s.UserEatItem(user, charData)

	// Действия в инвентаре
	case v.GetString("callback_char.goods_moving"):
		if len(charData) == 1 {
			userItems := r.GetInventoryItems(user.ID)
			msg.Text = s.MessageGoodsUserItems(user, userItems, 0)
			btns = s.GoodsInlineKeyboard(user, userItems, 0)
		} else {
			msg.Text, btns = s.GoodsMoving(charData, user)
		}
	case v.GetString("callback_char.dress_good"):
		msg.Text, btns = s.DressUserItem(user, charData)
	case v.GetString("callback_char.change_left_hand"), v.GetString("callback_char.change_right_hand"):
		user, userItem := r.UpdateUserHand(user, charData)
		charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		msg.Text, btns = s.GoodsMoving(charDataForOpenGoods, user)
		msg.Text = fmt.Sprintf("%s%sВы надели %s", msg.Text, v.GetString("msg_separator"), userItem.Item.View)
	case v.GetString("callback_char.take_off_good"):
		msg.Text, btns = s.UserTakeOffGood(user, charData)

	// Удаление, выкидывание, описание итема
	case v.GetString("callback_char.delete_item"):
		msg.Text, btns = s.UserDeleteItem(user, charData)
	case v.GetString("callback_char.count_of_throw_out"):
		msg.Text, btns = s.UserWantsToThrowOutItem(user, charData)
	case v.GetString("callback_char.throw_out_item"):
		msg.Text, btns = s.UserThrowOutItem(user, charData)
	case v.GetString("callback_char.description"):
		msg.Text = r.UserItem{ID: r.ToInt(charData[1])}.GetFullDescriptionOfUserItem()
		btns = s.DescriptionInlineButton(charData)

	// Крафтинг
	case v.GetString("callback_char.workbench"):
		msg.Text, btns = s.Workbench(nil, charData)
	case v.GetString("callback_char.receipt"):
		msg.Text = fmt.Sprintf("📖 *Рецепты*: 📖%s%s", v.GetString("msg_separator"), s.AllReceiptsMsg())
		btns = s.ReturnToWorkbench(charData)
	case v.GetString("callback_char.put_item"):
		userItem := r.GetUserItemsByType(user.ID, strings.Fields("food resource"))
		btns = s.ChooseUserItemKeyboard(userItem, charData)
		msg.Text = fmt.Sprintf("%s%sВыбери предмет:", s.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
	case v.GetString("callback_char.put_count_item"):
		btns = s.PutCountComponent(charData)
		msg.Text = fmt.Sprintf("%s%s⚠️ Сколько выкладываешь?", s.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
	case v.GetString("callback_char.make_new_item"):
		resp := s.GetReceiptFromData(charData)
		receipt := r.FindReceiptForUser(resp)
		msg.Text, btns = s.UserCraftItem(user, receipt, charData)

	// Использование надетых итемов
	case v.GetString("message.emoji.hand"), ItemLeftHand.View, ItemRightHand.View, v.GetString("message.emoji.fist"):
		msg1, _ = s.UserUseHandOrInstrumentMessage(user, charData)
		res := s.DirectionCell(user, charData[1])
		_, btns = s.ChoseInstrumentMessage(user, charData, res)
		fmt.Printf("%s\n\n%s", msg1, msg2)
		msg.Text = fmt.Sprintf("%s", msg1)
	case v.GetString("message.emoji.foot"):
		msg, btns = s.UserMoving(user, charData, charData[1])
	case ItemHead.View:
		res := s.DirectionCell(user, charData[1])
		text, err := r.UpdateUserInstrument(user, ItemHead)
		msg.Text = r.ViewItemInfo(res)
		if err != nil {
			msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), text)
		}
		btns = s.CancelButton()

	// Квесты
	case v.GetString("callback_char.quests"):
		msg.Text = v.GetString("user_location.tasks_menu_message")
		btns = s.AllQuestsMessageKeyboard(user)
	case v.GetString("callback_char.quest"):
		msg.Text, btns = s.OpenQuest(uint(r.ToInt(charData[1])), user)
	case v.GetString("callback_char.user_get_quest"):
		r.UserQuest{
			UserId:  user.ID,
			QuestId: uint(r.ToInt(charData[1])),
		}.GetOrCreateUserQuest()
		msg.Text, btns = s.OpenQuest(uint(r.ToInt(charData[1])), user)
	case v.GetString("callback_char.user_done_quest"):
		msg.Text, btns = s.UserDoneQuest(uint(r.ToInt(charData[1])), user)

	// Дом юзера
	case v.GetString("callback_char.buy_home"):
		err := user.CreateUserHouse()
		text := "Поздравляю с покупкой дома!"

		if err != nil {
			switch err.Error() {
			case "user doesn't have money enough":
				text = "Не хватает деняк! Прийдется еще поднакопить :( "
			default:
				text = "Не получилось :("
			}
		}

		msg.Text, btns = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), text)

	// Чатик
	case v.GetString("message.emoji.chat"):
		loc := s.DirectionCell(user, charData[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		msg.Text, btns = s.OpenChatKeyboard(cell, user)
	case v.GetString("callback_char.join_to_chat"):
		ui := make([]r.ChatUser, 1)
		ui[0] = r.Chat{ID: uint(r.ToInt(charData[1]))}.GetOrCreateChatUser(user)
		cell := r.Cell{ID: uint(r.ToInt(charData[3]))}.GetCell()
		msg.Text, btns = s.OpenChatKeyboard(cell, user)

		s.NotifyUsers(ui, v.GetString("main_info.message_user_sign_in_chat"))

		// ивент итемы
	case v.GetString("message.emoji.wrench"):
		loc := s.DirectionCell(user, charData[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		charWorkbench := strings.Fields("workbench usPoint 0 1stComp null 1 2ndComp null 1 3rdComp null 1")
		msg.Text, btns = s.Workbench(&cell, charWorkbench)
	case v.GetString("message.emoji.quest"):
		loc := s.DirectionCell(user, charData[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		msg.Text, btns = s.Quest(&cell, user)
	case v.GetString("message.emoji.wordle_game"):
		r.User{TgId: user.TgId, MenuLocation: "wordle"}.UpdateUser()
		msg.Text, btns = s.WordleMap(user)

		// Взаимодействие с предметами на карте, у которых нет действий
	case v.GetString("message.emoji.water"):
		msg.Text, btns = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%sТы не похож на Jesus! 👮", msg.Text, v.GetString("msg_separator"))
	case v.GetString("message.emoji.clock"):
		msg.Text, btns = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s\nЧасики тикают...", t.Format("15:04:05"))

	case v.GetString("message.emoji.casino"):
		msg.Text, btns = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%s💰💵🤑 Ставки на JOY CASINO дот COM! 🤑💵💰", msg.Text, v.GetString("msg_separator"))
	case v.GetString("message.emoji.forbidden"):
		msg.Text, btns = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%s🚫 Сюда нельзя! 🚫\"", msg.Text, v.GetString("msg_separator"))
	case "👨‍🔧":
		msg.Text, btns = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%s👨‍🔧 Зачем зашел за кассу? 😑", msg.Text, v.GetString("msg_separator"))
	case "/menu", v.GetString("user_location.menu"):
		msg.Text = "Меню:"
		btns = s.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	case "cancel":
		msg.Text, btns = r.GetMyMap(user)
	default:
		msg.Text, btns = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%sХммм....🤔", msg.Text, v.GetString("msg_separator"))
	}

	msg = tg.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, msg.Text)
	buttons = tg.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, btns)
	msg.ParseMode = "markdown"

	return msg, buttons
}

func SendUserMessageAllChatUsers(update tg.Update) ([]r.ChatUser, string) {
	user := r.GetOrCreateUser(update)
	chUser := r.GetChatOfUser(user)
	chatUsers := r.Chat{ID: chUser.ChatID}.GetChatUsers()

	var chUsWithoutSender []r.ChatUser
	for _, chatUser := range chatUsers {
		if chatUser.User.TgId != uint(update.Message.From.ID) {
			chUsWithoutSender = append(chUsWithoutSender, chatUser)
		}
	}

	replacer := strings.NewReplacer(
		"/start", fmt.Sprintf("<i>%s</i> %s <code>присоединился к чатику<code>", user.Avatar, user.Username),
	)
	userMsg := replacer.Replace(update.Message.Text)

	message := fmt.Sprintf("<code>От %s %s %s</code>%s%s", user.Avatar, user.Username, user.Avatar, v.GetString("msg_separator"), userMsg)

	return chUsWithoutSender, message
}

func gameWordle(update tg.Update, user r.User) (msgText string, btns tg.InlineKeyboardMarkup) {

	if update.CallbackQuery == nil {
		msgText, btns = s.UserSendNextWord(user, update.Message.Text)
		return msgText, btns
	}

	charData := strings.Fields(update.CallbackQuery.Data)

	switch charData[0] {
	case v.GetString("callback_char.wordle_regulations"):
		msgText = v.GetString("wordle.regulations")
		btns = s.WordleExitButton()
	case "wordleUserStatistic":
		msgText = r.GetWordleUserStatistic(user)
		btns = s.WordleExitButton()
	case "wordleMenu":
		msgText, btns = s.WordleMap(user)
	}

	return msgText, btns
}
