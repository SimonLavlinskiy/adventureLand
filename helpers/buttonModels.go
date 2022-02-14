package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/repository"
	"strings"
)

func BackpackInlineKeyboard(items []repository.UserItem, i int) tgbotapi.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(items[i].Item.View+" "+
				repository.ToString(*items[i].Count)+"шт."+
				"   +"+repository.ToString(*items[i].Item.Healing)+" ♥️️"+
				"   +"+repository.ToString(*items[i].Item.Satiety)+"\U0001F9C3", "callbackAnswerAlert"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍽 1шт", "eatFood "+repository.ToString(items[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", "backpackMoving "+repository.ToString(i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑 все!", "throwOutFood "+repository.ToString(items[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔻", "backpackMoving "+repository.ToString(i+1)),
		),
	)
}

func GoodsInlineKeyboard(user repository.User, userItems []repository.UserItem, i int) tgbotapi.InlineKeyboardMarkup {
	if len(userItems) == 0 {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Пусто...(", "emptyGoods"),
			),
		)
	}

	text, data := repository.IsDressedItem(user, userItems[i])

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userItems[i].Item.View+" "+
				repository.ToString(*userItems[i].Count)+"шт.", "callbackAnswerAlert"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, data+" "+repository.ToString(userItems[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", "goodsMoving "+repository.ToString(i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑", "throwOutGood "+repository.ToString(userItems[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔻", "goodsMoving "+repository.ToString(i+1)),
		),
	)
}

func EmodjiInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	var listOfAvatar []string
	listOfAvatar = strings.Fields("🐶 🐱 🐭 🐹 🐰 🦊 🐻 🐼 ‍️🐨 🐯 🦁 🐮 🐷 🐸 🐵 🐦 🐧 🐔 🐤 🐥 🦆 🐴 🦄 🐺 🐗 🐝 🦋 🐛 🐌 🐞 🪲 🪰 🐜 🕷 🪳 🦖 🦕 🐙 🦀 🐟 🐠 🐡 🦭")

	for x := 0; x < len(listOfAvatar); x = x + 8 {
		var row []tgbotapi.InlineKeyboardButton
		for i := 0; i < 8; i++ {
			sum := x + i
			if len(listOfAvatar) > sum {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(listOfAvatar[sum], "changeAvatar "+listOfAvatar[sum]))
			}
		}
		buttons = append(buttons, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func ProfileKeyboard(user repository.User) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📝 Изменить имя? 📝"),
			tgbotapi.NewKeyboardButton(user.Avatar+" Изменить аватар? "+user.Avatar),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Меню"),
		),
	)
}

func MainKeyboard(user repository.User) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗺 Карта 🗺"),
			tgbotapi.NewKeyboardButton(user.Avatar+" Профиль 👔"),
		),
	)
}
