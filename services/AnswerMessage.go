package services

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	r "project0/repository"
	"strings"
)

func DirectionCell(user r.User, direction string) r.Location {
	res := r.GetOrCreateMyLocation(user)

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

func MessageBackpackUserItems(userItems []r.UserItem, rowUser int, itemType string) string {
	var userItemMsg = fmt.Sprintf("🎒 *Рюкзачок* ➡️ *%s* \n \n", v.GetString(fmt.Sprintf("user_location.item_categories.%s", itemType)))

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
		switch itemType {
		case "food":
			userItemMsg += fmt.Sprintf("%s   %d%s     *HP*:  _%d_ ♥️️     *ST*:  _%d_ \U0001F9C3 ️\n", firstCell, *item.Count, item.Item.View, *item.Item.Healing, *item.Item.Satiety)
		case "resource", "sprout", "furniture":
			userItemMsg += fmt.Sprintf("%s   %s %d шт. - _%s_\n", firstCell, item.Item.View, *item.Count, item.Item.Name)
		default:
			userItemMsg += fmt.Sprintf("%s   %s %d шт.\n", firstCell, item.Item.View, *item.Count)
		}
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

func BackPackMoving(charData []string, user r.User) tg.MessageConfig {
	var msg tg.MessageConfig

	category := charData[2]
	userItems := r.GetBackpackItems(user.ID, category)

	var i int
	switch charData[1] {
	case "-1":
		i = len(userItems) - 1
	case fmt.Sprintf("%d", len(userItems)):
		i = 0
	default:
		i = r.ToInt(charData[1])
	}

	msg.Text = MessageBackpackUserItems(userItems, i, category)
	msg.ReplyMarkup = BackpackInlineKeyboard(userItems, i, category)

	return msg
}

func GoodsMoving(charData []string, user r.User) tg.MessageConfig {
	var msg tg.MessageConfig

	userItems := r.GetInventoryItems(user.ID)

	var i int
	switch charData[1] {
	case "-1":
		i = len(userItems) - 1
	case fmt.Sprintf("%d", len(userItems)):
		i = 0
	default:
		i = r.ToInt(charData[1])
	}

	msg.Text = MessageGoodsUserItems(user, userItems, i)
	msg.ReplyMarkup = GoodsInlineKeyboard(user, userItems, i)

	return msg
}

func UserEatItem(user r.User, charData []string) tg.MessageConfig {
	var msg tg.MessageConfig

	userItemId := r.ToInt(charData[1])

	userItem := r.UserItem{ID: userItemId}.UserGetUserItem()

	res := userItem.EatItem(user)
	charDataForOpenBackPack := strings.Fields(fmt.Sprintf("%s %s food", v.GetString("callback_char.backpack_moving"), charData[2]))
	msg = BackPackMoving(charDataForOpenBackPack, user)
	msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), res)

	return msg
}

func UserDeleteItem(user r.User, charData []string) tg.MessageConfig {
	var msg tg.MessageConfig

	userItemId := r.ToInt(charData[1])
	userItem := r.UserItem{ID: userItemId}.UserGetUserItem()

	if charData[4] == "false" {
		msg.ReplyMarkup = DeleteItem(charData)
		msg.Text = fmt.Sprintf("Вы точно хотите уничтожить %s %s _(%d шт.)_?", userItem.Item.View, userItem.Item.Name, *userItem.Count)
		return msg
	}

	countAfterUserThrowOutItem := 0
	var updateUserItemStruct = r.UserItem{
		ID:    userItemId,
		Count: &countAfterUserThrowOutItem,
	}

	user.UpdateUserItem(updateUserItemStruct)

	var charDataForOpenList []string
	if charData[3] == "good" {
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		UserTakeOffGood(user, charData)
		user = r.GetUser(r.User{TgId: user.TgId})
		msg = GoodsMoving(charDataForOpenList, user)
	} else {
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), charData[2], charData[3]))
		msg = BackPackMoving(charDataForOpenList, user)
	}

	msg.Text = fmt.Sprintf("%s%s🗑 Вы уничтожили %s%dшт.", msg.Text, v.GetString("msg_separator"), userItem.Item.View, *userItem.Count)

	return msg
}

func UsersHandItemsView(user r.User) (r.Item, r.Item, r.Item) {
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

func Quest(cell *r.Cell, user r.User) tg.MessageConfig {
	var msg tg.MessageConfig

	if !cell.IsQuest() {
		msg.Text = v.GetString("error.no_quest_item")
		return msg
	}

	msg.Text = v.GetString("user_location.tasks_menu_message")
	msg.ReplyMarkup = AllQuestsMessageKeyboard(user)

	return msg
}

func UserTakeOffGood(user r.User, charData []string) tg.MessageConfig {
	var msg tg.MessageConfig

	userItemId := r.ToInt(charData[1])
	userItem := r.UserItem{ID: userItemId}.UserGetUserItem()

	if user.HeadId != nil && userItem.ItemId == *user.HeadId {
		r.SetNullUserField(user, "head_id")
	} else if user.LeftHandId != nil && userItem.ItemId == *user.LeftHandId {
		r.SetNullUserField(user, "left_hand_id")
	} else if user.RightHandId != nil && userItem.ItemId == *user.RightHandId {
		r.SetNullUserField(user, "right_hand_id")
	} else if user.BodyId != nil && userItem.ItemId == *user.BodyId {
		r.SetNullUserField(user, "body_id")
	} else if user.FootId != nil && userItem.ItemId == *user.FootId {
		r.SetNullUserField(user, "foot_id")
	} else if user.ShoesId != nil && userItem.ItemId == *user.ShoesId {
		r.SetNullUserField(user, "shoes_id")
	}

	user = r.GetUser(r.User{TgId: user.TgId})

	charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
	msg = GoodsMoving(charDataForOpenGoods, user)
	msg.Text = fmt.Sprintf("%s%sВещь снята!", msg.Text, v.GetString("msg_separator"))

	return msg
}

func DressUserItem(user r.User, charData []string) tg.MessageConfig {
	var msg tg.MessageConfig

	userItemId := r.ToInt(charData[1])
	userItem := r.UserItem{ID: userItemId}.UserGetUserItem()
	changeHandItem := false

	var result = fmt.Sprintf("Вы надели %s", userItem.Item.View)

	switch *userItem.Item.DressType {
	case "hand":
		if user.LeftHandId == nil {
			clothes := &r.Clothes{LeftHandId: &userItem.ItemId}
			user = r.User{TgId: user.TgId, Clothes: *clothes}.UpdateUser()
		} else if user.RightHandId == nil {
			clothes := &r.Clothes{RightHandId: &userItem.ItemId}
			user = r.User{TgId: user.TgId, Clothes: *clothes}.UpdateUser()
		} else {
			result = "У вас заняты все руки! Что хочешь снять?"
			changeHandItem = true
		}
	case "head":
		if user.HeadId == nil {
			clothes := &r.Clothes{HeadId: &userItem.ItemId}
			user = r.User{TgId: user.TgId, Clothes: *clothes}.UpdateUser()
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "body":
		if user.BodyId == nil {
			clothes := &r.Clothes{BodyId: &userItem.ItemId}
			user = r.User{TgId: user.TgId, Clothes: *clothes}.UpdateUser()
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "foot":
		if user.FootId == nil {
			clothes := &r.Clothes{FootId: &userItem.ItemId}
			user = r.User{TgId: user.TgId, Clothes: *clothes}.UpdateUser()
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "shoes":
		if user.ShoesId == nil {
			clothes := &r.Clothes{ShoesId: &userItem.ItemId}
			user = r.User{TgId: user.TgId, Clothes: *clothes}.UpdateUser()
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	}

	if changeHandItem {
		msg.ReplyMarkup = ChangeItemInHandKeyboard(user, userItemId, charData[2])
	} else {
		charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		msg = GoodsMoving(charDataForOpenGoods, user)
	}

	msg.Text = fmt.Sprintf("%s%s%s", msg.Text, v.GetString("msg_separator"), result)

	return msg
}

func UserThrowOutItem(user r.User, charData []string) tg.MessageConfig {
	var msg tg.MessageConfig

	cellType := "item"
	userItem := r.UserItem{ID: r.ToInt(charData[1])}.UserGetUserItem()

	*userItem.Count = *userItem.Count - r.ToInt(charData[3])

	var msgText string

	if charData[4] == "other" && userItem.Item.Type == "chat" {
		cellType = "chat"
	}

	err := r.UpdateCellUnderUser(user, userItem, r.ToInt(charData[3]), cellType)
	if err != nil {
		msgText = fmt.Sprintf("%s%s", v.GetString("msg_separator"), err)
	} else {
		msgText = fmt.Sprintf("%sВы сбросили %s %sшт. на карту!", v.GetString("msg_separator"), userItem.Item.View, charData[3])
		user.UpdateUserItem(r.UserItem{ID: userItem.ID, Count: userItem.Count})
	}

	var charDataForOpenList []string
	if charData[4] == "good" {
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		if *userItem.Count == 0 {
			UserTakeOffGood(user, charData)
		}
		msg = GoodsMoving(charDataForOpenList, user)
	} else {
		charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), charData[2], charData[4]))
		msg = BackPackMoving(charDataForOpenList, user)
	}

	msg.Text = fmt.Sprintf("%s%s", msg.Text, msgText)

	return msg
}

func Workbench(cell *r.Cell, char []string) tg.MessageConfig {
	var msg tg.MessageConfig

	var charData []string
	if cell != nil && !cell.IsWorkbench() {
		msg.Text = "Здесь нет верстака!"
		return msg
	}

	if cell != nil {
		charData = strings.Fields("workbench usPoint 0 1stComp nil 0 2ndComp nil 0 3rdComp nil 0")
	} else {
		charData = strings.Fields(fmt.Sprintf("workbench usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", char[2], char[4], char[5], char[7], char[8], char[10], char[11]))
	}

	msg.Text = OpenWorkbenchMessage(charData)
	msg.ReplyMarkup = WorkbenchKeyboard(charData)

	return msg
}

func OpenWorkbenchMessage(char []string) string {
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

	msgText := fmt.Sprintf(
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

	return msgText
}

func getViewEmojiForMsg(char []string, i int) string {
	count := i + 5 + i*2

	if char[count] == "0" {
		return "\U0001F7EB"
	}

	return fmt.Sprintf("%s⃣", char[count])
}

func viewComponent(id string) string {
	if id != "nil" {
		component := r.UserItem{ID: r.ToInt(id)}.UserGetUserItem()
		return component.Item.View
	}
	return "⚪"
}

func UserWantsToThrowOutItem(user r.User, charData []string) tg.MessageConfig {
	var msg tg.MessageConfig

	userItem := r.UserItem{ID: r.ToInt(charData[1])}.UserGetUserItem()

	if userItem.CountUseLeft != nil && *userItem.CountUseLeft != *userItem.Item.CountUse {
		*userItem.Count = *userItem.Count - 1
	}

	if *userItem.Count == 0 {
		var charDataForOpenList []string
		if charData[3] == "good" {
			charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
			if *userItem.CountUseLeft == *userItem.Item.CountUse {
				UserTakeOffGood(user, charData)
			}
			msg = GoodsMoving(charDataForOpenList, user)
		} else {
			charDataForOpenList = strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), charData[2], charData[3]))
			msg = BackPackMoving(charDataForOpenList, user)
		}
		msg.Text = fmt.Sprintf("%s%sНельзя выкинуть на карту предмет, который уже был использован!", msg.Text, v.GetString("msg_separator"))
	} else {
		msg.ReplyMarkup = CountItemUserWantsToThrowKeyboard(charData, userItem)
		msg.Text = fmt.Sprintf("%sСколько %s ты хочешь скинуть на карту?", v.GetString("msg_separator"), userItem.Item.View)
	}

	return msg
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
		head = "〰️"
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

func PutCountComponent(char []string) tg.MessageConfig {
	var msg tg.MessageConfig

	userItemId := char[r.ToInt(char[2])+(4+r.ToInt(char[2])*2)] // char[x + (4+x*2 )] = char[4]
	userItem := r.UserItem{ID: r.ToInt(userItemId)}.UserGetUserItem()

	msg.ReplyMarkup = ChangeCountUserItemKeyboard(char, userItem)

	return msg
}

func UserCraftItem(user r.User, receipt *r.Receipt) (tg.MessageConfig, bool) {
	var msg tg.MessageConfig

	deletePrevMessage := true
	if receipt == nil {
		msg.Text = "Такого рецепта не существует!"
		msg.ReplyMarkup = nil
		deletePrevMessage = false
		return msg, deletePrevMessage
	}

	msg.ReplyMarkup = nil
	resultItem := r.UserItem{UserId: int(user.ID), ItemId: receipt.ItemResultID}.UserGetUserItem()

	if resultItem.Item.MaxCountUserHas != nil && *receipt.ItemResultCount+*resultItem.Count > *resultItem.Item.MaxCountUserHas {
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

func UserMoving(user r.User, char []string, charDirection string) tg.MessageConfig {
	var msg tg.MessageConfig

	var text string
	res := DirectionCell(user, charDirection)

	var inlineButtons tg.InlineKeyboardMarkup

	locMsg, err := r.UpdateLocation(char, res, user)
	msgMap, buttons := r.GetMyMap(user)

	if err != nil {
		if err.Error() == "user has not home" {
			inlineButtons = BuyHomeKeyboard()
			msg.ReplyMarkup = inlineButtons
			text = locMsg
		} else {
			text = fmt.Sprintf("%s%s%s", msgMap, v.GetString("msg_separator"), locMsg)
		}
	} else {
		msg.ReplyMarkup = buttons

		lighterMsg, err := user.CheckUserHasLighter()
		if err != nil {
			text = fmt.Sprintf("%s%s", v.GetString("msg_separator"), lighterMsg)
		}
		text = fmt.Sprintf("%s%s", msgMap, text)
	}

	msg.Text = text
	return msg
}

func UserUseHandOrInstrumentMessage(user r.User, char []string) tg.MessageConfig {
	var msg tg.MessageConfig

	res := DirectionCell(user, char[1])
	resultOfGetItem := r.UserGetItem(user, res, char)
	resText, buttons := r.GetMyMap(user)
	msg.Text = fmt.Sprintf("%s%s%s", resText, v.GetString("msg_separator"), resultOfGetItem)
	msg.ReplyMarkup = buttons

	return msg
}

func OpenQuest(questId uint, user r.User) tg.MessageConfig {
	var msg tg.MessageConfig

	quest := r.Quest{ID: questId}.GetQuest()
	userQuest := r.UserQuest{UserId: user.ID, QuestId: questId}.GetUserQuest()

	msg.Text = quest.QuestInfo(userQuest)
	msg.ReplyMarkup = OpenQuestKeyboard(quest, userQuest)

	return msg
}

func UserDoneQuest(questId uint, user r.User) tg.MessageConfig {
	var msg tg.MessageConfig

	userQuest := r.UserQuest{UserId: user.ID, QuestId: questId}.GetUserQuest()
	if !userQuest.Quest.Task.HasUserDoneTask(user) {
		msg.Text = v.GetString("errors.user_did_not_task")
		return msg
	}

	userQuest.UserDoneQuest(user)
	user.UserGetResult(userQuest.Quest.Result)

	msgQuestResult := msgQuestResult(userQuest.Quest.Result)

	msg = OpenQuest(questId, user)
	msg.Text = fmt.Sprintf("*Задание выполнено!*\n%s%s%s", msg.Text, v.GetString("msg_separator"), msgQuestResult)

	return msg
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

func ChoseInstrumentMessage(user r.User, char []string, cellLocation r.Location) tg.MessageConfig {
	var msg tg.MessageConfig

	if char[0] != v.GetString("message.emoji.exclamation_mark") {
		char = strings.Fields(fmt.Sprintf("❗ 🛠 ❓ %s %s", char[1], char[2]))
	}
	cell := r.Cell{MapsId: *cellLocation.MapsId, AxisX: *cellLocation.AxisX, AxisY: *cellLocation.AxisY}
	cell = cell.GetCell()

	buttons, err := ChooseInstrumentKeyboard(char, cell, user)

	if err == nil {
		msg.ReplyMarkup = buttons
		msg.Text = v.GetString("errors.chose_instrument_to_use")
		msg.ChatID = int64(user.TgId)
	} else {
		msg.Text = "Тут ничего нет..."
	}

	return msg
}

func msgQuestResult(result r.Result) string {
	result = result.GetResult()

	msg := "🏆 *Ты получил*:"
	if result.Item != nil {
		msg = fmt.Sprintf("%s\n_%s %s - %d шт._", msg, result.Item.View, result.Item.Name, *result.CountItem)
	}
	if result.SpecialItem != nil {
		msg = fmt.Sprintf("%s\n_%s %s - %d шт._", msg, result.SpecialItem.View, result.SpecialItem.Name, *result.SpecialItemCount)
	}

	return msg
}
