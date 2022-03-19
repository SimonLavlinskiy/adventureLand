package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	v "github.com/spf13/viper"
	"math/rand"
	"strings"
	"time"
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
		StartX:           5,
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
				}
				Cells = append(Cells, WallCell(x, y, m, "red"))
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

	j, _ := json.Marshal(Cells)
	fmt.Println(string(j))

	return Cells
}

func isCorner(x int, y int, m Map) bool {
	if x < 3 && y < 3 || x > m.SizeX-3 && y > m.SizeY-3 {
		return true
	}
	return false
}

func isWall(x int, y int) (bool, string) {
	var numberType string
	var wall bool

	if (x <= 2 && (y == 3 || y == 5 || y == 7)) || ((x == 0 || x == 2) && (y == 4 || y == 6)) {
		wall = true
	} else if (x >= 7 && (y == 3 || y == 5 || y == 7)) || ((x == 8 || x == 10) && (y == 4 || y == 6)) {
		wall = true
	} else if (y <= 1 && (x == 3 || x == 5 || x == 7)) || ((y == 0 || y == 2) && (x == 4 || x == 6)) || (y == 2 && x > 5) {
		wall = true
	} else if (y >= 7 && (x == 3 || x == 5 || x == 7)) || ((y == 8 || y == 10) && (x == 4 || x == 6)) {
		wall = true
	}

	if y == 2 && x == 5 {
		wall = false
	}

	if (y+x)%2 == 1 {
		numberType = "odd"
	} else {
		numberType = "even"
	}

	fmt.Println((y+x)%2, numberType, x, y)
	return wall, numberType
}

func isFloor(x int, y int) bool {
	if x >= 3 && x <= 7 && y >= 3 && y <= 7 {
		return true
	}
	return false
}

func isDoor(x int, y int) bool {
	if x == 5 && y == 2 {
		return true
	}
	return false
}

func isWindow(x int, y int) bool {
	if (x == 1 && y == 6) || (x == 9 && (y == 6 || y == 8)) {
		return true
	} else if (y == 1 && (x == 6 || x == 8)) || (y == 9 && (x == 6 || x == 8)) {
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
	tpId := 1
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
	t := time.Now()
	var str []string

	if day := "06" <= t.Format("15") && t.Format("15") < "19"; day {
		str = strings.Fields(v.GetString("elem.day_window"))
	} else if evening := "19" <= t.Format("15") && t.Format("15") < "23"; evening {
		str = strings.Fields(v.GetString("elem.evening_window"))
	} else if night := "23" <= t.Format("15") || t.Format("15") < "06"; night {
		str = strings.Fields(v.GetString("elem.night_window"))
	}

	randomNum := rand.Intn(len(str))
	view := str[randomNum]

	return Cell{
		MapsId:  int(m.ID),
		AxisX:   x,
		AxisY:   y,
		View:    view,
		CanStep: false,
		Type:    &cellType,
	}
}

func CreateUserHouse() error {
	m := MapHouse().CreateMap()
	cells := m.GetCellsOfEmptyHouse()

	err := CreateCells(cells)
	if err {
		return errors.New("map is not created")
	}

	return nil
}
