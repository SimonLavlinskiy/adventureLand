package mapController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/instrumentController"
	"project0/src/models"
)

func CalculateUserMapBorder(resLocation models.Location, resMap models.Map) models.UserMap {
	d := displaySize
	if d*2 > resMap.SizeX || d*2 > resMap.SizeY {
		if resMap.SizeX-resMap.SizeY > 0 {
			d = resMap.SizeY / 2
		} else {
			d = resMap.SizeX / 2
		}
	}
	mapSize := DefaultUserMap(resLocation, d)

	if *resLocation.AxisX < d {
		mapSize.LeftIndent = 0
		mapSize.RightIndent = d * 2
	}
	if *resLocation.AxisY < d {
		mapSize.DownIndent = 0
		mapSize.UpperIndent = d * 2
	}
	if mapSize.RightIndent >= resMap.SizeX && *resLocation.AxisX > d {
		mapSize.LeftIndent = resMap.SizeX - d*2
		mapSize.RightIndent = resMap.SizeX
	}
	if mapSize.UpperIndent >= resMap.SizeY && *resLocation.AxisY > d {
		mapSize.DownIndent = resMap.SizeY - d*2
		mapSize.UpperIndent = resMap.SizeY
	}

	return mapSize
}

func CalculateButtonMap(resLocation models.Location, resUser models.User, m map[[2]int]models.Cell) tg.InlineKeyboardMarkup {
	type Point = [2]int

	buttons := DefaultButtons(resUser.Avatar)

	var CellsAroundUser []models.Cell
	CellsAroundUser = append(CellsAroundUser,
		m[Point{*resLocation.AxisX, *resLocation.AxisY + 1}],
		m[Point{*resLocation.AxisX, *resLocation.AxisY - 1}],
		m[Point{*resLocation.AxisX + 1, *resLocation.AxisY}],
		m[Point{*resLocation.AxisX - 1, *resLocation.AxisY}],
	)

	buttons = PutButton(CellsAroundUser, buttons, resUser)

	return CreateMapKeyboard(buttons)
}

func PutButton(CellsAroundUser []models.Cell, btn models.MapButtons, user models.User) models.MapButtons {
	for i, cell := range CellsAroundUser {
		switch true {
		case cell.IsDefaultCell():
			text := cell.View
			switch i {
			case 0:
				btn.Up, btn.UpData = text, text
			case 1:
				btn.Down, btn.DownData = text, text
			case 2:
				btn.Right, btn.RightData = text, text
			case 3:
				btn.Left, btn.LeftData = text, text
			}
		case cell.IsTeleport() || cell.IsHome():
			button := fmt.Sprintf("%s%s", user.Avatar, cell.View)
			data := fmt.Sprintf("move %d", cell.ID)
			switch i {
			case 0:
				btn.Up = button
				btn.UpData = data
			case 1:
				btn.Down = button
				btn.DownData = data
			case 2:
				btn.Right = button
				btn.RightData = data
			case 3:
				btn.Left = button
				btn.LeftData = data
			}
		case cell.IsWorkbench() || cell.IsQuest() || cell.IsChat():
			var data string
			if cell.IsWorkbench() {
				data = fmt.Sprintf("wrench %d", cell.ID)
			} else if cell.IsQuest() {
				data = fmt.Sprintf("quests %d", cell.ID)
			} else if cell.IsChat() {
				data = fmt.Sprintf("chat %d", cell.ID)
			}
			switch i {
			case 0:
				btn.Up += cell.Item.View
				btn.UpData = data
			case 1:
				btn.Down += cell.Item.View
				btn.DownData = data
			case 2:
				btn.Right += cell.Item.View
				btn.RightData = data
			case 3:
				btn.Left += cell.Item.View
				btn.LeftData = data
			}
		case cell.IsItem() || cell.IsSwap() || cell.IsBox(user):
			switch i {
			case 0:
				btn.Up, btn.UpData = ViewItemButton(cell, user)
			case 1:
				btn.Down, btn.DownData = ViewItemButton(cell, user)
			case 2:
				btn.Right, btn.RightData = ViewItemButton(cell, user)
			case 3:
				btn.Left, btn.LeftData = ViewItemButton(cell, user)
			}
		case cell.IsWordleGame():
			switch i {
			case 0:
				btn.UpData = fmt.Sprintf("wordle_game up %s", cell.View)
				btn.Up = fmt.Sprintf("%s%s", btn.Up, cell.View)
			case 1:
				btn.DownData = fmt.Sprintf("wordle_game down %s", cell.View)
				btn.Down = fmt.Sprintf("%s%s", btn.Down, cell.View)
			case 2:
				btn.RightData = fmt.Sprintf("wordle_game right %s", cell.View)
				btn.Right = fmt.Sprintf("%s%s", btn.Right, cell.View)
			case 3:
				btn.LeftData = fmt.Sprintf("wordle_game left %s", cell.View)
				btn.Left = fmt.Sprintf("%s%s", btn.Left, cell.View)
			}
		case cell.ID == 0:
			switch i {
			case 0:
				btn.Up, btn.UpData = "🚫", "🚫"
			case 1:
				btn.Down, btn.DownData = "🚫", "🚫"
			case 2:
				btn.Right, btn.RightData = "🚫", "🚫"
			case 3:
				btn.Left, btn.LeftData = "🚫", "🚫"
			}
		default:
			switch i {
			case 0:
				btn.UpData = fmt.Sprintf("move %d", cell.ID)
			case 1:
				btn.DownData = fmt.Sprintf("move %d", cell.ID)
			case 2:
				btn.RightData = fmt.Sprintf("move %d", cell.ID)
			case 3:
				btn.LeftData = fmt.Sprintf("move %d", cell.ID)
			}
		}
	}

	return btn
}

func CreateMapKeyboard(buttons models.MapButtons) tg.InlineKeyboardMarkup {
	nearUsers := "Это не кнопка"

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Вещи 🧦", "goodsMoving"),
			tg.NewInlineKeyboardButtonData(buttons.Up, buttons.UpData),
			tg.NewInlineKeyboardButtonData("Рюкзак 🎒", "category"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(buttons.Left, buttons.LeftData),
			tg.NewInlineKeyboardButtonData(buttons.Center, buttons.Center),
			tg.NewInlineKeyboardButtonData(buttons.Right, buttons.RightData),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(nearUsers, nearUsers),
			tg.NewInlineKeyboardButtonData(buttons.Down, buttons.DownData),
			tg.NewInlineKeyboardButtonData(v.GetString("user_location.menu"), v.GetString("user_location.menu")),
		),
	)
}

func ViewItemButton(cell models.Cell, user models.User) (btn string, btnData string) {
	instrumentsUserCanUse := instrumentController.GetInstrumentsUserCanUse(user, cell)

	if len(instrumentsUserCanUse) > 0 {
		btn = fmt.Sprintf("🛠❓%s", cell.Item.View)
		btnData = fmt.Sprintf("chooseInstrument %d", cell.ID)
	} else {
		btn = cell.Item.View
		btnData = cell.Item.View
	}

	return btn, btnData
}
