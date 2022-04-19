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

	char := update.CallbackQuery.Data
	charData := strings.Fields(update.CallbackQuery.Data)

	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})

	if len(charData) == 1 && charData[0] == v.GetString("callback_char.cancel") {
		msg.Text, btns = r.GetMyMap(user)
		user = r.User{TgId: user.TgId, MenuLocation: "Карта"}.UpdateUser()
	}

	fmt.Println(charData)

	switch user.MenuLocation {
	case "wordle":
		msg.Text, btns = gameWordle(update, user)
	case "Меню":
		msg.Text, btns = menu(update, user)
	case "Профиль":
		msg.Text, btns = profile(update, user, charData)
	case "Карта":
		msg.Text, btns = mapsDoing(user, char)
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

func useCellWithoutDoing(user r.User, text string) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg, buttons = r.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), text)
	return msg, buttons
}

func openWordle(user r.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	r.User{TgId: user.TgId, MenuLocation: "wordle"}.UpdateUser()
	msg, buttons = s.WordleMap(user)
	return msg, buttons
}

func joinToChat(user r.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	ui := make([]r.ChatUser, 1)
	ui[0] = r.Chat{ID: uint(r.ToInt(charData[1]))}.GetOrCreateChatUser(user)
	cell := r.Cell{ID: uint(r.ToInt(charData[3]))}.GetCell()
	msg, buttons = s.OpenChatKeyboard(cell, user)

	s.NotifyUsers(ui, v.GetString("main_info.message_user_sign_in_chat"))
	return msg, buttons
}

func buyHome(user r.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	text := "Поздравляю с покупкой дома!"
	err := user.CreateUserHouse()
	if err != nil {
		switch err.Error() {
		case "user doesn't have money enough":
			text = "Не хватает деняк! Прийдется еще поднакопить :( "
		default:
			text = "Не получилось :("
		}
	}

	msg, buttons = r.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), text)
	return msg, buttons
}

func userGetQuest(user r.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	r.UserQuest{
		UserId:  user.ID,
		QuestId: uint(r.ToInt(charData[1])),
	}.GetOrCreateUserQuest()
	msg, buttons = s.OpenQuest(uint(r.ToInt(charData[1])), user)
	return msg, buttons
}

func listOfQuests(user r.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg = v.GetString("user_location.tasks_menu_message")
	buttons = s.AllQuestsMessageKeyboard(user)
	return msg, buttons
}

func userHeadItem(user r.User, cell r.Cell, ItemHead r.Item) (msg string, buttons tg.InlineKeyboardMarkup) {
	text, err := r.UpdateUserInstrument(user, ItemHead)
	msg = r.ViewItemInfo(cell)
	if err != nil {
		msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), text)
	}
	buttons = s.CancelButton()
	return msg, buttons
}

func useHandOrInstrument(user r.User, charData []string, cell r.Cell) (msg string, buttons tg.InlineKeyboardMarkup) {
	resultOfGetItem := r.UserGetItemUpdateModels(user, cell, charData)

	msgMap, _ := r.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msgMap, v.GetString("msg_separator"), resultOfGetItem)

	newCell := r.Cell{ID: cell.ID}.GetCell()
	_, buttons = s.ChoseInstrumentMessage(user, newCell)

	return msg, buttons
}

func userGetBox(user r.User, cell r.Cell) (msg string, buttons tg.InlineKeyboardMarkup) {
	var resultMsg string

	for _, instrument := range cell.Item.Instruments {
		if instrument.Type == "get" {
			user.UserGetResult(*instrument.Result)
			resultMsg = s.UserGetResultMsg(*instrument.Result)
		}
	}

	r.UserBox{BoxId: cell.Item.ID, UserId: user.ID}.CreateUserBox()

	msg, buttons = r.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), resultMsg)
	return msg, buttons
}

func craftItem(user r.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	resp := s.GetReceiptFromData(charData)
	receipt := r.FindReceiptForUser(resp)
	msg, buttons = s.UserCraftItem(user, receipt, charData)
	return msg, buttons
}

func changeCountComponent(charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	buttons = s.PutCountComponent(charData)
	msg = fmt.Sprintf("%s%s⚠️ Сколько выкладываешь?", s.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
	return msg, buttons
}

func changeComponent(user r.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	userItem := r.GetUserItemsByType(user.ID, strings.Fields("food resource"))
	buttons = s.ChooseUserItemKeyboard(userItem, charData)
	msg = fmt.Sprintf("%s%sВыбери предмет:", s.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
	return msg, buttons
}

func listOfReceipt(charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg = fmt.Sprintf("📖 *Рецепты*: 📖%s%s", v.GetString("msg_separator"), s.AllReceiptsMsg())
	buttons = s.ReturnToWorkbench(charData)
	return msg, buttons
}

func changeHand(user r.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	user, userItem := r.UpdateUserHand(user, charData)
	charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
	msg, buttons = s.GoodsMoving(charDataForOpenGoods, user)
	msg = fmt.Sprintf("%s%sВы надели %s", msg, v.GetString("msg_separator"), userItem.Item.View)
	return msg, buttons
}

func listOfGoods(user r.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	if len(charData) == 1 {
		userItems := r.GetInventoryItems(user.ID)
		msg = s.MessageGoodsUserItems(user, userItems, 0)
		buttons = s.GoodsInlineKeyboard(user, userItems, 0)
	} else {
		msg, buttons = s.GoodsMoving(charData, user)
	}
	return msg, buttons
}

func listOfBackpackItems(user r.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	if len(charData) == 1 {
		msg, buttons = s.BackpackCategoryKeyboard()
	} else {
		resUserItems := r.GetBackpackItems(user.ID, charData[1])
		msg = s.MessageBackpackUserItems(user, resUserItems, 0, charData[1])
		buttons = s.BackpackInlineKeyboard(resUserItems, 0, charData[1])
	}
	return msg, buttons
}

func mapWithUserInfo(user r.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg, buttons = r.GetMyMap(user)
	msg = fmt.Sprintf("%s\n\n%s", user.GetUserInfo(), msg)
	return msg, buttons
}

func userTouchItem(user r.User, cell r.Cell) (msg string, buttons tg.InlineKeyboardMarkup) {
	msgMap, _ := r.GetMyMap(user)
	msg, buttons = s.ChoseInstrumentMessage(user, cell)
	msg = fmt.Sprintf("%s%s%s", msgMap, v.GetString("msg_separator"), msg)
	return msg, buttons
}

func profile(update tg.Update, user r.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	if strings.Contains(update.CallbackQuery.Data, v.GetString("callback_char.change_avatar")) {
		res := r.User{TgId: user.TgId, Avatar: charData[1]}.UpdateUser()
		msg, buttons = userProfileLocation(update, res)
	}

	switch update.CallbackQuery.Data {
	case "cancelChangeName":
		user = r.User{TgId: user.TgId, Username: update.CallbackQuery.From.UserName}.UpdateUser()
		msg, buttons = userProfileLocation(update, user)
	case "📝 Изменить имя? 📝":
		r.User{TgId: user.TgId, Username: "waiting"}.UpdateUser()
		msg = "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️"
		buttons = s.CancelChangeNameButton(update.CallbackQuery.From.UserName)
	case "avatarList":
		msg = "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар..."
		buttons = s.EmojiInlineKeyboard()
	case "/menu", v.GetString("user_location.menu"):
		msg = "Меню:"
		buttons = s.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	}

	return msg, buttons
}

func menu(update tg.Update, user r.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	switch update.CallbackQuery.Data {
	case "/menu", v.GetString("user_location.menu"):
		msg = "Меню:"
		buttons = s.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	case "🗺 Карта 🗺":
		msg, buttons = r.GetMyMap(user)
		r.User{TgId: user.TgId, MenuLocation: "Карта"}.UpdateUser()
	case fmt.Sprintf("%s Профиль 👔", user.Avatar):
		msg = user.GetUserInfo()
		buttons = s.ProfileKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Профиль"}.UpdateUser()
	}

	return msg, buttons
}

func mapsDoing(user r.User, char string) (msg string, buttons tg.InlineKeyboardMarkup) {
	t := time.Now()
	charData := strings.Fields(char)
	_, _, ItemHead := s.UsersHandItemsView(user)

	switch charData[0] {

	// Действия/кнопки  на карте
	case "move":
		cell := r.Cell{ID: uint(r.ToInt(charData[1]))}.GetCell()
		msg, buttons = s.UserMoving(user, cell)
	case "chooseInstrument":
		cell := r.Cell{ID: uint(r.ToInt(charData[1]))}.GetCell()
		msg, buttons = userTouchItem(user, cell)
	case v.GetString("message.emoji.stop_use"):
		msg = v.GetString("errors.user_not_has_item_in_hand")
	case user.Avatar:
		msg, buttons = mapWithUserInfo(user)
	case "box":
		cell := r.Cell{ID: uint(r.ToInt(charData[1]))}.GetCell()
		msg, buttons = userGetBox(user, cell)

	// Использование надетых итемов
	case "hand", "fist", "item":
		cell := r.Cell{ID: uint(r.ToInt(charData[1]))}.GetCell()
		msg, buttons = useHandOrInstrument(user, charData, cell)
	case "step":
		cell := r.Cell{ID: uint(r.ToInt(charData[1]))}.GetCell()
		msg, buttons = s.UserMoving(user, cell)
	case "head":
		cell := r.Cell{ID: uint(r.ToInt(charData[1]))}.GetCell()
		msg, buttons = userHeadItem(user, cell, ItemHead)

	// Действия в рюкзаке
	case v.GetString("callback_char.category"):
		msg, buttons = listOfBackpackItems(user, charData)
	case v.GetString("callback_char.backpack_moving"):
		msg, buttons = s.BackPackMoving(charData, user)
	case v.GetString("callback_char.eat_food"):
		msg, buttons = s.UserEatItem(user, charData)

	// Действия в инвентаре
	case v.GetString("callback_char.goods_moving"):
		msg, buttons = listOfGoods(user, charData)
	case v.GetString("callback_char.dress_good"):
		msg, buttons = s.DressUserItem(user, charData)
	case v.GetString("callback_char.change_left_hand"), v.GetString("callback_char.change_right_hand"):
		msg, buttons = changeHand(user, charData)
	case v.GetString("callback_char.take_off_good"):
		msg, buttons = s.UserTakeOffGood(user, charData)

	// Удаление, выкидывание, описание итема
	case v.GetString("callback_char.delete_item"):
		msg, buttons = s.UserDeleteItem(user, charData)
	case v.GetString("callback_char.count_of_throw_out"):
		msg, buttons = s.UserWantsToThrowOutItem(user, charData)
	case v.GetString("callback_char.throw_out_item"):
		msg, buttons = s.UserThrowOutItem(user, charData)
	case v.GetString("callback_char.description"):
		msg = r.UserItem{ID: r.ToInt(charData[1])}.GetFullDescriptionOfUserItem()
		buttons = s.DescriptionInlineButton(charData)

	// Крафтинг
	case "wrench":
		cell := r.Cell{ID: uint(r.ToInt(charData[1]))}.GetCell()
		msg, buttons = s.Workbench(&cell, charData)
	case v.GetString("callback_char.workbench"):
		msg, buttons = s.Workbench(nil, charData)
	case v.GetString("callback_char.receipt"):
		msg, buttons = listOfReceipt(charData)
	case v.GetString("callback_char.put_item"):
		msg, buttons = changeComponent(user, charData)
	case v.GetString("callback_char.put_count_item"):
		msg, buttons = changeCountComponent(charData)
	case v.GetString("callback_char.make_new_item"):
		msg, buttons = craftItem(user, charData)

	// Квесты
	case "quests":
		msg, buttons = listOfQuests(user)
	case v.GetString("callback_char.quest"):
		msg, buttons = s.OpenQuest(uint(r.ToInt(charData[1])), user)
	case v.GetString("callback_char.user_get_quest"):
		msg, buttons = userGetQuest(user, charData)
	case v.GetString("callback_char.user_done_quest"):
		msg, buttons = s.UserDoneQuest(uint(r.ToInt(charData[1])), user)

	// Дом юзера
	case v.GetString("callback_char.buy_home"):
		msg, buttons = buyHome(user)

	// Чатик
	case "chat":
		cell := r.Cell{ID: uint(r.ToInt(charData[1]))}.GetCell()
		msg, buttons = s.OpenChatKeyboard(cell, user)
	case v.GetString("callback_char.join_to_chat"):
		msg, buttons = joinToChat(user, charData)

		// вордле
	case "wordle_game":
		msg, buttons = openWordle(user)

	// Взаимодействие с предметами на карте, у которых нет действий
	case v.GetString("message.emoji.water"):
		msg, buttons = useCellWithoutDoing(user, "Ты не похож на Jesus! 👮")
	case v.GetString("message.emoji.clock"):
		msg, buttons = useCellWithoutDoing(user, fmt.Sprintf("%s\nЧасики тикают...", t.Format("15:04:05")))
	case v.GetString("message.emoji.casino"):
		msg, buttons = useCellWithoutDoing(user, "💰💵🤑 Ставки на JOY CASINO дот COM! 🤑💵💰")
	case v.GetString("message.emoji.forbidden"):
		msg, buttons = useCellWithoutDoing(user, "🚫 Сюда нельзя! 🚫")
	case v.GetString("message.emoji.cassir"):
		msg, buttons = useCellWithoutDoing(user, "‍🔧 Зачем зашел за кассу? 😑")

	case "/menu", v.GetString("user_location.menu"):
		msg = "Меню:"
		buttons = s.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	case "cancel":
		msg, buttons = r.GetMyMap(user)
	default:
		msg, buttons = r.GetMyMap(user)
		msg = fmt.Sprintf("%s%sХммм....🤔", msg, v.GetString("msg_separator"))
	}

	return msg, buttons
}
