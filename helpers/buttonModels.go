package helpers

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	r "project0/repository"
	str "strings"
)

func BackpackInlineKeyboard(items []r.UserItem, i int, backpackType string) tg.InlineKeyboardMarkup {
	switch backpackType {
	case "food":
		return FoodListBackpackInlineKeyboard(items, i)
	case "sprout":
		return SproutListBackpackInlineKeyboard(items, i, backpackType)
	default:
		return DefaultListBackpackInlineKeyboard(items, i, backpackType)
	}
}

func FoodListBackpackInlineKeyboard(items []r.UserItem, i int) tg.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт.   +%d ♥️️   +%d\U0001F9C3", items[i].Item.View, *items[i].Count, *items[i].Item.Healing, *items[i].Item.Satiety),
				fmt.Sprintf("%s %d %d food", v.GetString("callback_char.description"), items[i].ID, i)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("🍽 1шт", fmt.Sprintf("%s %d %d", v.GetString("callback_char.eat_food"), items[i].ID, i)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d food", v.GetString("callback_char.backpack_moving"), i-1)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d food", v.GetString("callback_char.throw_out_item"), items[i].ID, i)),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d food", v.GetString("callback_char.backpack_moving"), i+1)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %d %d food", v.GetString("callback_char.delete_item"), items[i].ID, i)),
		),
	)
}

func DescriptionInlineButton(char []string) tg.InlineKeyboardMarkup {
	switch char[3] {
	case "food":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("🍽 1шт", fmt.Sprintf("%s %s %s", v.GetString("callback_char.eat_food"), char[1], char[2])),
				tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %s %s food", v.GetString("callback_char.throw_out_item"), char[1], char[2])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %s %s food", v.GetString("callback_char.delete_item"), char[1], char[2])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s food", v.GetString("callback_char.backpack_moving"), char[2])),
			),
		)
	case "resource":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %s %s %s", v.GetString("callback_char.throw_out_item"), char[1], char[2], char[3])),
				tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %s %s %s", v.GetString("callback_char.delete_item"), char[1], char[2], char[3])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), char[2], char[3])),
			),
		)
	case "sprout":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("👋\U0001F9A0🗺", fmt.Sprintf("%s %s %s %s", v.GetString("callback_char.throw_out_item"), char[1], char[2], char[3])),
				tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %s %s %s", v.GetString("callback_char.delete_item"), char[1], char[2], char[3])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), char[2], char[3])),
			),
		)
	case "furniture":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("👋\U0001F9A0🗺", fmt.Sprintf("%s %s %s %s", v.GetString("callback_char.throw_out_item"), char[1], char[2], char[3])),
				tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %s %s %s", v.GetString("callback_char.delete_item"), char[1], char[2], char[3])),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s %s", v.GetString("callback_char.backpack_moving"), char[2], char[3])),
			),
		)
	case "good":
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("%s %s good", v.GetString("callback_char.goods_moving"), char[2])),
			),
		)
	default:
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
			),
		)
	}
}

func DefaultListBackpackInlineKeyboard(items []r.UserItem, i int, itemType string) tg.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт.", items[i].Item.View, *items[i].Count),
				fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.description"), items[i].ID, i, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.throw_out_item"), items[i].ID, i, items[i].Item.Type)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i-1, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.delete_item"), items[i].ID, i, itemType)),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i+1, itemType)),
		),
	)
}

func SproutListBackpackInlineKeyboard(items []r.UserItem, i int, itemType string) tg.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт. - %s", items[i].Item.View, *items[i].Count, *items[i].Item.Description),
				fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.description"), items[i].ID, i, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋\U0001F9A0🗺", fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.throw_out_item"), items[i].ID, i, itemType)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i-1, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.delete_item"), items[i].ID, i, itemType)),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i+1, itemType)),
		),
	)
}

func ChangeItemInHandKeyboard(user r.User, itemId int, charData2 string) tg.InlineKeyboardMarkup {
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

func GoodsInlineKeyboard(user r.User, userItems []r.UserItem, i int) tg.InlineKeyboardMarkup {
	if len(userItems) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyGoods"),
			),
		)
	}

	text, data := user.IsDressedItem(userItems[i])
	itemDescription := "Описания нет("
	if userItems[i].Item.Description != nil {
		itemDescription = *userItems[i].Item.Description
	}

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт. %s  %s", userItems[i].Item.View, *userItems[i].Count, userItems[i].Item.Name, itemDescription),
				fmt.Sprintf("%s %d %d good", v.GetString("callback_char.description"), userItems[i].ID, i),
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(text, fmt.Sprintf("%s %d %d", data, userItems[i].ID, i)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d", v.GetString("callback_char.goods_moving"), i-1)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d good", v.GetString("callback_char.throw_out_item"), userItems[i].ID, i)),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d", v.GetString("callback_char.goods_moving"), i+1)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %d %d good", v.GetString("callback_char.delete_item"), userItems[i].ID, i)),
		),
	)
}

func CountItemUserWantsToThrowKeyboard(buttonData []string, userItem r.UserItem) tg.InlineKeyboardMarkup {
	maxCountItem := *userItem.Count
	var buttons [][]tg.InlineKeyboardButton

	for x := 1; x < 10; x = x + 5 {
		var row []tg.InlineKeyboardButton
		if x > maxCountItem {
			break
		}
		for y := 0; y < 5; y++ {
			if x+y > maxCountItem {
				break
			}
			row = append(row, tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", x+y),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.count_of_throw_out"), buttonData[1], buttonData[2], x+y, buttonData[3])),
			)
		}
		buttons = append(buttons, row)
	}

	for y := 20; y <= maxCountItem; y = y + 20 {
		var row []tg.InlineKeyboardButton
		if y < maxCountItem {
			row = append(row, tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", y),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.count_of_throw_out"), buttonData[1], buttonData[2], y, buttonData[3])),
			)
		}
		x := y + 10
		if y < maxCountItem {
			row = append(row, tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d шт.", x),
				fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.count_of_throw_out"), buttonData[1], buttonData[2], x, buttonData[3])),
			)
		}
		buttons = append(buttons, row)
	}

	var row []tg.InlineKeyboardButton
	row = append(row, tg.NewInlineKeyboardButtonData("Все!",
		fmt.Sprintf("%s %s %s %d %s", v.GetString("callback_char.count_of_throw_out"), buttonData[1], buttonData[2], maxCountItem, buttonData[3])),
	)
	buttons = append(buttons, row)

	return tg.NewInlineKeyboardMarkup(buttons...)
}

func EmojiInlineKeyboard() tg.InlineKeyboardMarkup {
	var buttons [][]tg.InlineKeyboardButton
	var listOfAvatar []string
	listOfAvatar = str.Fields(v.GetString("message.list_of_avatar"))

	for x := 0; x < len(listOfAvatar); x = x + 8 {
		var row []tg.InlineKeyboardButton
		for i := 0; i < 8; i++ {
			sum := x + i
			if len(listOfAvatar) > sum {
				row = append(row, tg.NewInlineKeyboardButtonData(listOfAvatar[sum],
					fmt.Sprintf("%s %s", v.GetString("callback_char.change_avatar"), listOfAvatar[sum])),
				)
			}
		}
		buttons = append(buttons, row)
	}

	return tg.NewInlineKeyboardMarkup(buttons...)
}

func ProfileKeyboard(user r.User) tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("📝 Изменить имя? 📝"),
			tg.NewKeyboardButton(fmt.Sprintf("%s Изменить аватар? %s", user.Avatar, user.Avatar)),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Меню"),
		),
	)
}

func MainKeyboard(user r.User) tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("🗺 Карта 🗺"),
			tg.NewKeyboardButton(user.Avatar+" Профиль 👔"),
		),
	)
}

func ChooseInstrumentKeyboard(char []string, cell r.Cell, user r.User) tg.InlineKeyboardMarkup {
	instruments := r.GetInstrumentsUserCanUse(user, cell)

	if len(instruments) != 0 {
		var row []tg.InlineKeyboardButton

		for instrument, i := range instruments {
			if cell.Item.Cost != nil && *cell.Item.Cost > 0 && (i == "hand" || i == "swap") && cell.NeedPay {
				row = append(row, tg.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s ( %d💰 )", instrument, *cell.Item.Cost),
					fmt.Sprintf("%s %s %s", instrument, char[3], char[4])),
				)
			} else {
				row = append(row, tg.NewInlineKeyboardButtonData(
					instrument,
					fmt.Sprintf("%s %s %s", instrument, char[3], char[4])),
				)
			}
		}

		return tg.NewInlineKeyboardMarkup(
			row,
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Отмена", v.GetString("callback_char.cancel")),
			),
		)
	}

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("На карту?", v.GetString("callback_char.cancel")),
		),
	)
}

func WorkbenchKeyboard(char []string) tg.InlineKeyboardMarkup {
	leftArrow := "⬅️"
	rightArrow := "➡️"
	userPointer := r.ToInt(char[2])

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

	putItem := tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData("Вставить предмет!", putItemData))
	changeItem := tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData("✏️ Изменить", putItemData),
		tg.NewInlineKeyboardButtonData("🔢 Кол-во?", putCountItemData))

	ButtonManageItem := putItem

	if (userPointer == 0 && char[4] != "nil") || (userPointer == 1 && char[7] != "nil") || (userPointer == 2 && char[10] != "nil") {
		ButtonManageItem = changeItem
	}

	//"workbench usPoint: 0 1stComp: nil 0 2ndComp: nil 0 3rdComp: nil 0"

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("✨⚡️ Слепить! ⚡️✨", makeNewItem),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(leftArrow, leftArrowData),
			tg.NewInlineKeyboardButtonData(rightArrow, rightArrowData),
		),
		ButtonManageItem,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Рецепты 📚", v.GetString("callback_char.receipt")),
		),
	)
}

func ChooseUserItemKeyboard(userItem []r.UserItem, char []string) tg.InlineKeyboardMarkup {
	var buttons [][]tg.InlineKeyboardButton

	var itemData string

	for x := 0; x < len(userItem); x = x + 5 {

		var row []tg.InlineKeyboardButton

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
				row = append(row, tg.NewInlineKeyboardButtonData(userItem[x+i].Item.View, itemData))
			}
		}
		buttons = append(buttons, row)
	}

	return tg.NewInlineKeyboardMarkup(buttons...)
}

func ChangeCountUserItemKeyboard(charData []string, item r.UserItem) tg.InlineKeyboardMarkup {
	charDone := fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.workbench"), charData[2], charData[4], charData[5], charData[7], charData[8], charData[10], charData[11])
	itemCount := r.ToInt(charData[r.ToInt(charData[2])+(5+r.ToInt(charData[2])*2)])
	maxCountItem := item.Count

	appData := str.Fields(fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), charData[2], charData[4], charData[5], charData[7], charData[8], charData[10], charData[11]))
	subData := str.Fields(fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), charData[2], charData[4], charData[5], charData[7], charData[8], charData[10], charData[11]))

	subCount, appCount := fmt.Sprintf("%d", itemCount), fmt.Sprintf("%d", itemCount)

	if itemCount > 0 {
		subCount = fmt.Sprintf("%d", itemCount-1)
	}
	if itemCount < *maxCountItem {
		appCount = fmt.Sprintf("%d", itemCount+1)
	}

	subData[r.ToInt(charData[2])+(5+r.ToInt(charData[2])*2)] = subCount
	appData[r.ToInt(charData[2])+(5+r.ToInt(charData[2])*2)] = appCount

	subButData := fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), subData[2], subData[4], subData[5], subData[7], subData[8], subData[10], subData[11])
	appButData := fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.put_count_item"), appData[2], appData[4], appData[5], appData[7], appData[8], appData[10], appData[11])

	subtractButton := tg.NewInlineKeyboardButtonData(fmt.Sprintf("%s⃣%s", subCount, item.Item.View), subButData)
	appendButton := tg.NewInlineKeyboardButtonData(fmt.Sprintf("%s⃣%s", appCount, item.Item.View), appButData)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(fmt.Sprintf("✅ Готово: %d⃣%s", itemCount, item.Item.View), charDone),
		),
		tg.NewInlineKeyboardRow(
			subtractButton,
			appendButton,
		),
	)
}

func AllQuestsMessageKeyboard(u r.User) tg.InlineKeyboardMarkup {
	quests := r.Quest{}.GetQuests()
	userQuests := r.User{ID: u.ID}.GetUserQuests()

	type statusQuest struct {
		status string
		quest  r.Quest
	}

	m := map[uint]statusQuest{}
	for _, quest := range quests {
		m[quest.ID] = statusQuest{status: "new", quest: quest}
	}

	for _, uq := range userQuests {
		m[uq.QuestId] = statusQuest{status: uq.Status, quest: uq.Quest}
	}

	if len(quests) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "cancel"),
			),
		)
	}

	var result [][]tg.InlineKeyboardButton

	for _, i := range m {
		status := v.GetString(fmt.Sprintf("quest_statuses.%s_emoji", i.status))
		result = append(result,
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s - Задание: «%s»", status, i.quest.Name),
					fmt.Sprintf("quest %d", i.quest.ID),
				),
			),
		)
	}

	return tg.NewInlineKeyboardMarkup(result...)
}

func OpenQuestKeyboard(q r.Quest, uq r.UserQuest) tg.InlineKeyboardMarkup {
	if uq.Status == "" {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Взять в работу", fmt.Sprintf("user_get_quest %d", q.ID)),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", "quests"),
			),
		)
	}

	switch uq.Status {
	case "processed":
		var buttonStatus tg.InlineKeyboardButton
		if q.Task.HasUserDoneTask(uq.User) {
			buttonStatus = tg.NewInlineKeyboardButtonData("Готово! Я всё сделаль!", fmt.Sprintf("user_done_quest %d", uq.QuestId))
		} else {
			buttonStatus = tg.NewInlineKeyboardButtonData("Еще в работе... Прийду потом", "quests")
		}

		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				buttonStatus,
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Назад", "quests"),
			),
		)
	}

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Назад", "quests"),
		),
	)
}

func BuyHomeKeyboard() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(fmt.Sprintf("🏘 Купить дом! 🏘 (%d 💰)", v.GetInt("main_info.cost_of_house")), "buyHome"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Отмена", "cancel"),
		),
	)
}

func BackpackCategoryKeyboard() (tg.InlineKeyboardMarkup, string) {
	categories := str.Fields(v.GetString("user_location.item_categories.categories"))

	var rows [][]tg.InlineKeyboardButton

	for _, category := range categories {
		rows = append(rows, tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s", v.GetString(fmt.Sprintf("user_location.item_categories.%s", category))),
			fmt.Sprintf("%s", fmt.Sprintf("category %s", category)),
		)))
	}

	return tg.NewInlineKeyboardMarkup(
		rows...,
	), fmt.Sprintf("🎒 *Рюкзачок*\n%s", v.GetString("user_location.item_categories.category_title"))
}

func OpenChatKeyboard(cell r.Cell, user r.User) (tg.InlineKeyboardMarkup, string) {
	var button tg.InlineKeyboardButton
	msgText := "Присоединяйся и общайтесь!"

	if !cell.IsChat() {
		msgText = "Здесь нет чата! Поищи в другом месте..."
		button = tg.NewInlineKeyboardButtonData("Жаль...", "cancel")
	} else {
		userChat := cell.Chat.GetChatUser(user)

		if userChat == nil {
			button = tg.NewInlineKeyboardButtonData("Присоединиться к беседе", fmt.Sprintf("joinToChat %d cell %d", *cell.ChatId, cell.ID))
		} else {
			button = tg.NewInlineKeyboardButtonURL("Перейти в беседу", "https://t.me/AdventureChatBot")
		}
	}

	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			button,
		),
	)
	return keyboard, msgText
}
