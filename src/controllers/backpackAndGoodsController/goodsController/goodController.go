package goodsController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/userController"
	"project0/src/controllers/userItemController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"strings"
)

func ChangeHand(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	user, userItem := userController.UpdateUserHand(user, charData)
	charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
	msg, buttons = GoodsMoving(charDataForOpenGoods, user)
	msg = fmt.Sprintf("%s%sВы надели %s", msg, v.GetString("msg_separator"), userItem.Item.View)
	return msg, buttons
}

func ListOfGoods(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	if len(charData) == 1 {
		userItems := userItemController.GetInventoryItems(user.ID)
		msg = MessageGoodsUserItems(user, userItems, 0)
		buttons = GoodsInlineKeyboard(user, userItems, 0)
	} else {
		msg, buttons = GoodsMoving(charData, user)
	}
	return msg, buttons
}

func MessageGoodsUserItems(user models.User, userItems []models.UserItem, rowUser int) string {
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
		userItemMsg += fmt.Sprintf("%s  %s %d шт.  %s %s... (%d/%d)\n", firstCell, item.Item.View, *item.Count, res, strings.Split(item.Item.Name, " ")[0], *item.CountUseLeft, *item.Item.CountUse)

	}

	return userItemMsg
}

func GoodsMoving(charData []string, user models.User) (msgText string, buttons tg.InlineKeyboardMarkup) {
	userItems := userItemController.GetInventoryItems(user.ID)

	var i int
	switch charData[1] {
	case "-1":
		i = len(userItems) - 1
	case fmt.Sprintf("%d", len(userItems)):
		i = 0
	default:
		i = helpers.ToInt(charData[1])
	}

	msgText = MessageGoodsUserItems(user, userItems, i)
	buttons = GoodsInlineKeyboard(user, userItems, i)

	return msgText, buttons
}

func messageUserDressedGoods(user models.User) string {
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
		body = "👕"
	}
	if user.Foot != nil {
		foot = user.Foot.View
	} else {
		foot = "\U0001FA73"
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

func UserTakeOffClothes(user models.User, charData []string) (msgText string, buttons tg.InlineKeyboardMarkup) {
	userItemId := helpers.ToInt(charData[1])
	userItem := models.UserItem{ID: userItemId}.GetOrCreateUserItem()

	if user.HeadId != nil && userItem.ItemId == *user.HeadId {
		repositories.SetNullUserField(user, "head_id")
	} else if user.LeftHandId != nil && userItem.ItemId == *user.LeftHandId {
		repositories.SetNullUserField(user, "left_hand_id")
	} else if user.RightHandId != nil && userItem.ItemId == *user.RightHandId {
		repositories.SetNullUserField(user, "right_hand_id")
	} else if user.BodyId != nil && userItem.ItemId == *user.BodyId {
		repositories.SetNullUserField(user, "body_id")
	} else if user.FootId != nil && userItem.ItemId == *user.FootId {
		repositories.SetNullUserField(user, "foot_id")
	} else if user.ShoesId != nil && userItem.ItemId == *user.ShoesId {
		repositories.SetNullUserField(user, "shoes_id")
	}

	user = repositories.GetUser(models.User{TgId: user.TgId})

	charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
	msgText, buttons = GoodsMoving(charDataForOpenGoods, user)
	msgText = fmt.Sprintf("%s%sВещь снята!", msgText, v.GetString("msg_separator"))

	return msgText, buttons
}

func DressUserItem(user models.User, charData []string) (msgText string, buttons tg.InlineKeyboardMarkup) {

	userItemId := helpers.ToInt(charData[1])
	userItem := models.UserItem{ID: userItemId}.GetOrCreateUserItem()
	changeHandItem := false

	var result = fmt.Sprintf("Вы надели %s", userItem.Item.View)

	switch *userItem.Item.DressType {
	case "hand":
		if user.LeftHandId == nil {
			clothes := &models.Clothes{LeftHandId: &userItem.ItemId}
			user = repositories.UpdateUser(models.User{TgId: user.TgId, Clothes: *clothes})
		} else if user.RightHandId == nil {
			clothes := &models.Clothes{RightHandId: &userItem.ItemId}
			user = repositories.UpdateUser(models.User{TgId: user.TgId, Clothes: *clothes})
		} else {
			result = "У вас заняты все руки! Что хочешь снять?"
			changeHandItem = true
		}
	case "head":
		if user.HeadId == nil {
			clothes := &models.Clothes{HeadId: &userItem.ItemId}
			user = repositories.UpdateUser(models.User{TgId: user.TgId, Clothes: *clothes})
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "body":
		if user.BodyId == nil {
			clothes := &models.Clothes{BodyId: &userItem.ItemId}
			user = repositories.UpdateUser(models.User{TgId: user.TgId, Clothes: *clothes})
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "foot":
		if user.FootId == nil {
			clothes := &models.Clothes{FootId: &userItem.ItemId}
			user = repositories.UpdateUser(models.User{TgId: user.TgId, Clothes: *clothes})
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	case "shoes":
		if user.ShoesId == nil {
			clothes := &models.Clothes{ShoesId: &userItem.ItemId}
			user = repositories.UpdateUser(models.User{TgId: user.TgId, Clothes: *clothes})
		} else {
			result = v.GetString("errors.user_has_other_item")
		}
	}

	if changeHandItem {
		buttons = ChangeItemInHandKeyboard(user, userItemId, charData[2])
	} else {
		charDataForOpenGoods := strings.Fields(fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData[2]))
		msgText, buttons = GoodsMoving(charDataForOpenGoods, user)
	}

	msgText = fmt.Sprintf("%s%s%s", msgText, v.GetString("msg_separator"), result)

	return msgText, buttons
}

func ChangeItemInHandKeyboard(user models.User, itemId int, charData2 string) tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("❓ %s ❔", user.LeftHand.View),
				fmt.Sprintf("%s %d %s", v.GetString("callback_char.change_left_hand"), itemId, charData2),
			),
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("❓ %s ❔", user.RightHand.View),
				fmt.Sprintf("%s %d %s", v.GetString("callback_char.change_right_hand"), itemId, charData2),
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Отмена", fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData2)),
		),
	)
}
