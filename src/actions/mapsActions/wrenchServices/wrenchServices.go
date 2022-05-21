package wrenchServices

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"strings"
)

func AllReceiptsMsg() string {
	receipts := repositories.GetReceipts()
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

func PutCountComponent(char []string) (buttons tg.InlineKeyboardMarkup) {
	userItemId := char[helpers.ToInt(char[2])+(4+helpers.ToInt(char[2])*2)] // char[x + (4+x*2 )] = char[4]
	userItem := models.UserItem{ID: helpers.ToInt(userItemId)}.UserGetUserItem()

	buttons = ChangeCountUserItemKeyboard(char, userItem)
	return buttons
}

func UserCraftItem(user models.User, receipt *models.Receipt, charData []string) (msgText string, buttons tg.InlineKeyboardMarkup) {

	if receipt == nil {
		msgText, buttons = Workbench(nil, charData)
		msgText = fmt.Sprintf("%s%sТакого рецепта не существует!", msgText, v.GetString("msg_separator"))
		return msgText, buttons
	}

	resultItem := models.UserItem{UserId: int(user.ID), ItemId: receipt.ItemResultID}.UserGetUserItem()

	if resultItem.Item.MaxCountUserHas != nil && *receipt.ItemResultCount+*resultItem.Count > *resultItem.Item.MaxCountUserHas {
		msgText, buttons = Workbench(nil, charData)
		msgText = fmt.Sprintf("%s%sВы не можете иметь больше, чем %d %s!\nСори... такие правила(", msgText, v.GetString("msg_separator"), *resultItem.Item.MaxCountUserHas, resultItem.Item.View)
		return msgText, buttons
	}

	if receipt.Component1ID != nil && receipt.Component1Count != nil {
		userItem := models.UserItem{UserId: int(user.ID), ItemId: *receipt.Component1ID}.UserGetUserItem()
		countItem1 := *userItem.Count - *receipt.Component1Count
		user.UpdateUserItem(models.UserItem{ID: userItem.ID, ItemId: *receipt.Component1ID, Count: &countItem1}) // CountUseLeft: resultItem.CountUseLeft
	}
	if receipt.Component2ID != nil && receipt.Component2Count != nil {
		userItem := models.UserItem{UserId: int(user.ID), ItemId: *receipt.Component2ID}.UserGetUserItem()
		countItem2 := *userItem.Count - *receipt.Component2Count
		user.UpdateUserItem(models.UserItem{ID: userItem.ID, ItemId: *receipt.Component2ID, Count: &countItem2}) // CountUseLeft: resultItem.CountUseLeft
	}
	if receipt.Component3ID != nil && receipt.Component3Count != nil {
		userItem := models.UserItem{UserId: int(user.ID), ItemId: *receipt.Component3ID}.UserGetUserItem()
		countItem3 := *userItem.Count - *receipt.Component3Count
		user.UpdateUserItem(models.UserItem{ID: userItem.ID, ItemId: *receipt.Component3ID, Count: &countItem3}) // CountUseLeft: resultItem.CountUseLeft
	}

	if *resultItem.Count == 0 {
		resultItem.CountUseLeft = resultItem.Item.CountUse
	}
	*resultItem.Count = *resultItem.Count + *receipt.ItemResultCount
	user.UpdateUserItem(models.UserItem{ID: resultItem.ID, Count: resultItem.Count, CountUseLeft: resultItem.CountUseLeft})

	charData = strings.Fields("workbench usPoint 0 1stComp nil 0 2ndComp nil 0 3rdComp nil 0")
	msgText, buttons = Workbench(nil, charData)
	msgText = fmt.Sprintf("%s%sСупер! Ты получил %s %d шт. %s!", msgText, v.GetString("msg_separator"), resultItem.Item.View, *receipt.ItemResultCount, receipt.ItemResult.Name)

	return msgText, buttons
}

func Workbench(cell *models.Cell, char []string) (msgText string, buttons tg.InlineKeyboardMarkup) {
	var charData []string
	if cell != nil && !cell.IsWorkbench() {
		msgText = "Здесь нет верстака!"
		return msgText, buttons
	}

	if cell != nil {
		charData = strings.Fields("workbench usPoint 0 1stComp nil 0 2ndComp nil 0 3rdComp nil 0")
	} else {
		charData = strings.Fields(fmt.Sprintf("workbench usPoint %s 1stComp %s %s 2ndComp %s %s 3rdComp %s %s", char[2], char[4], char[5], char[7], char[8], char[10], char[11]))
	}

	msgText = OpenWorkbenchMessage(charData)
	buttons = WorkbenchKeyboard(charData)

	return msgText, buttons
}

func OpenWorkbenchMessage(char []string) string {
	// char = workbench usPoint 0 1stComp: id 0 2ndComp id 0 3rdComp id 0

	fstCnt := getViewEmojiForMsg(char, 0)
	secCnt := getViewEmojiForMsg(char, 1)
	trdCnt := getViewEmojiForMsg(char, 2)

	fstComponentView := viewComponent(char[4])
	secComponentView := viewComponent(char[7])
	trdComponentView := viewComponent(char[10])

	cellUser := helpers.ToInt(char[2])
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
		component := models.UserItem{ID: helpers.ToInt(id)}.UserGetUserItem()
		return component.Item.View
	}
	return "⚪"
}
