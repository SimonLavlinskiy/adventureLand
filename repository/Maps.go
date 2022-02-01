package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strconv"
	"strings"
)

type Map struct {
	ID    uint
	Name  string
	SizeX uint
	SizeY uint
}

type MapButtons struct {
	Up    string
	Left  string
	Right string
	Down  string
}

func DefaultButtons() MapButtons {
	return MapButtons{
		Up:    "🔼",
		Left:  "◀️",
		Right: "▶️",
		Down:  "🔽",
	}
}

type UserMap struct {
	leftIndent  uint
	rightIndent uint
	upperIndent uint
	downIndent  uint
}

var displayMapSize = uint(5)

func DefaultUserMap(location Location) UserMap {
	return UserMap{
		leftIndent:  Dif(location.AxisX, displayMapSize),
		rightIndent: Sum(location.AxisX, displayMapSize),
		upperIndent: Sum(location.AxisY, displayMapSize),
		downIndent:  Dif(location.AxisY, displayMapSize),
	}
}

func GetMap(update tgbotapi.Update) Map {
	resLocation := GetOrCreateMyLocation(update)
	result := Map{}

	err := config.Db.Where(&Map{Name: resLocation.Map}).FirstOrCreate(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetMyMap(update tgbotapi.Update) (tgbotapi.MessageConfig, MapButtons) {
	resUser := GetOrCreateUser(update)
	resLocation := GetOrCreateMyLocation(update)
	resMap := GetMap(update)
	buttons := DefaultButtons()
	mapSize := DefaultUserMap(resLocation)
	var Maps [50][]string

	if resLocation.AxisX < mapSize.leftIndent {
		mapSize.leftIndent = uint(1)
		mapSize.rightIndent = displayMapSize * 2
	}
	if resLocation.AxisY < mapSize.downIndent {
		mapSize.downIndent = uint(1)
		mapSize.upperIndent = displayMapSize * 2
	}
	if mapSize.rightIndent > resMap.SizeX && resLocation.AxisX > displayMapSize {
		mapSize.leftIndent = resMap.SizeX - displayMapSize*2
		mapSize.rightIndent = resMap.SizeX
	}
	if mapSize.upperIndent > resMap.SizeY && resLocation.AxisY > displayMapSize {
		mapSize.downIndent = resMap.SizeY - displayMapSize*2
		mapSize.upperIndent = resMap.SizeY
	}

	messageMap := "*Карта*: _" + resLocation.Map + "_\n*X*: _" + strconv.FormatUint(uint64(resLocation.AxisX), 10) + "_  *Y*: _" + strconv.FormatUint(uint64(resLocation.AxisY), 10) + "_\n"

	var result []Cellule

	err := config.Db.Where("map = '" + resLocation.Map + "' and axis_x >= " + ToString(mapSize.leftIndent) + " and axis_x <= " + ToString(mapSize.rightIndent) + " and axis_y >= " + ToString(mapSize.downIndent) + " and axis_y <= " + ToString(mapSize.upperIndent)).Order("axis_x ASC").Order("axis_y ASC").Find(&result).Error
	if err != nil {
		panic(err)
	}

	type Point = [2]uint

	m := map[Point]Cellule{}
	for _, cell := range result {
		m[Point{cell.AxisX - 1, cell.AxisY - 1}] = cell
	}

	if cell := m[Point{resLocation.AxisX - 1, resLocation.AxisY}]; !cell.CanStep {
		buttons.Up = "🚫" + cell.View + "🚫"
	}
	if cell := m[Point{resLocation.AxisX - 1, resLocation.AxisY - 2}]; !cell.CanStep {
		buttons.Down = "🚫" + cell.View + "🚫"
	}
	if cell := m[Point{resLocation.AxisX, resLocation.AxisY - 1}]; !cell.CanStep {
		buttons.Right = "🚫" + cell.View + "🚫"
	}
	if cell := m[Point{resLocation.AxisX - 2, resLocation.AxisY - 1}]; !cell.CanStep {
		buttons.Left = "🚫" + cell.View + "🚫"
	}

	m[Point{resLocation.AxisX - 1, resLocation.AxisY - 1}] = Cellule{View: resUser.Avatar}

	for i := range Maps {
		for y := 0; y < 13; y++ {
			Maps[i] = append(Maps[i], m[Point{uint(y), uint(i)}].View)
		}
	}

	for _, row := range Maps {
		messageMap = strings.Join(row, ``) + "\n" + messageMap
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, messageMap), buttons
}

func Sum(x uint, y uint) uint {
	return x + y
}

func Dif(x uint, y uint) uint {
	return x - y
}

func ToString(uint uint) string {
	return strconv.FormatUint(uint64(uint), 10)
}
