package itemCellController

import (
	"fmt"
	"project0/src/controllers/resultController"
	"project0/src/models"
	"project0/src/repositories"
)

func DestructItem(itemCell models.ItemCell, user models.User, instrument models.Instrument) (result string) {

	var itemDestructionHp int
	if itemCell.DestructionHp == nil {
		itemDestructionHp = *itemCell.Item.DestructionHp
	} else {
		itemDestructionHp = *itemCell.DestructionHp
	}

	itemDestructionHp -= *instrument.Good.Destruction

	if !isItemCrushed(itemCell, itemDestructionHp) {
		repositories.UpdateItemCellDestructHp(itemCell, itemDestructionHp)

		result = fmt.Sprintf("Попробуй еще.. (%s)", itemHpLeft(itemCell, instrument))
		return result
	}

	if instrument.Result == nil {
		return "А тут ничего нет... 🤔 Хм..."
	}

	resultController.UserGetResult(user, *instrument.Result)
	UpdateItemCellAfterDestruction(itemCell, instrument)

	result = fmt.Sprintf("\nТы получил %s %d шт. %s", instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)
	return result
}

func isItemCrushed(itemCell models.ItemCell, ItemHp int) bool {
	if itemCell.Item.DestructionHp != nil && ItemHp <= 0 {
		return true
	} else {
		return false
	}
}

func itemHpLeft(itemCell models.ItemCell, instrument models.Instrument) (result string) {
	maxCountHit := int(float64(*itemCell.Item.DestructionHp / *instrument.Good.Destruction))
	var countHitLeft int

	if itemCell.DestructionHp != nil {
		countHitLeft = int(float64(*itemCell.DestructionHp / *instrument.Good.Destruction))
	} else {
		countHitLeft = int(float64(*itemCell.Item.DestructionHp / *instrument.Good.Destruction))
	}

	for i := 1; i <= maxCountHit; i++ {
		if i < countHitLeft {
			result += instrument.Good.View
		} else {
			result += "✔️"
		}
	}
	return result
}
