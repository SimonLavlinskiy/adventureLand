package helpers

import (
	"fmt"
	v "github.com/spf13/viper"
	"math/rand"
	"project0/src/models"
	"strings"
)

func IsCorner(x int, y int, m models.Map) bool {
	if x < 3 && (y < 3 || y > m.SizeY-3) || x > m.SizeX-3 && (y < 3 || y > m.SizeY-3) {
		return true
	}
	return false
}

func IsWall(x int, y int) (bool, string) {
	var numberType string
	var wall bool

	if y == 0 && x > 2 && x <= 7 {
		wall = true
	} else if y == 1 && (x == 3 || x == 5 || x == 7) {
		wall = true
	} else if y == 2 && x > 3 && x <= 7 {
		wall = true
	} else if (y == 3 || y == 5 || y == 7) && (x < 3 || x > 7) {
		wall = true
	} else if (y == 4 || y == 6) && (x == 0 || x == 2 || x == 8 || x == 10) {
		wall = true
	} else if y == 8 && x > 2 && x <= 7 {
		wall = true
	} else if y == 9 && (x == 3 || x == 5 || x == 7) {
		wall = true
	} else if y == 10 && x > 2 && x <= 7 {
		wall = true
	}

	if (y+x)%2 == 1 {
		numberType = "odd"
	} else {
		numberType = "even"
	}

	return wall, numberType
}

func IsFloor(x int, y int) bool {
	if x >= 3 && x <= 7 && y >= 3 && y <= 7 {
		return true
	}
	return false
}

func IsDoor(x int, y int) bool {
	if x == 3 && y == 2 {
		return true
	}
	return false
}

func IsWindow(x int, y int) bool {
	if (x == 1 && (y == 4 || y == 6)) || (x == 9 && (y == 4 || y == 6)) {
		return true
	} else if (y == 1 || y == 9) && (x == 4 || x == 6) {
		return true
	}
	return false
}

func CornerCell(x int, y int, m models.Map) models.Cell {
	cellType := "cell"
	return models.Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    v.GetString("elem.black"),
		CanStep: false,
		Type:    &cellType,
	}
}

func WallCell(x int, y int, m models.Map, color string) models.Cell {
	cellType := "cell"
	return models.Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    v.GetString(fmt.Sprintf("elem.%s", color)),
		CanStep: false,
		Type:    &cellType,
	}
}

func FloorCell(x int, y int, m models.Map, color string) models.Cell {
	cellType := "cell"
	return models.Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    v.GetString(fmt.Sprintf("elem.%s", color)),
		CanStep: true,
		Type:    &cellType,
	}
}

func DoorCell(x int, y int, m models.Map) models.Cell {
	cellType := "teleport"
	tpId := 5
	return models.Cell{
		MapsId:     int(m.ID),
		AxisX:      x,
		AxisY:      y,
		View:       v.GetString("elem.door"),
		CanStep:    false,
		Type:       &cellType,
		TeleportID: &tpId,
	}
}

func WindowCell(x int, y int, m models.Map) models.Cell {
	cellType := "cell"

	windows := strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("elem.day_window"), v.GetString("elem.evening_window"), v.GetString("elem.night_window")))

	randomNum := rand.Intn(len(windows))
	view := windows[randomNum]

	return models.Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    view,
		CanStep: false,
		Type:    &cellType,
	}
}
