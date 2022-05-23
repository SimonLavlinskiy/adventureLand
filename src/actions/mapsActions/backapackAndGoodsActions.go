package mapsActions

import (
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/backpackAndGoodsController"
	"project0/src/controllers/backpackAndGoodsController/backpackController"
	"project0/src/controllers/backpackAndGoodsController/goodsController"
	"project0/src/models"
	"project0/src/services/helpers"
)

func CheckBackpackAndGoodsAction(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup, err error) {
	switch charData[0] {
	// Действия в рюкзаке
	case v.GetString("callback_char.category"):
		msg, buttons = backpackController.ListOfBackpackItems(user, charData)
	case v.GetString("callback_char.backpack_moving"):
		msg, buttons = backpackController.BackPackMoving(charData, user)
	case v.GetString("callback_char.eat_food"):
		msg, buttons = backpackController.UserEatItem(user, charData)

	// Действия в инвентаре
	case v.GetString("callback_char.goods_moving"):
		msg, buttons = goodsController.ListOfGoods(user, charData)
	case v.GetString("callback_char.dress_good"):
		msg, buttons = goodsController.DressUserItem(user, charData)
	case v.GetString("callback_char.change_left_hand"), v.GetString("callback_char.change_right_hand"):
		msg, buttons = goodsController.ChangeHand(user, charData)
	case v.GetString("callback_char.take_off_good"):
		msg, buttons = goodsController.UserTakeOffGood(user, charData)

	// Удаление, выкидывание, описание итема
	case v.GetString("callback_char.delete_item"):
		msg, buttons = backpackAndGoodsController.UserDeleteItem(user, charData)
	case v.GetString("callback_char.count_of_throw_out"):
		msg, buttons = backpackAndGoodsController.UserWantsToThrowOutItem(user, charData)
	case v.GetString("callback_char.throw_out_item"):
		msg, buttons = backpackAndGoodsController.UserThrowOutItem(user, charData)
	case v.GetString("callback_char.description"):
		msg = models.UserItem{ID: helpers.ToInt(charData[1])}.GetFullDescriptionOfUserItem()
		buttons = backpackAndGoodsController.DescriptionInlineButton(charData)
	default:
		err = errors.New("not good or backpack actions")
	}

	return msg, buttons, err
}
