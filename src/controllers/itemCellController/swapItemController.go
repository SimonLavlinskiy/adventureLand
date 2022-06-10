package itemCellController

import (
	"fmt"
	"project0/src/controllers/resultController"
	"project0/src/controllers/userItemController"
	"project0/src/models"
	"project0/src/repositories"
)

func SwapItem(user models.User, cell models.Cell, instrument models.Instrument, userInstrument models.Item) string {
	resultController.UserGetResult(user, *instrument.Result)
	result := fmt.Sprintf("Ты получил %s %d шт. %s", instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)

	if cell.ItemCell.Item.Cost != nil {
		updateUserMoney := *user.Money - *cell.ItemCell.Item.Cost
		repositories.UpdateUser(models.User{TgId: user.TgId, Money: &updateUserMoney})
	}

	instrumentMsg, err := userItemController.SubCountUsingFromInstrument(user, userInstrument)
	if err != nil {
		result += instrumentMsg
	}
	UpdateItemCellAfterDestruction(*cell.ItemCell, instrument)

	return result
}
