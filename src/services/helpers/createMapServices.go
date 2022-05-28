package helpers

import (
	"encoding/csv"
	"fmt"
	v "github.com/spf13/viper"
	"log"
	"math/rand"
	"os"
	"project0/config"
	r "project0/src/models"
)

func CreateMap() {
	cellType := "cell"
	cell := r.Cell{
		MapsId:  54,
		View:    v.GetString("elem.green"),
		CanStep: true,
		Type:    &cellType,
		NeedPay: false,
	}

	for x := 0; x <= 7; x++ {
		var cells []r.Cell

		for y := 0; y <= 7; y++ {
			cell.AxisX = x
			cell.AxisY = y
			cells = append(cells, cell)
		}

		err := config.Db.Create(cells).Error

		if err != nil {
			fmt.Println("не получилось ", x)
		}
	}
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func UpdateMap() {
	csvCells := readCsvFile("./files/minimap.csv")
	mapId := 54
	var csvCellsReverted [][]string
	for i := len(csvCells) - 1; i >= 0; i-- {
		csvCellsReverted = append(csvCellsReverted, csvCells[i])
	}

	for y, row := range csvCellsReverted {
		var updateCells []r.Cell

		for x, c := range row {
			cell := getCellInterface(c, mapId)
			cell.AxisX = x
			cell.AxisY = y
			updateCells = append(updateCells, cell)
		}

		for _, upd := range updateCells {
			config.Db.Where(r.Cell{MapsId: mapId, AxisY: upd.AxisY, AxisX: upd.AxisX}).Updates(&upd)
		}

		fmt.Printf("row %d updated!\n", y)
	}
}

func getCellInterface(cell string, mapId int) r.Cell {
	cellTypeCell := "cell"
	cellTypeItem := "item"
	cellTypeSwap := "swap"

	cellItemWater := 22
	cellItemTree := 12
	cellItemPalm := 41
	cellItemElka := 43
	cellItemStone := 9

	countItem := rand.Intn(100)
	cellItemsId := []int{12, 39, 40}

	BlueCell := r.Cell{
		MapsId:  mapId,
		View:    "\U0001F7E6",
		CanStep: false,
		Type:    &cellTypeSwap,
		ItemID:  &cellItemWater,
		NeedPay: false,
	}

	GreenCell := r.Cell{
		MapsId:  mapId,
		View:    "\U0001F7E9",
		CanStep: true,
		ItemID:  nil,
		Type:    &cellTypeCell,
		NeedPay: false,
	}

	YellowCell := r.Cell{
		MapsId:  mapId,
		View:    "\U0001F7E8",
		CanStep: true,
		ItemID:  nil,
		Type:    &cellTypeCell,
		NeedPay: false,
	}

	WhiteCell := r.Cell{
		MapsId:  mapId,
		View:    "⬜️",
		CanStep: true,
		ItemID:  nil,
		Type:    &cellTypeCell,
		NeedPay: false,
	}

	BrownCell := r.Cell{
		MapsId:  mapId,
		View:    "\U0001F7EB",
		CanStep: true,
		ItemID:  nil,
		Type:    &cellTypeCell,
		NeedPay: false,
	}

	TreeCell := r.Cell{
		MapsId:    mapId,
		View:      "🌳",
		CanStep:   true,
		Type:      &cellTypeItem,
		ItemID:    &cellItemTree,
		ItemCount: &countItem,
		NeedPay:   false,
	}

	PalmaCell := r.Cell{
		MapsId:    mapId,
		View:      "🌴",
		CanStep:   true,
		Type:      &cellTypeItem,
		ItemCount: &countItem,
		NeedPay:   false,
		ItemID:    &cellItemPalm,
	}

	ElkaCell := r.Cell{
		MapsId:    mapId,
		View:      "🌲",
		CanStep:   true,
		ItemID:    &cellItemElka,
		Type:      &cellTypeItem,
		ItemCount: &countItem,
		NeedPay:   false,
	}

	FlowerCell := r.Cell{
		MapsId:  mapId,
		View:    "🌻",
		CanStep: true,
		ItemID:  nil,
		Type:    &cellTypeCell,
		NeedPay: false,
	}

	StoneCell := r.Cell{
		MapsId:    mapId,
		View:      "\U0001FAA8",
		CanStep:   true,
		ItemID:    &cellItemStone,
		Type:      &cellTypeItem,
		ItemCount: &countItem,
		NeedPay:   false,
	}

	switch cell {
	case "0":
		return BlueCell
	case "1":
		return GreenCell
	case "2":
		return YellowCell
	case "3":
		return WhiteCell
	case "4":
		return BrownCell
	case "5":
		z := rand.Intn(3)
		TreeCell.ItemID = &cellItemsId[z]
		return TreeCell
	case "51":
		TreeCell.ItemID = &cellItemTree
		return TreeCell
	case "6":
		return PalmaCell
	case "7":
		return ElkaCell
	case "8":
		return FlowerCell
	case "9":
		return StoneCell
	}

	return r.Cell{}
}
