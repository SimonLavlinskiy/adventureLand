package houseController

import (
	"errors"
	v "github.com/spf13/viper"
	"project0/src/controllers/userController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
)

func MapHouse() models.Map {
	return models.Map{
		Name:             "Дом, милый дом!",
		SizeX:            10,
		SizeY:            10,
		StartX:           3,
		StartY:           3,
		DayType:          "alwaysDay",
		EmptySpaceSymbol: "〰️",
	}
}

func CreateUserHouse(u models.User) error {
	if *u.Money <= v.GetInt("main_info.cost_of_house") {
		return errors.New("user doesn't have money enough")
	}
	m := MapHouse().CreateMap()
	cells := GetCellsOfEmptyHouse(m)

	err := repositories.CreateCells(cells)
	if err {
		return errors.New("userMapConfiguration is not created")
	}

	userController.UserBuyHome(u, m)

	return nil
}

func GetCellsOfEmptyHouse(m models.Map) []models.Cell {
	var Cells []models.Cell

	for x := 0; x <= m.SizeX; x++ {
		for y := 0; y <= m.SizeY; y++ {

			if helpers.IsCorner(x, y, m) {
				Cells = append(Cells, helpers.CornerCell(x, y, m))
			}

			if isW, numType := helpers.IsWall(x, y); isW {
				if numType == "even" {
					Cells = append(Cells, helpers.WallCell(x, y, m, "white"))
				} else {
					Cells = append(Cells, helpers.WallCell(x, y, m, "red"))
				}
			}

			if helpers.IsFloor(x, y) {
				Cells = append(Cells, helpers.FloorCell(x, y, m, "brown"))
			}

			if helpers.IsDoor(x, y) {
				Cells = append(Cells, helpers.DoorCell(x, y, m))
			}

			if helpers.IsWindow(x, y) {
				Cells = append(Cells, helpers.WindowCell(x, y, m))
			}
		}
	}

	return Cells
}
