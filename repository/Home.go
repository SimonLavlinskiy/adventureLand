package repository

import (
	"errors"
	"fmt"
	v "github.com/spf13/viper"
	"math/rand"
	"strings"
)

type House struct {
	HouseMap Map
	Cells    []Cell
}

type HouseWithGarden struct {
	HouseMap    Map
	HouseCells  []Cell
	GardenMap   Map
	GardenCells []Cell
}

func MapHouse() Map {
	return Map{
		Name:             "Дом, милый дом!",
		SizeX:            10,
		SizeY:            10,
		StartX:           3,
		StartY:           3,
		DayType:          "alwaysDay",
		EmptySpaceSymbol: "〰️",
	}
}

func (m Map) GetCellsOfEmptyHouse() []Cell {
	var Cells []Cell

	for x := 0; x <= m.SizeX; x++ {
		for y := 0; y <= m.SizeY; y++ {

			if isCorner(x, y, m) {
				Cells = append(Cells, CornerCell(x, y, m))
			}

			if isW, numType := isWall(x, y); isW {
				if numType == "even" {
					Cells = append(Cells, WallCell(x, y, m, "white"))
				} else {
					Cells = append(Cells, WallCell(x, y, m, "red"))
				}
			}

			if isFloor(x, y) {
				Cells = append(Cells, FloorCell(x, y, m, "brown"))
			}

			if isDoor(x, y) {
				Cells = append(Cells, DoorCell(x, y, m))
			}

			if isWindow(x, y) {
				Cells = append(Cells, WindowCell(x, y, m))
			}
		}
	}

	return Cells
}

func isCorner(x int, y int, m Map) bool {
	if x < 3 && (y < 3 || y > m.SizeY-3) || x > m.SizeX-3 && (y < 3 || y > m.SizeY-3) {
		return true
	}
	return false
}

func isWall(x int, y int) (bool, string) {
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

func isFloor(x int, y int) bool {
	if x >= 3 && x <= 7 && y >= 3 && y <= 7 {
		return true
	}
	return false
}

func isDoor(x int, y int) bool {
	if x == 3 && y == 2 {
		return true
	}
	return false
}

func isWindow(x int, y int) bool {
	if (x == 1 && (y == 4 || y == 6)) || (x == 9 && (y == 4 || y == 6)) {
		return true
	} else if (y == 1 || y == 9) && (x == 4 || x == 6) {
		return true
	}
	return false
}

func CornerCell(x int, y int, m Map) Cell {
	cellType := "cell"
	return Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    v.GetString("elem.black"),
		CanStep: false,
		Type:    &cellType,
	}
}

func WallCell(x int, y int, m Map, color string) Cell {
	cellType := "cell"
	return Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    v.GetString(fmt.Sprintf("elem.%s", color)),
		CanStep: false,
		Type:    &cellType,
	}
}

func FloorCell(x int, y int, m Map, color string) Cell {
	cellType := "cell"
	return Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    v.GetString(fmt.Sprintf("elem.%s", color)),
		CanStep: true,
		Type:    &cellType,
	}
}

func DoorCell(x int, y int, m Map) Cell {
	cellType := "teleport"
	tpId := 5
	return Cell{
		MapsId:     int(m.ID),
		AxisX:      x,
		AxisY:      y,
		View:       v.GetString("elem.door"),
		CanStep:    false,
		Type:       &cellType,
		TeleportID: &tpId,
	}
}

func WindowCell(x int, y int, m Map) Cell {
	cellType := "cell"

	windows := strings.Fields(fmt.Sprintf("%s %s %s", v.GetString("elem.day_window"), v.GetString("elem.evening_window"), v.GetString("elem.night_window")))

	randomNum := rand.Intn(len(windows))
	view := windows[randomNum]

	return Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    view,
		CanStep: false,
		Type:    &cellType,
	}
}

func (u User) CreateUserHouse() error {
	if *u.Money <= v.GetInt("main_info.cost_of_house") {
		return errors.New("user doesn't have money enough")
	}
	m := MapHouse().CreateMap()
	cells := m.GetCellsOfEmptyHouse()

	err := CreateCells(cells)
	if err {
		return errors.New("map is not created")
	}

	return nil
}
