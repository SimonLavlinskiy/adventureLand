package actionsController

import (
	"errors"
	"fmt"
	"project0/src/controllers/itemCellController"
	"project0/src/controllers/userItemController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"strings"
)

func UpdateModelsAfterHandOrInstrumentActions(user models.User, cell models.Cell, charData []string) string {
	var newUserItem models.UserItem
	var instrument *models.Instrument
	var err error

	if charData[0] == "item" || charData[0] == "fist" {
		instrumentId := uint(helpers.ToInt(charData[2]))

		err, instrument = checkItemsOnNeededInstrument(cell, instrumentId)
		if err != nil {
			return "Предмет не поддается под таким инструментом"
		}
	}

	if charData[0] == "takeContain" {
		newUserItem = repositories.GetOrCreateUserItem(user, *cell.ItemCell.ContainedItem)
	} else if charData[0] == "hand" {
		newUserItem = repositories.GetOrCreateUserItem(user, *cell.ItemCell.Item)
	} else if instrument.ResultId != nil {
		newUserItem = repositories.GetOrCreateUserItem(user, *instrument.Result.Item)
		fmt.Println(3)
	}

	switch true {
	case isUserHasMaxCountItem(newUserItem):
		return "У тебя уже есть такой!"

	case !canUserPayItem(user, cell) && cell.NeedPay:
		return "Не хватает деняк!"

	case strings.Contains(charData[0], "hand"):
		UpdateUserItem(cell, user, newUserItem)
		return UserGetItemWithHand(cell)

	case strings.Contains(charData[0], "takeContain"):
		UpdateUserItem(cell, user, newUserItem)
		return UserGetContainedItem(*cell.ItemCell)

	case strings.Contains(charData[0], "fist"):
		return itemCellController.DestructItem(*cell.ItemCell, user, *instrument)

	case strings.Contains(charData[0], "item"):
		return UserGetItemWithInstrument(cell, user, *instrument)

	}

	return "Нельзя взять!"
}

func checkItemsOnNeededInstrument(cell models.Cell, instrumentId uint) (error, *models.Instrument) {
	for _, instrument := range cell.ItemCell.Item.Instruments {
		if instrument.Good.ID == instrumentId {
			res := models.Instrument{ID: instrument.ID}.GetInstrument()
			return nil, &res
		}
	}

	return errors.New("user has not instrument"), nil
}

func UserGetItemWithInstrument(cell models.Cell, user models.User, instrument models.Instrument) (result string) {
	var instrumentMsg string

	err, userInstrument := user.CheckUserHasInstrument(instrument)
	if err != nil {
		return "Нет инструмента в руках"
	}

	switch instrument.Type {
	case "destruction":
		result = itemCellController.DestructItem(*cell.ItemCell, user, instrument)

		if instrumentMsg, err = userItemController.SubCountUsingFromInstrument(user, userInstrument); err != nil {
			result = result + instrumentMsg
		}

	case "growing":
		if result, err = itemCellController.GrowingItem(*cell.ItemCell, user, instrument); err == nil {
			if instrumentMsg, err = userItemController.SubCountUsingFromInstrument(user, userInstrument); err != nil {
				result = fmt.Sprintf("%s\n%s", result, instrumentMsg)
			}
		}

	case "swap":
		result = itemCellController.SwapItem(user, cell, instrument, userInstrument)

	}

	return result
}

func UserGetItemWithHand(cell models.Cell) string {
	var textCountLeft string

	if *cell.Type != "swap" && *cell.ItemCell.ItemCount > 1 {
		*cell.ItemCell.ItemCount = *cell.ItemCell.ItemCount - 1
		cell.ItemCell.UpdateItemCell()
		textCountLeft = fmt.Sprintf("(Осталось лежать еще %d)", *cell.ItemCell.ItemCount)
	} else if cell.ItemCell.ContainedItem != nil {
		*cell.ItemCell = cell.ItemCell.BecomeToContainedItem()
		cell.ItemCell.SpecialUpdateItemCell()
	} else {
		*cell.ItemCell.ItemCount = *cell.ItemCell.ItemCount - 1
		cell.ItemCell.DestructionHp = nil
		cell.ItemCell.BrokenTime = nil
		cell.ItemCell.GrowingTime = nil
		cell.ItemCell.SpecialUpdateItemCell()
	}

	return fmt.Sprintf("Ты получил %s 1 шт. %s", cell.ItemCell.Item.View, textCountLeft)
}

func UserGetContainedItem(itemCell models.ItemCell) (result string) {
	*itemCell.ContainedItemCount = *itemCell.ContainedItemCount - 1
	if *itemCell.ContainedItemCount >= 1 {
		itemCell.UpdateItemCell()
		result = fmt.Sprintf("(Осталось еще %d%s)", *itemCell.ContainedItemCount, itemCell.ContainedItem.View)
	} else if *itemCell.ContainedItemCount <= 0 {
		itemCell.ContainedItemId = nil
		itemCell.ContainedItemCount = nil
		itemCell.ContainedItemBrokenTime = nil
		itemCell.SpecialUpdateItemCell()
	}

	result = fmt.Sprintf("Ты получил %s 1 шт. %s", itemCell.ContainedItem.View, result)

	return result
}

func UpdateUserItem(cell models.Cell, user models.User, userItem models.UserItem) models.UserItem {
	sumCountItem := *userItem.Count + 1

	updateUserMoney := *user.Money

	if cell.NeedPay {
		updateUserMoney -= *cell.ItemCell.Item.Cost
	}

	var countUseLeft = userItem.Item.CountUse

	if userItem.CountUseLeft != nil {
		countUseLeft = userItem.CountUseLeft
	}

	if *userItem.Count == 0 && userItem.Item.CountUse != nil {
		*countUseLeft = *userItem.Item.CountUse
	}

	models.User{ID: user.ID}.UpdateUserItem(models.UserItem{ID: userItem.ID, Count: &sumCountItem, CountUseLeft: countUseLeft})
	repositories.UpdateUser(models.User{TgId: user.TgId, Money: &updateUserMoney})

	return userItem
}

func canUserPayItem(user models.User, cell models.Cell) bool {
	return cell.ItemCell.Item.Cost == nil || *user.Money >= *cell.ItemCell.Item.Cost
}

func isUserHasMaxCountItem(item models.UserItem) bool {
	if item.Item.MaxCountUserHas == nil || *item.Count < *item.Item.MaxCountUserHas {
		return false
	}
	return true
}
