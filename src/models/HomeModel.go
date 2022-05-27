package models

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
