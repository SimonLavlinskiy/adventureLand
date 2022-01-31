package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strconv"
	"strings"
)

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

func GetMyMap(update tgbotapi.Update) (tgbotapi.MessageConfig, MapButtons) {
	var Maps [12][]string
	resUser := GetOrCreateUser(update)
	resLocation := GetOrCreateMyLocation(update)
	buttons := DefaultButtons()

	messageMap := "*Карта*: _" + resLocation.Map + "_\n*X*: _" + strconv.FormatUint(uint64(resLocation.AxisX), 10) + "_  *Y*: _" + strconv.FormatUint(uint64(resLocation.AxisY), 10) + "_\n"

	var result []Cellule

	err := config.Db.Where(Cellule{Map: resLocation.Map}).Order("axis_x ASC").Order("axis_y ASC").Find(&result).Error
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
		for y := 0; y < 7; y++ {
			Maps[i] = append(Maps[i], m[Point{uint(y), uint(i)}].View)
		}
	}

	for _, row := range Maps {
		messageMap = strings.Join(row, ``) + "\n" + messageMap
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, messageMap), buttons
}
