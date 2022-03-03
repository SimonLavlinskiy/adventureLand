package helpers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
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
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт.   +%d ♥️️   +%d\U0001F9C3", items[i].Item.View, *items[i].Count, *items[i].Item.Healing, *items[i].Item.Satiety),
				fmt.Sprintf("%s %d", v.GetString("callback_char.description"), items[i].ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍽 1шт", fmt.Sprintf("%s %d %d", v.GetString("callback_char.eat_food"), items[i].ID, i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d", v.GetString("callback_char.backpack_moving"), i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d backpack", v.GetString("callback_char.throw_out_item"), items[i].ID, i)),
			tgbotapi.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d", v.GetString("callback_char.backpack_moving"), i+1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %d %d backpack", v.GetString("callback_char.delete_item"), items[i].ID, i)),
		),
	)
}

func ChangeItemInHand(user repository.User, itemId int, charData2 string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("❓ %s ❔", user.LeftHand.View),
				fmt.Sprintf("%s %d %s", v.GetString("callback_char.change_left_hand"), itemId, charData2),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("❓ %s ❔", user.RightHand.View),
				fmt.Sprintf("%s %d %s", v.GetString("callback_char.change_right_hand"), itemId, charData2),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", fmt.Sprintf("%s %s", v.GetString("callback_char.goods_moving"), charData2)),
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
	itemDescription := "Описания нет("
	if userItems[i].Item.Description != nil {
		itemDescription = *userItems[i].Item.Description
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт. %s  %s", userItems[i].Item.View, *userItems[i].Count, userItems[i].Item.Name, itemDescription),
				fmt.Sprintf("%s %d", v.GetString("callback_char.description"), userItems[i].ID),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, fmt.Sprintf("%s %d %d", data, userItems[i].ID, i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d", v.GetString("callback_char.goods_moving"), i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d good", v.GetString("callback_char.throw_out_item"), userItems[i].ID, i)),
			tgbotapi.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d", v.GetString("callback_char.goods_moving"), i+1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %d %d good", v.GetString("callback_char.delete_item"), userItems[i].ID, i)),
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
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", x+y),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.count_of_delete"), buttonData[1], buttonData[2], x+y, buttonData[3])),
			)
		}
		buttons = append(buttons, row)
	}

	for y := 20; y <= maxCountItem; y = y + 20 {
		var row []tgbotapi.InlineKeyboardButton
		if y < maxCountItem {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", y),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.count_of_delete"), buttonData[1], buttonData[2], y, buttonData[3])),
			)
		}
		x := y + 10
		if y < maxCountItem {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", x),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.count_of_delete"), buttonData[1], buttonData[2], x, buttonData[3])),
			)
		}
		buttons = append(buttons, row)
	}

	var row []tgbotapi.InlineKeyboardButton
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Все!",
		fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.count_of_delete"), buttonData[1], buttonData[2], maxCountItem, buttonData[3])),
	)
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
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(listOfAvatar[sum],
					fmt.Sprintf("%s %s", v.GetString("callback_char.change_avatar"), listOfAvatar[sum])),
				)
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
			tgbotapi.NewKeyboardButton(fmt.Sprintf("%s Изменить аватар? %s", user.Avatar, user.Avatar)),
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
			if cell.Item.Cost != nil && *cell.Item.Cost > 0 && instrument == v.GetString("callback_char.hand") {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s ( %d💰 )", instrument, *cell.Item.Cost),
					fmt.Sprintf("%s %s %s", instrument, char[3], char[4])),
				)
			} else {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(instrument,
					fmt.Sprintf("%s %s %s", instrument, char[3], char[4])),
				)
			}
		}

		return tgbotapi.NewInlineKeyboardMarkup(
			row,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отмена", v.GetString("callback_char.cancel")),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На карту?", v.GetString("callback_char.cancel")),
		),
	)
}

func WorkbenchButton(char []string) tgbotapi.InlineKeyboardMarkup {
	leftArrow := "⬅️"
	rightArrow := "➡️"
	userPointer := repository.ToInt(char[2])

	defaultData := fmt.Sprintf("usPoint %d 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", userPointer, char[4], char[5], char[7], char[8], char[10], char[11])
	rightArrowData := fmt.Sprintf("%s usPoint %d 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.workbench"), userPointer+1, char[4], char[5], char[7], char[8], char[10], char[11])
	leftArrowData := fmt.Sprintf("%s usPoint %d 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.workbench"), userPointer-1, char[4], char[5], char[7], char[8], char[10], char[11])
	putItemData := fmt.Sprintf("%s %s", v.GetString("callback_char.put_item"), defaultData)
	putCountItemData := fmt.Sprintf("%s %s", v.GetString("callback_char.put_count_item"), defaultData)

	makeNewItem := fmt.Sprintf("%s %s", v.GetString("callback_char.make_new_item"), defaultData)

	if userPointer == 0 {
		leftArrow = "✖️"
		leftArrowData = "nothing"
	} else if userPointer == 2 {
		rightArrow = "✖️"
		rightArrowData = "nothing"
	}

	putItem := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Вставить предмет!", putItemData))
	changeItem := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✏️ Изменить", putItemData),
		tgbotapi.NewInlineKeyboardButtonData("🔢 Кол-во?", putCountItemData))

	ButtonManageItem := putItem

	if (userPointer == 0 && char[4] != "nil") || (userPointer == 1 && char[7] != "nil") || (userPointer == 2 && char[10] != "nil") {
		ButtonManageItem = changeItem
	}

	//"workbench usPoint: 0 1stComp: nil 0 2ndComp: nil 0 3rdComp: nil 0"

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✨⚡️ Слепить! ⚡️✨", makeNewItem),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(leftArrow, leftArrowData),
			tgbotapi.NewInlineKeyboardButtonData(rightArrow, rightArrowData),
		),
		ButtonManageItem,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Рецепты 📚", v.GetString("callback_char.receipt")),
		),
	)
}

func ChooseUserItemButton(userItem []repository.UserItem, char []string) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton

	var itemData string

	for x := 0; x < len(userItem); x = x + 5 {

		var row []tgbotapi.InlineKeyboardButton

		for i := 0; i < 5; i++ {
			if i+x < len(userItem) {
				switch char[2] {
				case "0":
					itemData = fmt.Sprintf("%s usPoint %s 1stComp %d %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), char[2], userItem[x+i].ID, char[5], char[7], char[8], char[10], char[11])
				case "1":
					itemData = fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %d %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), char[2], char[4], char[5], userItem[x+i].ID, char[8], char[10], char[11])
				case "2":
					itemData = fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %d %s", v.GetString("callback_char.put_count_item"), char[2], char[4], char[5], char[7], char[8], userItem[x+i].ID, char[11])
				}
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(userItem[x+i].Item.View, itemData))
			}
		}
		buttons = append(buttons, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func ChangeCountUserItem(charData []string, item repository.UserItem) tgbotapi.InlineKeyboardMarkup {
	charDone := fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.workbench"), charData[2], charData[4], charData[5], charData[7], charData[8], charData[10], charData[11])
	itemCount := repository.ToInt(charData[repository.ToInt(charData[2])+(5+repository.ToInt(charData[2])*2)])
	maxCountItem := item.Count

	appData := strings.Fields(fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), charData[2], charData[4], charData[5], charData[7], charData[8], charData[10], charData[11]))
	subData := strings.Fields(fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), charData[2], charData[4], charData[5], charData[7], charData[8], charData[10], charData[11]))

	subCount, appCount := fmt.Sprintf("%d", itemCount), fmt.Sprintf("%d", itemCount)

	if itemCount > 0 {
		subCount = fmt.Sprintf("%d", itemCount-1)
	}
	if itemCount < *maxCountItem {
		appCount = fmt.Sprintf("%d", itemCount+1)
	}

	subData[repository.ToInt(charData[2])+(5+repository.ToInt(charData[2])*2)] = subCount
	appData[repository.ToInt(charData[2])+(5+repository.ToInt(charData[2])*2)] = appCount

	subButData := fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), subData[2], subData[4], subData[5], subData[7], subData[8], subData[10], subData[11])
	appButData := fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), appData[2], appData[4], appData[5], appData[7], appData[8], appData[10], appData[11])

	subtractButton := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s⃣%s", subCount, item.Item.View), subButData)
	appendButton := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s⃣%s", appCount, item.Item.View), appButData)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("✅ Готово: %d⃣%s", itemCount, item.Item.View), charDone),
		),
		tgbotapi.NewInlineKeyboardRow(
			subtractButton,
			appendButton,
		),
	)
}
