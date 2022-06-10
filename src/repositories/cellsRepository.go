package repositories

import (
	"fmt"
	"project0/config"
	"project0/src/models"
)

func UpdateCellUnderUser(cell models.Cell, itemCell models.ItemCell, cellType string) {
	err := config.Db.
		Model(models.Cell{}).
		Where(&models.Cell{AxisX: cell.AxisX, AxisY: cell.AxisY, MapsId: cell.MapsId}).
		Update("item_cell_id", itemCell.ID).
		Update("type", cellType).
		Error
	if err != nil {
		fmt.Println(err)
	}
}

func GetFullMap(id int) []models.Cell {
	var results []models.Cell

	err := config.Db.
		Preload("Item").
		Where(models.Cell{MapsId: id}).
		Find(&results).
		Error

	if err != nil {
		fmt.Println("Map not found!")
	}

	return results
}

func CreateCells(cells []models.Cell) bool {
	err := config.Db.
		Create(cells).
		Error

	if err != nil {
		return true
	}

	return false
}

func UpdateItemCellDestructHp(itemCell models.ItemCell, destructionHp int) {
	err := config.Db.Model(&models.ItemCell{}).
		Where(&models.ItemCell{ID: itemCell.ID}).
		Update("destruction_hp", destructionHp).
		Update("growing_time", nil).
		Update("last_growing", nil).
		Error

	if err != nil {
		panic(err)
	}
}

func GetCellsUserMap(mapSize models.UserMap, userLocation models.Location) (resultCell []models.Cell) {
	err := config.Db.
		Preload("Teleport").
		Preload("ItemCell").
		Preload("ItemCell.Item").
		Preload("ItemCell.Item.Instruments").
		Preload("ItemCell.Item.Instruments.Good").
		Preload("ItemCell.Item.Instruments.Result").
		Preload("ItemCell.Item.Instruments.NextStageItem").
		Preload("ItemCell.Item.Instruments.GrowingItem").
		Preload("ItemCell.ContainedItem").
		Where(models.Cell{MapsId: *userLocation.MapsId}).
		Where("axis_x >= " + fmt.Sprintf("%d", mapSize.LeftIndent)).
		Where("axis_x <= " + fmt.Sprintf("%d", mapSize.RightIndent)).
		Where("axis_y >= " + fmt.Sprintf("%d", mapSize.DownIndent)).
		Where("axis_y <= " + fmt.Sprintf("%d", mapSize.UpperIndent)).
		Order("axis_x ASC").
		Order("axis_y ASC").
		Find(&resultCell).
		Error

	if err != nil {
		fmt.Printf("GetCellUserMap error: %s", err)
	}

	return resultCell
}
