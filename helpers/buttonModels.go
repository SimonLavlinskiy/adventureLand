package helpers

import (
	"fmt"
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
				"   +"+repository.ToString(*items[i].Item.Satiety)+"\U0001F9C3", "description "+repository.ToString(items[i].ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍽 1шт", "eatFood "+repository.ToString(items[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", "backpackMoving "+repository.ToString(i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👋🗑🗺", "throwOutItem "+repository.ToString(items[i].ID)+" "+repository.ToString(i)+" backpack"),
			tgbotapi.NewInlineKeyboardButtonData("🔻", "backpackMoving "+repository.ToString(i+1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💥🗑💥", "deleteItem "+repository.ToString(items[i].ID)+" "+repository.ToString(i)+" backpack"),
		),
	)
}

func ChangeItemInHand(user repository.User, itemId int, charData2 string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ "+user.LeftHand.View+" ❔", "changeLeftHand "+repository.ToString(itemId)+" "+charData2),
			tgbotapi.NewInlineKeyboardButtonData("❔ "+user.RightHand.View+" ❓", "changeRightHand "+repository.ToString(itemId)+" "+charData2),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "goodsMoving "+charData2),
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
	itemDesciption := "Описания нет("
	if userItems[i].Item.Description != nil {
		itemDesciption = *userItems[i].Item.Description
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userItems[i].Item.View+" "+repository.ToString(*userItems[i].Count)+"шт. "+userItems[i].Item.Name+"  "+itemDesciption,
				"description "+repository.ToString(userItems[i].ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, data+" "+repository.ToString(userItems[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", "goodsMoving "+repository.ToString(i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👋🗑🗺", "throwOutItem "+repository.ToString(userItems[i].ID)+" "+repository.ToString(i)+" good"),
			tgbotapi.NewInlineKeyboardButtonData("🔻", "goodsMoving "+repository.ToString(i+1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💥🗑💥", "deleteItem "+repository.ToString(userItems[i].ID)+" "+repository.ToString(i)+" good"),
		),
	)
}

func CountItemUserWantsToThrow(buttonData []string, userItem repository.UserItem) tgbotapi.InlineKeyboardMarkup {
	maxCountItem := *userItem.Count
	var buttons [][]tgbotapi.InlineKeyboardButton

	for x := 1; x < 10; x = x + 5 {
		var row []tgbotapi.InlineKeyboardButton
		if x > maxCountItem {
			break
		}
		for y := 0; y < 5; y++ {
			if x+y > maxCountItem {
				break
			}
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(repository.ToString(x+y)+"шт.", fmt.Sprintf("countOfDelete %s %s %s %s", buttonData[1], buttonData[2], repository.ToString(x+y), buttonData[3])))
		}
		buttons = append(buttons, row)
	}

	for y := 20; y <= maxCountItem; y = y + 20 {
		var row []tgbotapi.InlineKeyboardButton
		if y < maxCountItem {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(repository.ToString(y)+"шт.", fmt.Sprintf("countOfDelete %s %s %s %s", buttonData[1], buttonData[2], repository.ToString(y), buttonData[3])))
		}
		x := y + 10
		if y < maxCountItem {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(repository.ToString(x)+"шт.", fmt.Sprintf("countOfDelete %s %s %s %s", buttonData[1], buttonData[2], repository.ToString(x), buttonData[3])))
		}
		buttons = append(buttons, row)
	}

	var row []tgbotapi.InlineKeyboardButton
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Все!", fmt.Sprintf("countOfDelete %s %s %s %s", buttonData[1], buttonData[2], repository.ToString(maxCountItem), buttonData[3])))
	buttons = append(buttons, row)

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
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

func ChooseInstrument(char []string, cell repository.Cellule, user repository.User) tgbotapi.InlineKeyboardMarkup {
	instruments := repository.GetInstrumentsUserCanUse(user, cell)

	if len(instruments) != 0 {
		var row []tgbotapi.InlineKeyboardButton

		for _, instrument := range instruments {
			if cell.Item.Cost != nil && *cell.Item.Cost > 0 && instrument == "👋" {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(instrument+" ( "+repository.ToString(*cell.Item.Cost)+"💰 )", instrument+" "+char[3]+" "+char[4]))
			} else {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(instrument, instrument+" "+char[3]+" "+char[4]))
			}
		}

		return tgbotapi.NewInlineKeyboardMarkup(
			row,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На карту?", "cancel"),
		),
	)
}
