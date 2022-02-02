package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strconv"
	"strings"
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
		Left:   "◀️",
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

func GetMap(m Map) Map {
	result := Map{}
	err := config.Db.Where(&m).FirstOrCreate(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetMyMap(update tgbotapi.Update) (textMessage string, buttons MapButtons) {
	resUser := GetOrCreateUser(update)
	resLocation := GetOrCreateMyLocation(update)
	resMap := GetUserMap(update)
	buttons = DefaultButtons(resUser.Avatar)
	mapSize := CalculateUserMapBorder(resLocation, resMap)

	var result []Cellule

	err := config.Db.Where("map = '" + resLocation.Map + "' and axis_x >= " + ToString(mapSize.leftIndent) + " and axis_x <= " + ToString(mapSize.rightIndent) + " and axis_y >= " + ToString(mapSize.downIndent) + " and axis_y <= " + ToString(mapSize.upperIndent)).Order("axis_x ASC").Order("axis_y ASC").Find(&result).Error
	if err != nil {
		panic(err)
	}

	var Maps [][]string
	Maps = make([][]string, resMap.SizeY+1)

	type Point = [2]int
	m := map[Point]Cellule{}

	for _, cell := range result {
		m[Point{cell.AxisX, cell.AxisY}] = cell
	}

	if cell := m[Point{*resLocation.AxisX, *resLocation.AxisY + 1}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Up = "🚫"
		} else if cell.Type == "teleport" {
			buttons.Up += cell.View + "🚶‍♂️"
		} else {
			buttons.Up = cell.View
		}
	}
	if cell := m[Point{*resLocation.AxisX, *resLocation.AxisY - 1}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Down = "🚫"
		} else if cell.Type == "teleport" {
			buttons.Down += cell.View + "🚶‍♂️"
		} else {
			buttons.Down = cell.View
		}
	}
	if cell := m[Point{*resLocation.AxisX + 1, *resLocation.AxisY}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Right = "🚫"
		} else {
			buttons.Right = cell.View
		}
	}
	if cell := m[Point{*resLocation.AxisX - 1, *resLocation.AxisY}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Left = "🚫"
		} else {
			buttons.Left = cell.View
		}
	}

	m[Point{*resLocation.AxisX, *resLocation.AxisY}] = Cellule{View: resUser.Avatar}

	for i := range Maps {
		for x := mapSize.leftIndent; x <= mapSize.rightIndent; x++ {
			if m[Point{x, i}].ID != 0 || m[Point{x, i}] == m[Point{*resLocation.AxisX, *resLocation.AxisY}] {
				Maps[i] = append(Maps[i], m[Point{x, i}].View)
			} else {
				Maps[i] = append(Maps[i], "\U0001FAA8")
			}
		}
	}

	messageMap := "*Карта*: _" + resLocation.Map + "_ *X*: _" + ToString(*resLocation.AxisX) + "_  *Y*: _" + ToString(*resLocation.AxisY) + "_"

	for i, row := range Maps {
		if i >= mapSize.downIndent && i <= mapSize.upperIndent {
			messageMap = strings.Join(row, ``) + "\n" + messageMap
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
