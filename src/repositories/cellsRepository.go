package repositories

import (
	"fmt"
	"project0/config"
	"project0/src/models"
	"time"
)

func UpdateCellUnderUser(cell models.Cell, userItem models.UserItem, count int, cellType string) {
	var nextStateTime time.Time
	if userItem.Item.GrowingUpTime != nil {
		nextStateTime = time.Now().Add(time.Duration(*userItem.Item.GrowingUpTime) * time.Minute)
	}
	err := config.Db.
		Model(models.Cell{}).
		Where(&models.Cell{AxisX: cell.AxisX, AxisY: cell.AxisY, MapsId: cell.MapsId}).
		Update("item_id", userItem.ItemId).
		Update("item_count", count).
		Update("type", cellType).
		Update("destruction_hp", nil).
		Update("next_state_time", nextStateTime).
		Update("last_growing", nil).
		Update("prev_item_id", nil).
		Update("prev_item_count", nil).
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

func UpdateCellDestructHp(cell models.Cell, destructionHp int) {
	err := config.Db.Model(&models.Cell{}).
		Where(&models.Cell{ID: cell.ID}).
		Update("destruction_hp", destructionHp).
		Update("next_state_time", nil).
		Update("last_growing", nil).
		Error

	if err != nil {
		panic(err)
	}
}

func GetCellsUserMap(mapSize models.UserMap, userLocation models.Location) (resultCell []models.Cell) {
	err := config.Db.
		Preload("Item").
		Preload("Teleport").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.Result").
		Preload("Item.Instruments.NextStageItem").
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
