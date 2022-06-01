package itemController

import (
	"errors"
	"fmt"
	"project0/src/controllers/resultController"
	"project0/src/controllers/userItemController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"time"
)

func checkItemsOnNeededInstrument(cell models.Cell, instrumentId uint) (error, *models.Instrument) {
	for _, instrument := range cell.Item.Instruments {
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
		result = DestructItem(cell, user, instrument)
		instrumentMsg, err = userItemController.UpdateUserInstrument(user, userInstrument)
		if err != nil {
			result = result + instrumentMsg
		}
	case "growing":
		result, err = GrowingItem(cell, user, instrument)
		if err == nil {
			instrumentMsg, err = userItemController.UpdateUserInstrument(user, userInstrument)
		}
		if err != nil {
			result = fmt.Sprintf("%s\n%s", result, instrumentMsg)
		}
	case "swap":
		result = swapItem(user, cell, instrument, userInstrument)
	}

	return result
}

func UserGetItemWithHand(cell models.Cell, user models.User, userGetItem models.UserItem) string {
	sumCountItem := *userGetItem.Count + 1
	updateUserMoney := *user.Money

	if cell.NeedPay {
		updateUserMoney -= *cell.Item.Cost
	}

	var countUseLeft = userGetItem.Item.CountUse

	if userGetItem.CountUseLeft != nil {
		countUseLeft = userGetItem.CountUseLeft
	}
	if *userGetItem.Count == 0 && userGetItem.Item.CountUse != nil {
		*countUseLeft = *userGetItem.Item.CountUse
	}

	models.User{ID: user.ID}.UpdateUserItem(models.UserItem{ID: userGetItem.ID, Count: &sumCountItem, CountUseLeft: countUseLeft})
	repositories.UpdateUser(models.User{TgId: user.TgId, Money: &updateUserMoney})

	var textCountLeft string
	if *cell.Type != "swap" && (*cell.ItemCount > 1 || cell.PrevItemID == nil) {
		countAfterUserGetItem := *cell.ItemCount - 1
		cell.ItemCount = &countAfterUserGetItem
		cell.UpdateCell()

		if countAfterUserGetItem != 0 {
			textCountLeft = fmt.Sprintf("(Осталось лежать еще %d)", countAfterUserGetItem)
		}
	} else if cell.PrevItemID != nil {
		cell.UpdateCellOnPrevItem()
	}

	return fmt.Sprintf("Ты получил %s 1 шт. %s", userGetItem.Item.View, textCountLeft)
}

func itemHpLeft(cell models.Cell, instrument models.Instrument) string {
	maxCountHit := int(float64(*cell.Item.DestructionHp / *instrument.Good.Destruction))
	var countHitLeft int

	if cell.DestructionHp != nil {
		countHitLeft = int(float64(*cell.DestructionHp / *instrument.Good.Destruction))
	} else {
		countHitLeft = int(float64(*cell.Item.DestructionHp / *instrument.Good.Destruction))
	}

	var result string
	for i := 1; i <= maxCountHit; i++ {
		if i < countHitLeft {
			result += instrument.Good.View
		} else {
			result += "✔️"
		}
	}
	return result
}

func DestructItem(cell models.Cell, user models.User, instrument models.Instrument) (result string) {
	var itemDestructionHp int
	if cell.DestructionHp == nil {
		itemDestructionHp = *cell.Item.DestructionHp
	} else {
		itemDestructionHp = *cell.DestructionHp
	}

	itemDestructionHp -= *instrument.Good.Destruction

	if !isItemCrushed(cell, itemDestructionHp) {
		repositories.UpdateCellDestructHp(cell, itemDestructionHp)
		result = fmt.Sprintf("Попробуй еще.. (%s)", itemHpLeft(cell, instrument))
		return result
	}

	result = "А тут ничего нет... 🤔 Хм..."

	if instrument.Result == nil {
		return result
	}

	updateUserMoney := *user.Money - *cell.Item.Cost

	resultController.UserGetResult(user, *instrument.Result)
	repositories.UpdateUser(models.User{TgId: user.TgId, Money: &updateUserMoney})
	cell.UpdateCellAfterDestruction(instrument)

	result = fmt.Sprintf("\nТы получил %s %d шт. %s", instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)
	return result
}

func isItemGrown(cell models.Cell, updateItemTime time.Time) bool {
	if cell.Item.Growing != nil && updateItemTime.Before(time.Now()) {
		return true
	} else {
		return false
	}
}

func isItemCrushed(cell models.Cell, ItemHp int) bool {
	if cell.Item.DestructionHp != nil && ItemHp <= 0 {
		return true
	} else {
		return false
	}
}

func UserGetItemUpdateModels(user models.User, cell models.Cell, charData []string) string {
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

	if instrument == nil || instrument.ResultId == nil {
		newUserItem = repositories.GetOrCreateUserItem(user, *cell.Item)
	} else {
		newUserItem = repositories.GetOrCreateUserItem(user, *instrument.Result.Item)
	}

	if isUserHasMaxCountItem(newUserItem) {
		return "У тебя уже есть такой!"
	}

	if !canUserPayItem(user, cell) && cell.NeedPay {
		return "Не хватает деняк!"
	}

	if charData[0] == "hand" {
		return UserGetItemWithHand(cell, user, newUserItem)
	} else if charData[0] == "fist" {
		return DestructItem(cell, user, *instrument)
	} else if len(cell.Item.Instruments) != 0 {
		return UserGetItemWithInstrument(cell, user, *instrument)
	}

	return "Нельзя взять!"
}

func canUserPayItem(user models.User, cell models.Cell) bool {
	return cell.Item.Cost == nil || *user.Money >= *cell.Item.Cost
}

func isUserHasMaxCountItem(item models.UserItem) bool {
	if item.Item.MaxCountUserHas == nil || *item.Count < *item.Item.MaxCountUserHas {
		return false
	}
	return true
}

func ViewItemInfo(cell models.Cell) string {
	var itemInfo string
	var dressType string

	if cell.Item == nil {
		return "Ячейка пустая"
	}

	if cell.Item.DressType != nil {
		switch *cell.Item.DressType {
		case "hand":
			dressType = "(Для рук)"
		case "head":
			dressType = "(Головной убор)"
		case "body":
			dressType = "(Верхняя одежда)"
		case "shoes":
			dressType = "(Обувь)"
		case "foot":
			dressType = "(Штанихи)"
		}
	}

	itemInfo = fmt.Sprintf("%s *%s* _%s_\n", cell.Item.View, cell.Item.Name, dressType)
	if cell.ItemCount != nil {
		itemInfo += fmt.Sprintf("*Кол-во*: _%d шт._\n", *cell.ItemCount)
	}
	itemInfo = itemInfo + fmt.Sprintf("*Описание*: `%s`\n", *cell.Item.Description)

	if cell.Item.Healing != nil && *cell.Item.Healing != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Здоровье*: `+%d♥️`\n", *cell.Item.Healing)
	}
	if cell.Item.Damage != nil && *cell.Item.Damage != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Атака*: `+%d`💥️\n", *cell.Item.Damage)
	}
	if cell.Item.Satiety != nil && *cell.Item.Satiety != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Сытость*: `+%d`\U0001F9C3️\n", *cell.Item.Satiety)
	}
	if cell.Item.Cost != nil && *cell.Item.Cost != 0 && cell.NeedPay {
		itemInfo = itemInfo + fmt.Sprintf("*Стоимость*: `%d`💰\n", *cell.Item.Cost)
	}
	if cell.Item.Destruction != nil && *cell.Item.Destruction != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Сила*: `%d %s`\n", *cell.Item.Destruction, cell.Item.View)
	}
	if cell.DestructionHp != nil && cell.Item.DestructionHp != nil && *cell.Item.DestructionHp != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Прочность*: `%d`\n", *cell.DestructionHp)
	} else if cell.Item.DestructionHp != nil && *cell.Item.DestructionHp != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Прочность*: `%d`\n", *cell.Item.DestructionHp)
	}
	if cell.Item.Growing != nil && cell.NextStateTime != nil {
		t := cell.NextStateTime.Sub(time.Now())
		h := t.Truncate(time.Hour).Hours()
		m := t.Truncate(time.Minute).Minutes() - t.Truncate(time.Hour).Minutes()
		itemInfo = itemInfo + fmt.Sprintf("*\U0001F973 Вырастет через*: _%vч %vм_\n", h, m)
	} else if cell.Item.Growing != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Время роста*: `%d мин.`\n", *cell.Item.Growing)
	}
	if cell.Item.IntervalGrowing != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Интервал ускорения роста*: `раз в %d мин.`\n", *cell.Item.IntervalGrowing)
	}
	if cell.LastGrowing != nil {
		t := time.Now().Sub(*cell.LastGrowing)
		m := t.Truncate(time.Minute).Minutes()
		itemInfo = itemInfo + fmt.Sprintf("*Последнее ускорение:* %vм назад\n", m)
	}
	if len(cell.Item.Instruments) != 0 {
		var itemsInstrument string
		for _, i := range cell.Item.Instruments {
			if i.GoodId != nil {
				itemsInstrument = itemsInstrument + fmt.Sprintf("%s - `%s`\n", i.Good.View, i.Good.Name)
			}
		}
		itemInfo = itemInfo + fmt.Sprintf("*Чем можно взаимодествовать*:\n%s", itemsInstrument)
	}

	return itemInfo
}

func swapItem(user models.User, cell models.Cell, instrument models.Instrument, userInstrument models.Item) string {
	resultController.UserGetResult(user, *instrument.Result)
	result := fmt.Sprintf("Ты получил %s %d шт. %s", instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)

	updateUserMoney := *user.Money - *cell.Item.Cost

	repositories.UpdateUser(models.User{TgId: user.TgId, Money: &updateUserMoney})

	instrumentMsg, err := userItemController.UpdateUserInstrument(user, userInstrument)
	if err != nil {
		result += instrumentMsg
	}
	cell.UpdateCellAfterDestruction(instrument)

	return result
}
