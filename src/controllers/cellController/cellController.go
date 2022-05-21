package cellController

import (
	"errors"
	"project0/config"
	"project0/src/models"
	"project0/src/repositories"
)

func UpdateCellUnderUserWhenUserThrowItem(user models.User, userItem models.UserItem, count int, cellType string) error {
	userLocation := repositories.GetOrCreateMyLocation(user)

	cell := models.Cell{AxisX: *userLocation.AxisX, AxisY: *userLocation.AxisY, MapsId: *userLocation.MapsId}
	cell = cell.GetCell()
	if cell.ItemCount != nil && *cell.ItemCount > 0 {
		return errors.New("В этой ячейке уже есть предмет, перейди на другую ячейку")
	}

	repositories.UpdateCellUnderUser(cell, userItem, count, cellType)

	if cellType == "chat" {
		timeOut := userItem.Item.GetItemEndTime()
		chat := repositories.CreateChat(timeOut)
		config.Db.Model(models.Cell{}).
			Where(&models.Cell{AxisX: *userLocation.AxisX, AxisY: *userLocation.AxisY, MapsId: *userLocation.MapsId}).
			Update("chat_id", chat.ID)
	}

	return nil

}
