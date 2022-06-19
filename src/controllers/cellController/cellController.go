package cellController

import (
	"errors"
	"project0/config"
	"project0/src/controllers/itemCellController"
	"project0/src/models"
	"project0/src/repositories"
)

func UpdateCellUnderUserWhenUserThrowItem(user models.User, userItem models.UserItem, count int, cellType string) error {
	cell, userLocation := GetCellUnderUser(user)

	if cell.ItemCellId != nil && cell.ItemCell.ItemCount != nil && *cell.ItemCell.ItemCount > 0 {
		return errors.New("В этой ячейке уже есть предмет, перейди на другую ячейку")
	}

	itemCell := itemCellController.UpdateItemCellUnderUser(cell, userItem, count)
	repositories.UpdateCellUnderUser(cell, itemCell, cellType)

	if cellType == "chat" {
		timeOut := userItem.Item.GetItemEndTime()
		chat := repositories.CreateChat(timeOut)
		config.Db.Model(models.Cell{}).
			Where(&models.Cell{AxisX: *userLocation.AxisX, AxisY: *userLocation.AxisY, MapsId: *userLocation.MapsId}).
			Update("chat_id", chat.ID)
	}

	return nil

}

func GetCellUnderUser(user models.User) (models.Cell, models.Location) {
	userLocation := repositories.GetOrCreateMyLocation(user)
	cell := models.Cell{AxisX: *userLocation.AxisX, AxisY: *userLocation.AxisY, MapsId: *userLocation.MapsId}
	cell = cell.GetCell()

	return cell, userLocation
}
