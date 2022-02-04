package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strconv"
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

	err := config.Db.Where(&Map{Name: resLocation.Map}).FirstOrCreate(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetMyMap(update tgbotapi.Update) (textMessage string, buttons tgbotapi.ReplyKeyboardMarkup) {
	currentTime := time.Now()
	day := "06" <= currentTime.Format("15") && currentTime.Format("15") <= "23"

	resUser := GetOrCreateUser(update)
	resLocation := GetOrCreateMyLocation(update)
	resMap := GetUserMap(update)
	mapSize := CalculateUserMapBorder(resLocation, resMap)
	messageMap := "*Карта*: _" + resLocation.Map + "_ *X*: _" + ToString(*resLocation.AxisX) + "_  *Y*: _" + ToString(*resLocation.AxisY) + "_"

	type Point = [2]int
	m := map[Point]Cellule{}

	var result []Cellule

	err := config.Db.Where("map = '" + resLocation.Map + "' and axis_x >= " + ToString(mapSize.leftIndent) + " and axis_x <= " + ToString(mapSize.rightIndent) + " and axis_y >= " + ToString(mapSize.downIndent) + " and axis_y <= " + ToString(mapSize.upperIndent)).Order("axis_x ASC").Order("axis_y ASC").Find(&result).Error
	if err != nil {
		panic(err)
	}

	for _, cell := range result {
		m[Point{cell.AxisX, cell.AxisY}] = cell
	}

	buttons = CalculateButtonMap(resLocation, resUser, m)

	var Maps [][]string
	Maps = make([][]string, resMap.SizeY+1)

	m[Point{*resLocation.AxisX, *resLocation.AxisY}] = Cellule{View: resUser.Avatar, ID: m[Point{*resLocation.AxisX, *resLocation.AxisY}].ID}

	for y := range Maps {
		for x := mapSize.leftIndent; x <= mapSize.rightIndent; x++ {
			if day || resLocation.Map != "Main Place" {
				if m[Point{x, y}].ID != 0 || m[Point{x, y}] == m[Point{*resLocation.AxisX, *resLocation.AxisY}] {
					Maps[y] = append(Maps[y], m[Point{x, y}].View)
				} else {
					Maps[y] = append(Maps[y], "\U0001FAA8")
				}
			} else {
				if calculateNightMap(resLocation, x, y) && m[Point{x, y}].ID != 0 {
					Maps[y] = append(Maps[y], m[Point{x, y}].View)
				} else if (y+x)%2 == 1 || m[Point{x, y}].ID != 0 && (y+x)%2 == 1 {
					Maps[y] = append(Maps[y], "⬛️")
				} else if (y+x)%2 == 0 || m[Point{x, y}].ID != 0 && (y+x)%2 == 0 {
					Maps[y] = append(Maps[y], "✨")
				} else {
					Maps[y] = append(Maps[y], "\U0001FAA8")
				}
			}
		}
	}

	for i, row := range Maps {
		if i >= mapSize.downIndent && i <= mapSize.upperIndent {
			messageMap = strings.Join(row, "") + "\n" + messageMap
		}
	}

	return messageMap, buttons
}

func ToString(int int) string {
	return strconv.FormatInt(int64(int), 10)
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

	if cell := m[Point{*resLocation.AxisX, *resLocation.AxisY + 1}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Up = "🚫"
		} else if cell.Type == "teleport" {
			buttons.Up += " " + resUser.Avatar + " " + cell.View
		} else {
			buttons.Up = cell.View
		}
	}
	if cell := m[Point{*resLocation.AxisX, *resLocation.AxisY - 1}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Down = "🚫"
		} else if cell.Type == "teleport" {
			buttons.Down += " " + resUser.Avatar + " " + cell.View
		} else {
			buttons.Down = cell.View
		}
	}
	if cell := m[Point{*resLocation.AxisX + 1, *resLocation.AxisY}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Right = "🚫"
		} else if cell.Type == "teleport" {
			buttons.Right += " " + resUser.Avatar + " " + cell.View
		} else {
			buttons.Right = cell.View
		}
	}
	if cell := m[Point{*resLocation.AxisX - 1, *resLocation.AxisY}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Left = "🚫"
		} else if cell.Type == "teleport" {
			buttons.Left += " " + resUser.Avatar + " " + cell.View
		} else {
			buttons.Left = cell.View
		}
	}

	return CreateMoveKeyboard(buttons)
}

func CreateMoveKeyboard(buttons MapButtons) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬛"),
			tgbotapi.NewKeyboardButton(buttons.Up),
			tgbotapi.NewKeyboardButton("⬛"),
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

func calculateNightMap(location Location, x int, y int) bool {
	if *location.AxisX == x-2 && (*location.AxisY == y-1 || *location.AxisY == y || *location.AxisY == y+1) {
		return true
	}
	if *location.AxisX == x-1 && (*location.AxisY == y-2 || *location.AxisY == y-1 || *location.AxisY == y || *location.AxisY == y+1 || *location.AxisY == y+2) {
		return true
	}
	if *location.AxisX == x && (*location.AxisY == y-2 || *location.AxisY == y-1 || *location.AxisY == y || *location.AxisY == y+1 || *location.AxisY == y+2) {
		return true
	}
	if *location.AxisX == x+1 && (*location.AxisY == y-2 || *location.AxisY == y-1 || *location.AxisY == y || *location.AxisY == y+1 || *location.AxisY == y+2) {
		return true
	}
	if *location.AxisX == x+2 && (*location.AxisY == y-1 || *location.AxisY == y || *location.AxisY == y+1) {
		return true
	}
	return false
}
