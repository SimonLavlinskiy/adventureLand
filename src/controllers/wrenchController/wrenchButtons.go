package wrenchController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
	"project0/src/services/helpers"
	str "strings"
)

func ChangeCountUserItemKeyboard(charData []string, item models.UserItem) tg.InlineKeyboardMarkup {
	charDone := fmt.Sprintf("%s usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.workbench"), charData[2], charData[4], charData[5], charData[7], charData[8], charData[10], charData[11])
	itemCount := helpers.ToInt(charData[helpers.ToInt(charData[2])+(5+helpers.ToInt(charData[2])*2)])
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

	subData[helpers.ToInt(charData[2])+(5+helpers.ToInt(charData[2])*2)] = subCount
	appData[helpers.ToInt(charData[2])+(5+helpers.ToInt(charData[2])*2)] = appCount

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

func WorkbenchKeyboard(char []string) tg.InlineKeyboardMarkup {
	leftArrow := "⬅️"
	rightArrow := "➡️"
	userPointer := helpers.ToInt(char[2])

	defaultData := fmt.Sprintf("usPoint %d 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", userPointer, char[4], char[5], char[7], char[8], char[10], char[11])
	rightArrowData := fmt.Sprintf("%s usPoint %d 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.workbench"), userPointer+1, char[4], char[5], char[7], char[8], char[10], char[11])
	leftArrowData := fmt.Sprintf("%s usPoint %d 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", v.GetString("callback_char.workbench"), userPointer-1, char[4], char[5], char[7], char[8], char[10], char[11])
	putItemData := fmt.Sprintf("%s %s", v.GetString("callback_char.put_item"), defaultData)
	putCountItemData := fmt.Sprintf("%s %s", v.GetString("callback_char.put_count_item"), defaultData)

	makeNewItem := fmt.Sprintf("%s %s", v.GetString("callback_char.make_new_item"), defaultData)

	if userPointer == 0 {
		leftArrow = "✖️"
		leftArrowData = fmt.Sprintf("workbench %s", defaultData)
	} else if userPointer == 2 {
		rightArrow = "✖️"
		rightArrowData = fmt.Sprintf("workbench %s", defaultData)
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
			tg.NewInlineKeyboardButtonData("Рецепты 📚", fmt.Sprintf("%s %s", v.GetString("callback_char.receipt"), defaultData)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
		),
	)
}

func ReturnToWorkbench(char []string) tg.InlineKeyboardMarkup {
	defaultData := fmt.Sprintf("workbench usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", char[2], char[4], char[5], char[7], char[8], char[10], char[11])
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Назад", defaultData),
		),
	)
}
