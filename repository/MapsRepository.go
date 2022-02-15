package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strings"
	"time"
)

type Map struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"embedded"`
	SizeX  int    `gorm:"embedded"`
	SizeY  int    `gorm:"embedded"`
	StartX int    `gorm:"embedded"`
	StartY int    `gorm:"embedded"`
}

type MapButtons struct {
	Up     string
	Left   string
	Right  string
	Down   string
	Center string
}

func DefaultButtons(center string) MapButtons {
	return MapButtons{
		Up:     "🔼",
		Left:   "◀️️",
		Right:  "▶️",
		Down:   "🔽",
		Center: center,
	}
}

type UserMap struct {
	leftIndent  int
	rightIndent int
	upperIndent int
	downIndent  int
}

var displaySize = 6

func DefaultUserMap(location Location, displaySize int) UserMap {
	return UserMap{
		leftIndent:  *location.AxisX - displaySize,
		rightIndent: *location.AxisX + displaySize,
		upperIndent: *location.AxisY + displaySize,
		downIndent:  *location.AxisY - displaySize,
	}
}

func GetUserMap(update tgbotapi.Update) Map {
	resLocation := GetOrCreateMyLocation(update)
	result := Map{}

	err := config.Db.Where(&Map{Name: resLocation.Map}).First(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetMyMap(update tgbotapi.Update) (textMessage string, buttons tgbotapi.ReplyKeyboardMarkup) {
	resUser := GetUser(User{TgId: uint(update.Message.From.ID)})
	resLocation := GetOrCreateMyLocation(update)
	resMap := GetUserMap(update)
	mapSize := CalculateUserMapBorder(resLocation, resMap)
	messageMap := "*Карта*: _" + resLocation.Map + "_ *X*: _" + ToString(*resLocation.AxisX) + "_  *Y*: _" + ToString(*resLocation.AxisY) + "_"

	type Point = [2]int
	m := map[Point]Cellule{}

	var result []Cellule

	err := config.Db.
		Preload("Item").
		Preload("Teleport").
		Preload("Item.CanTakeWith").
		Where(Cellule{Map: resLocation.Map}).
		Where("axis_x >= " + ToString(mapSize.leftIndent)).
		Where("axis_x <= " + ToString(mapSize.rightIndent)).
		Where("axis_y >= " + ToString(mapSize.downIndent)).
		Where("axis_y <= " + ToString(mapSize.upperIndent)).
		Order("axis_x ASC").
		Order("axis_y ASC").
		Find(&result).Error

	if err != nil {
		panic(err)
	}

	for _, cell := range result {
		m[Point{cell.AxisX, cell.AxisY}] = cell
	}

	buttons = CalculateButtonMap(resLocation, resUser, m)

	m[Point{*resLocation.AxisX, *resLocation.AxisY}] = Cellule{View: resUser.Avatar, ID: m[Point{*resLocation.AxisX, *resLocation.AxisY}].ID}

	Maps := configurationMap(mapSize, resMap, resLocation, resUser, m)

	for i, row := range Maps {
		if i >= mapSize.downIndent && i <= mapSize.upperIndent {
			messageMap = strings.Join(row, "") + "\n" + messageMap
		}
	}

	return messageMap, buttons
}

func CalculateUserMapBorder(resLocation Location, resMap Map) UserMap {
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
		mapSize.leftIndent = 0
		mapSize.rightIndent = d * 2
	}
	if *resLocation.AxisY < d {
		mapSize.downIndent = 0
		mapSize.upperIndent = d * 2
	}
	if mapSize.rightIndent >= resMap.SizeX && *resLocation.AxisX > d {
		mapSize.leftIndent = resMap.SizeX - d*2
		mapSize.rightIndent = resMap.SizeX
	}
	if mapSize.upperIndent >= resMap.SizeY && *resLocation.AxisY > d {
		mapSize.downIndent = resMap.SizeY - d*2
		mapSize.upperIndent = resMap.SizeY
	}

	return mapSize
}

func CalculateButtonMap(resLocation Location, resUser User, m map[[2]int]Cellule) tgbotapi.ReplyKeyboardMarkup {
	type Point = [2]int

	buttons := DefaultButtons(resUser.Avatar)

	var CellsAroundUser = []Cellule{}
	CellsAroundUser = append(CellsAroundUser,
		m[Point{*resLocation.AxisX, *resLocation.AxisY + 1}],
		m[Point{*resLocation.AxisX, *resLocation.AxisY - 1}],
		m[Point{*resLocation.AxisX + 1, *resLocation.AxisY}],
		m[Point{*resLocation.AxisX - 1, *resLocation.AxisY}],
	)

	buttons = PutButton(CellsAroundUser, buttons, resUser)

	return CreateMapKeyboard(buttons)
}

func CreateMapKeyboard(buttons MapButtons) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🧥🎒"),
			tgbotapi.NewKeyboardButton(buttons.Up),
			tgbotapi.NewKeyboardButton("🎒"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.Left),
			tgbotapi.NewKeyboardButton(buttons.Center),
			tgbotapi.NewKeyboardButton(buttons.Right),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬛"),
			tgbotapi.NewKeyboardButton(buttons.Down),
			tgbotapi.NewKeyboardButton("Меню"),
		),
	)
}

func configurationMap(mapSize UserMap, resMap Map, resLocation Location, user User, m map[[2]int]Cellule) [][]string {
	currentTime := time.Now()
	day := "06" <= currentTime.Format("15") && currentTime.Format("15") <= "23"
	type Point [2]int
	Maps := make([][]string, resMap.SizeY+1)

	if day || resLocation.Map != "Main Place" {
		DayMap(Maps, mapSize, m, resLocation)
	} else {
		NigthMap(Maps, mapSize, m, resLocation, user)
	}
	return Maps
}

func DayMap(Maps [][]string, mapSize UserMap, m map[[2]int]Cellule, resLocation Location) {
	type Point [2]int

	for y := range Maps {
		for x := mapSize.leftIndent; x <= mapSize.rightIndent; x++ {
			if m[Point{x, y}].ID != 0 || m[Point{x, y}] == m[Point{*resLocation.AxisX, *resLocation.AxisY}] {
				appendVisibleUserZone(m, x, y, Maps)
			} else {
				Maps[y] = append(Maps[y], "\U0001FAA8")
			}
		}
	}
}

func NigthMap(Maps [][]string, mapSize UserMap, m map[[2]int]Cellule, resLocation Location, user User) {
	type Point [2]int

	for y := range Maps {
		for x := mapSize.leftIndent; x <= mapSize.rightIndent; x++ {
			if calculateNightMap(user, resLocation, x, y) && m[Point{x, y}].ID != 0 {
				appendVisibleUserZone(m, x, y, Maps)
			} else if (y+x)%2 == 1 || m[Point{x, y}].ID != 0 && (y+x)%2 == 1 {
				Maps[y] = append(Maps[y], "⬛️")
			} else if (y+x)%2 == 0 || m[Point{x, y}].ID != 0 && (y+x)%2 == 0 {
				Maps[y] = append(Maps[y], "✨")
			} else {
				Maps[y] = append(Maps[y], "")
			}
		}
	}
}

func calculateNightMap(user User, location Location, x int, y int) bool {

	if (*location.AxisX == x-2 || *location.AxisX == x+2) && (*location.AxisY >= y-1 && *location.AxisY <= y+1) {
		return true
	}
	if (*location.AxisX >= x-1 && *location.AxisX <= x+1) && (*location.AxisY >= y-2 && *location.AxisY <= y+2) {
		return true
	}

	if user.LeftHandId != nil && user.LeftHand.Type == "light" || user.RightHandId != nil && user.RightHand.Type == "light" {
		if (*location.AxisX == x-4 || *location.AxisX == x+4) && (*location.AxisY >= y-1 && *location.AxisY <= y+1) {
			return true
		}
		if (*location.AxisX == x-3 || *location.AxisX == x+3) && (*location.AxisY >= y-2 && *location.AxisY <= y+2) {
			return true
		}
		if (*location.AxisX == x-2 || *location.AxisX == x+2) && (*location.AxisY >= y-3 && *location.AxisY <= y+3) {
			return true
		}
		if (*location.AxisX >= x-1 && *location.AxisX <= x+1) && (*location.AxisY >= y-4 && *location.AxisY <= y+4) {
			return true
		}
	}

	return false
}

func appendVisibleUserZone(m map[[2]int]Cellule, x int, y int, Maps [][]string) {
	type Point [2]int

	if m[Point{x, y}].Type != nil &&
		*m[Point{x, y}].Type == "item" &&
		m[Point{x, y}].ItemID != nil &&
		*m[Point{x, y}].CountItem != 0 {
		Maps[y] = append(Maps[y], m[Point{x, y}].Item.View)
	} else {
		Maps[y] = append(Maps[y], m[Point{x, y}].View)
	}
}

func PutButton(CellsAroundUser []Cellule, buttons MapButtons, resUser User) MapButtons {

	switch true {
	case IsDefaultCell(CellsAroundUser[0]):
		buttons.Up = CellsAroundUser[0].View
	case IsTeleport(CellsAroundUser[0]):
		buttons.Up += " " + resUser.Avatar + " " + CellsAroundUser[0].View
	case IsItem(CellsAroundUser[0]):
		buttons.Up = isItemCost(CellsAroundUser[0], buttons.Up, resUser)
	case CellsAroundUser[0].ID == 0:
		buttons.Up = "🚫"
	}

	switch true {
	case IsDefaultCell(CellsAroundUser[1]):
		buttons.Down = CellsAroundUser[1].View
	case IsTeleport(CellsAroundUser[1]):
		buttons.Down += " " + resUser.Avatar + " " + CellsAroundUser[1].View
	case IsItem(CellsAroundUser[1]):
		buttons.Down = isItemCost(CellsAroundUser[1], buttons.Down, resUser)
	case CellsAroundUser[1].ID == 0:
		buttons.Down = "🚫"
	}

	switch true {
	case IsDefaultCell(CellsAroundUser[2]):
		buttons.Right = CellsAroundUser[2].View
	case IsTeleport(CellsAroundUser[2]):
		buttons.Right += " " + resUser.Avatar + " " + CellsAroundUser[2].View
	case IsItem(CellsAroundUser[2]):
		buttons.Right = isItemCost(CellsAroundUser[2], buttons.Right, resUser)
	case CellsAroundUser[2].ID == 0:
		buttons.Right = "🚫"
	}

	switch true {
	case IsDefaultCell(CellsAroundUser[3]):
		buttons.Left = CellsAroundUser[3].View
	case IsTeleport(CellsAroundUser[3]):
		buttons.Left += " " + resUser.Avatar + " " + CellsAroundUser[3].View
	case IsItem(CellsAroundUser[3]):
		buttons.Left = isItemCost(CellsAroundUser[3], buttons.Left, resUser)
	case CellsAroundUser[3].ID == 0:
		buttons.Left = "🚫"
	}

	return buttons
}

func IsDefaultCell(cell Cellule) bool {
	if cell.Type != nil && *cell.Type == "cell" && !cell.CanStep {
		return true
	}
	return false
}

func IsTeleport(cell Cellule) bool {
	if cell.Type != nil && *cell.Type == "teleport" && cell.TeleportID != nil {
		return true
	}
	return false
}

func IsItem(cell Cellule) bool {
	if cell.Type != nil && *cell.Type == "item" && cell.ItemID != nil && *cell.CountItem > 0 {
		return true
	}
	return false
}

func IsSpecialItem(cell Cellule, user User) string {
	if user.LeftHand != nil && user.LeftHand.Type == cell.Item.CanTakeWith.Type {
		return user.LeftHand.View
	} else if user.RightHand != nil && user.RightHand.Type == cell.Item.CanTakeWith.Type {
		return user.RightHand.View
	}
	return "🚷"
}

func isItemCost(cell Cellule, button string, resUser User) string {
	if cell.Item.Cost != nil && *cell.Item.Cost > 0 {
		if cell.Item.CanTakeWith != nil {
			res := IsSpecialItem(cell, resUser)
			button = res + " " + button + " " + cell.Item.View + " (" + ToString(*cell.Item.Cost) + "💰)"
		} else {
			button = "👋 " + button + " " + cell.Item.View + " ( " + ToString(*cell.Item.Cost) + "💰 )"
		}
	} else {
		if cell.Item.CanTakeWith != nil {
			res := IsSpecialItem(cell, resUser)
			button = res + " " + button + " " + cell.Item.View
		} else {
			button = "👋 " + button + " " + cell.Item.View
		}
	}
	return button
}
