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

var msgs []tg.MessageConfig

func messageResolver(update tg.Update) []tg.MessageConfig {
	msgs = []tg.MessageConfig{}
	user := r.GetOrCreateUser(update)

	fmt.Println(user.Username + " делает действие!")

	switch user.MenuLocation {
	case v.GetString("user_location.menu"):
		msgs = userMenuLocation(update, user)
	case v.GetString("user_location.maps"):
		msgs = userMapLocation(update, user)
	case v.GetString("user_location.profile"):
		msgs = userProfileLocation(update, user)
	default:
		msgs = userMenuLocation(update, user)
	}

	for i := range msgs {
		msgs[i].ParseMode = "markdown"
		msgs[i].ChatID = update.Message.Chat.ID
	}

	return msgs
}

func userMenuLocation(update tg.Update, user r.User) []tg.MessageConfig {
	var msg tg.MessageConfig
	newMessage := update.Message.Text

	switch newMessage {
	case "🗺 Карта 🗺":
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msgs = append(msgs, msg)
		r.User{TgId: user.TgId, MenuLocation: "Карта"}.UpdateUser()
	case fmt.Sprintf("%s Профиль 👔", user.Avatar):
		msg.Text = user.GetUserInfo()
		msg.ReplyMarkup = s.ProfileKeyboard(user)
		msgs = append(msgs, msg)
		r.User{TgId: user.TgId, MenuLocation: "Профиль"}.UpdateUser()
	default:
		msg.Text = "Меню"
		msg.ReplyMarkup = s.MainKeyboard(user)
		msgs = append(msgs, msg)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
	}

	for i := range msgs {
		msgs[i].ChatID = update.Message.Chat.ID
	}
	return msgs
}

func userProfileLocation(update tg.Update, user r.User) []tg.MessageConfig {
	var msg tg.MessageConfig
	newMessage := update.Message.Text

	if user.Username == "waiting" {
		r.User{TgId: user.TgId, Username: newMessage}.UpdateUser()
		msg.Text = user.GetUserInfo()
		msg.ReplyMarkup = s.ProfileKeyboard(user)
		msgs = append(msgs, msg)
	} else {
		switch newMessage {
		case "📝 Изменить имя? 📝":
			r.User{TgId: user.TgId, Username: "waiting"}.UpdateUser()
			msg.Text = "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️"
			msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
			msgs = append(msgs, msg)
		case fmt.Sprintf("%s Изменить аватар? %s", user.Avatar, user.Avatar):
			msg.Text = "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар..."
			msg.ReplyMarkup = s.EmojiInlineKeyboard()
			msgs = append(msgs, msg)
		case "/menu", v.GetString("user_location.menu"):
			msg.Text = "Меню"
			msg.ReplyMarkup = s.MainKeyboard(user)
			msgs = append(msgs, msg)
			r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
		default:
			msg.Text = user.GetUserInfo()
			msg.ReplyMarkup = s.ProfileKeyboard(user)
			msgs = append(msgs, msg)
		}
	}

	for i := range msgs {
		msgs[i].ChatID = update.Message.Chat.ID
	}

	return msgs
}

func userMapLocation(update tg.Update, user r.User) []tg.MessageConfig {
	newMessage := update.Message.Text
	char := strings.Fields(newMessage)

	fmt.Println(newMessage)

	if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Text == v.GetString("wordle.text_awaiting_msg") {
		msgs = s.UserSendNextWord(user, newMessage)
	} else if len(char) > 1 {
		msgs = useSpecialCell(char, user)
	} else if len(char) == 1 {
		msgs = useDefaultCell(update, user)
	} else {
		var msg tg.MessageConfig
		msg.ReplyToMessageID = update.Message.MessageID
		msg.Text = "🤨 Не пойму... 🧐"
		msgs = append(msgs, msg)
	}

	return msgs
}

func useSpecialCell(char []string, user r.User) []tg.MessageConfig {
	ItemLeftHand, ItemRightHand, ItemHead := s.UsersHandItemsView(user)
	var msg tg.MessageConfig

	// При нажатии кнопок
	switch char[0] {
	case v.GetString("message.doing.up"), v.GetString("message.doing.down"), v.GetString("message.doing.left"), v.GetString("message.doing.right"):
		msgs = append(msgs, s.UserMoving(user, char, char[0]))
	case v.GetString("message.emoji.foot"):
		msgs = append(msgs, s.UserMoving(user, char, char[1]))
	case v.GetString("message.emoji.hand"), ItemLeftHand.View, ItemRightHand.View:
		msgs = append(msgs, s.UserUseHandOrInstrumentMessage(user, char))
	case v.GetString("message.emoji.exclamation_mark"):
		cell := s.DirectionCell(user, char[3])
		msgs = append(msgs, s.ChoseInstrumentMessage(user, char, cell))
	case v.GetString("message.emoji.stop_use"):
		msgs = append(msgs, tg.MessageConfig{Text: v.GetString("errors.user_not_has_item_in_hand")})
	case "Рюкзак":
		msgs = append(msgs, s.BackpackCategoryKeyboard())
	case "Вещи":
		userItems := r.GetInventoryItems(user.ID)
		msg.Text = s.MessageGoodsUserItems(user, userItems, 0)
		msg.ReplyMarkup = s.GoodsInlineKeyboard(user, userItems, 0)
		msgs = append(msgs, msg)
	case v.GetString("message.emoji.online"):
		userOnline := true
		user = r.User{TgId: user.TgId, OnlineMap: &userOnline}.UpdateUser()
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%sОнлайн включен!", msg.Text, v.GetString("msg_separator"))
		msgs = append(msgs, msg)
	case v.GetString("message.emoji.offline"):
		userOnline := false
		user = r.User{TgId: user.TgId, OnlineMap: &userOnline}.UpdateUser()
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%sОнлайн выключен!", msg.Text, v.GetString("msg_separator"))
		msgs = append(msgs, msg)
	case ItemHead.View:
		res := s.DirectionCell(user, char[1])
		text, err := r.UpdateUserInstrument(user, ItemHead)
		if err != nil {
			msg.Text = fmt.Sprintf("%s%s%s", r.ViewItemInfo(res), v.GetString("msg_separator"), text)
		} else {
			msg.Text = r.ViewItemInfo(res)
		}
		msgs = append(msgs, msg)

		// ивент итемы
	case v.GetString("message.emoji.wrench"):
		loc := s.DirectionCell(user, char[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		charWorkbench := strings.Fields("workbench usPoint 0 1stComp null 0 2ndComp null 0 3rdComp null 0")
		msgs = append(msgs, s.Workbench(&cell, charWorkbench))
	case v.GetString("message.emoji.quest"):
		loc := s.DirectionCell(user, char[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		msgs = append(msgs, s.Quest(&cell, user))
	case v.GetString("message.emoji.wordle_game"):
		msgs = s.WordleMap(user)

		// Чатик
	case v.GetString("message.emoji.chat"):
		loc := s.DirectionCell(user, char[1])
		cell := r.Cell{MapsId: *loc.MapsId, AxisX: *loc.AxisX, AxisY: *loc.AxisY}.GetCell()
		msg.ReplyMarkup, msg.Text = s.OpenChatKeyboard(cell, user)
		msgs = append(msgs, msg)

	default:
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s\n\nНет инструмента в руке!", msg.Text)
		msgs = append(msgs, msg)
	}

	return msgs
}

func useDefaultCell(update tg.Update, user r.User) []tg.MessageConfig {
	var msg tg.MessageConfig
	newMessage := strings.Fields(update.Message.Text)
	currentTime := time.Now()

	// Взаимодействие с предметами на карте, у которых нет действий
	switch newMessage[0] {
	case v.GetString("message.doing.up"), v.GetString("message.doing.down"), v.GetString("message.doing.left"), v.GetString("message.doing.right"):
		msgs = append(msgs, s.UserMoving(user, newMessage, newMessage[0]))
	case v.GetString("message.emoji.water"):
		msg.Text = "Ты не похож на Jesus! 👮‍♂️"
		msgs = append(msgs, msg)
	case v.GetString("message.emoji.clock"):
		msg.Text = fmt.Sprintf("%s\nЧасики тикают...", currentTime.Format("15:04:05"))
		msgs = append(msgs, msg)
	case user.Avatar:
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s\n\n%s", user.GetUserInfo(), msg.Text)
		msgs = append(msgs, msg)
	case "/menu", v.GetString("user_location.menu"):
		msg.Text = "Меню"
		msg.ReplyMarkup = s.MainKeyboard(user)
		r.User{TgId: user.TgId, MenuLocation: "Меню"}.UpdateUser()
		msgs = append(msgs, msg)
	case v.GetString("message.emoji.casino"):
		msg.Text = "💰💵🤑 Ставки на JOY CASINO дот COM! 🤑💵💰"
		msgs = append(msgs, msg)
	case v.GetString("message.emoji.forbidden"):
		msg.Text = "🚫 Сюда нельзя! 🚫"
		msgs = append(msgs, msg)
	default:
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%sХммм....🤔", msg.Text, v.GetString("msg_separator"))
		msgs = append(msgs, msg)
	}

	for i := range msgs {
		msgs[i].ChatID = update.Message.Chat.ID
	}

	return msgs
}

func callBackResolver(update tg.Update) ([]tg.MessageConfig, bool) {
	msgs = []tg.MessageConfig{}
	var msg tg.MessageConfig
	buttons := tg.ReplyKeyboardMarkup{}
	charData := strings.Fields(update.CallbackQuery.Data)
	deletePrevMessage := true

	userTgId := r.GetUserTgId(update)
	user := r.GetUser(r.User{TgId: userTgId})
	ItemLeftHand, ItemRightHand, ItemHead := s.UsersHandItemsView(user)

	if len(charData) == 1 && charData[0] == v.GetString("callback_char.cancel") {
		msg.Text, msg.ReplyMarkup = r.GetMyMap(user)
		msgs = append(msgs, msg)
	}

	fmt.Println(charData)

	switch charData[0] {

	// Действия в рюкзаке
	case v.GetString("callback_char.category"):
		resUserItems := r.GetBackpackItems(user.ID, charData[1])
		msg.Text = s.MessageBackpackUserItems(resUserItems, 0, charData[1])
		msg.ReplyMarkup = s.BackpackInlineKeyboard(resUserItems, 0, charData[1])
		msgs = append(msgs, msg)
	case v.GetString("callback_char.backpack_moving"):
		msgs = append(msgs, s.BackPackMoving(charData, user))
	case v.GetString("callback_char.eat_food"):
		msgs = append(msgs, s.UserEatItem(user, charData))

	// Действия в инвентаре
	case v.GetString("callback_char.goods_moving"):
		msgs = append(msgs, s.GoodsMoving(charData, user))
	case v.GetString("callback_char.dress_good"):
		msgs = append(msgs, s.DressUserItem(user, charData))
	case v.GetString("callback_char.change_left_hand"), v.GetString("callback_char.change_right_hand"):
		user, userItem := r.UpdateUserHand(user, charData)
		charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		msg = s.GoodsMoving(charDataForOpenGoods, user)
		msg.Text = fmt.Sprintf("%s%sВы надели %s", msg.Text, v.GetString("msg_separator"), userItem.Item.View)
		msgs = append(msgs, msg)
	case v.GetString("callback_char.take_off_good"):
		msgs = append(msgs, s.UserTakeOffGood(user, charData))

	// Удаление, выкидывание, описание итема
	case v.GetString("callback_char.delete_item"):
		msgs = append(msgs, s.UserDeleteItem(user, charData))
	case v.GetString("callback_char.count_of_throw_out"):
		msgs = append(msgs, s.UserWantsToThrowOutItem(user, charData))
	case v.GetString("callback_char.throw_out_item"):
		msgs = append(msgs, s.UserThrowOutItem(user, charData))
	case v.GetString("callback_char.description"):
		msg.Text = r.UserItem{ID: r.ToInt(charData[1])}.GetFullDescriptionOfUserItem()
		msg.ReplyMarkup = s.DescriptionInlineButton(charData)
		msgs = append(msgs, msg)

	// Профиль
	case v.GetString("callback_char.change_avatar"):
		res := r.User{TgId: user.TgId, Avatar: charData[1]}.UpdateUser()
		msg.Text = res.GetUserInfo()
		msg.ReplyMarkup = s.ProfileKeyboard(res)
		msgs = append(msgs, msg)

	// Крафтинг
	case v.GetString("callback_char.workbench"):
		msgs = append(msgs, s.Workbench(nil, charData))
	case v.GetString("callback_char.receipt"):
		msg.Text = fmt.Sprintf("📖 *Рецепты*: 📖%s%s", v.GetString("msg_separator"), s.AllReceiptsMsg())
		msg.ReplyMarkup = nil
		deletePrevMessage = false
		msgs = append(msgs, msg)
	case v.GetString("callback_char.put_item"):
		userItem := r.GetUserItemsByType(user.ID, strings.Fields("food resource"))
		msg.ReplyMarkup = s.ChooseUserItemKeyboard(userItem, charData)
		msg.Text = fmt.Sprintf("%s%sВыбери предмет:", s.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
		msgs = append(msgs, msg)
	case v.GetString("callback_char.put_count_item"):
		msg = s.PutCountComponent(charData)
		msg.Text = fmt.Sprintf("%s%s⚠️ Сколько выкладываешь?", s.OpenWorkbenchMessage(charData), v.GetString("msg_separator"))
		msgs = append(msgs, msg)
	case v.GetString("callback_char.make_new_item"):
		resp := s.GetReceiptFromData(charData)
		receipt := r.FindReceiptForUser(resp)
		msg, deletePrevMessage = s.UserCraftItem(user, receipt)
		msgs = append(msgs, msg)

	// Использование надетых итемов
	case v.GetString("message.emoji.hand"), ItemLeftHand.View, ItemRightHand.View:
		msgs = append(msgs, s.UserUseHandOrInstrumentMessage(user, charData))
		res := s.DirectionCell(user, charData[1])
		msgs = append(msgs, s.ChoseInstrumentMessage(user, charData, res))
	case v.GetString("message.emoji.foot"):
		msgs = append(msgs, s.UserMoving(user, charData, charData[1]))
	case ItemHead.View:
		res := s.DirectionCell(user, charData[1])
		text, err := r.UpdateUserInstrument(user, ItemHead)
		msg.Text = r.ViewItemInfo(res)
		if err != nil {
			msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), text)
		}
		msgs = append(msgs, msg)

	// Квесты
	case v.GetString("callback_char.quests"):
		msg.Text = v.GetString("user_location.tasks_menu_message")
		msg.ReplyMarkup = s.AllQuestsMessageKeyboard(user)
		msgs = append(msgs, msg)
	case v.GetString("callback_char.quest"):
		msgs = append(msgs, s.OpenQuest(uint(r.ToInt(charData[1])), user))
	case v.GetString("callback_char.user_get_quest"):
		r.UserQuest{
			UserId:  user.ID,
			QuestId: uint(r.ToInt(charData[1])),
		}.GetOrCreateUserQuest()
		msgs = append(msgs, s.OpenQuest(uint(r.ToInt(charData[1])), user))
	case v.GetString("callback_char.user_done_quest"):
		msgs = append(msgs, s.UserDoneQuest(uint(r.ToInt(charData[1])), user))

	// Wordle
	case v.GetString("callback_char.wordle_regulations"):
		msg.Text = v.GetString("wordle.regulations")
		deletePrevMessage = false
		msgs = append(msgs, msg)
	case "awaitWord":
		msg.Text = v.GetString("wordle.text_awaiting_msg")
		msg.ReplyMarkup = tg.ForceReply{ForceReply: true}
		deletePrevMessage = false
		msgs = append(msgs, msg)
	case "wordleUserStatistic":
		msg.Text = r.GetWordleUserStatistic(user)
		deletePrevMessage = false
		msgs = append(msgs, msg)

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

		msg.Text, buttons = r.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), text)
		msg.ReplyMarkup = buttons
		msgs = append(msgs, msg)

	// Чатик
	case v.GetString("callback_char.join_to_chat"):
		ui := make([]r.ChatUser, 1)
		ui[0] = r.Chat{ID: uint(r.ToInt(charData[1]))}.GetOrCreateChatUser(user)
		cell := r.Cell{ID: uint(r.ToInt(charData[3]))}.GetCell()
		msg.ReplyMarkup, msg.Text = s.OpenChatKeyboard(cell, user)
		msgs = append(msgs, msg)
		s.NotifyUsers(ui, v.GetString("main_info.message_user_sign_in_chat"))
	}

	for i := range msgs {
		msgs[i].ParseMode = "markdown"
		msgs[i].ChatID = update.CallbackQuery.Message.Chat.ID
	}

	return msgs, deletePrevMessage
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
