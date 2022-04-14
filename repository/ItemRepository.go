package repository

import (
	"errors"
	"fmt"
	"project0/config"
	"time"
)

type Item struct {
	ID                uint         `gorm:"primaryKey"`
	Name              string       `gorm:"embedded"`
	Description       *string      `gorm:"embedded"`
	View              string       `gorm:"embedded"`
	Type              string       `gorm:"embedded"`
	Cost              *int         `gorm:"embedded"`
	Healing           *int         `gorm:"embedded"`
	Damage            *int         `gorm:"embedded"`
	Satiety           *int         `gorm:"embedded"`
	Destruction       *int         `gorm:"embedded"`
	DestructionHp     *int         `gorm:"embedded"`
	GrowingUpTime     *int         `gorm:"embedded"`
	Growing           *int         `gorm:"embedded"`
	IntervalGrowing   *int         `gorm:"embedded"`
	CanTake           bool         `gorm:"embedded"`
	CanStep           bool         `gorm:"embedded"`
	CanDestructByHand bool         `gorm:"embedded"`
	Instruments       []Instrument `gorm:"many2many:instrument_item;"`
	DressType         *string      `gorm:"embedded"`
	IsBackpack        bool         `gorm:"embedded"`
	IsInventory       bool         `gorm:"embedded"`
	MaxCountUserHas   *int         `gorm:"embedded"`
	CountUse          *int         `gorm:"embedded"`
}

type InstrumentItem struct {
	ItemID       int `gorm:"primaryKey"`
	InstrumentID int `gorm:"primaryKey"`
}

func (i Item) GetItemEndTime() time.Time {
	return time.Now().Add(time.Duration(*i.Growing) * time.Minute)
}

func UserGetItem(user User, LocationStruct Location, char []string) string {
	resultCell := Cell{MapsId: *LocationStruct.MapsId, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY}
	resultCell = resultCell.GetCell()

	if resultCell.ItemID == nil {
		return "Не получилось..."
	}

	return UserGetItemUpdateModels(user, resultCell, char[0])
}

func checkItemsOnNeededInstrument(cell Cell, msgInstrumentView string) (error, *Instrument) {
	for _, instrument := range cell.Item.Instruments {
		if instrument.Good.View == msgInstrumentView {
			res := Instrument{ID: instrument.ID}.GetInstrument()
			return nil, &res
		}
	}
	if msgInstrumentView == "👋" && cell.Item.CanTake {
		return nil, nil
	}

	return errors.New("user has not instrument"), nil
}

func UserGetItemWithInstrument(cell Cell, user User, instrument Instrument) string {
	var result string
	var instrumentMsg string
	var err error

	err, userInstrument := user.CheckUserHasInstrument(instrument)
	if err != nil {
		return "Нет инструмента в руках"
	}

	switch instrument.Type {
	case "destruction":
		itemMsg := DestructItem(cell, user, instrument)
		instrumentMsg, err = UpdateUserInstrument(user, userInstrument)
		if err != nil {
			result = itemMsg + instrumentMsg
		}
		result = itemMsg
	case "growing":
		itemMsg, err := GrowingItem(cell, user, instrument)
		if err == nil {
			if instrumentMsg, err = UpdateUserInstrument(user, userInstrument); err != nil {
				result = itemMsg + instrumentMsg
			}
		}
		result = itemMsg
	case "swap":
		result = swapItem(user, cell, instrument, userInstrument)
	}

	return result
}

func UserGetItemWithHand(cell Cell, user User, userGetItem UserItem) string {
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

	User{ID: user.ID}.UpdateUserItem(UserItem{ID: userGetItem.ID, Count: &sumCountItem, CountUseLeft: countUseLeft})
	User{TgId: user.TgId, Money: &updateUserMoney}.UpdateUser()

	textCountLeft := ""
	if *cell.Type != "swap" && (*cell.ItemCount > 1 || cell.PrevItemID == nil) {
		countAfterUserGetItem := *cell.ItemCount - 1
		Cell{ItemCount: &countAfterUserGetItem}.UpdateCell(cell.ID)
		if countAfterUserGetItem != 0 {
			textCountLeft = fmt.Sprintf("(Осталось лежать еще %d)", countAfterUserGetItem)
		}
	} else if cell.PrevItemID != nil {
		cell.UpdateCellOnPrevItem()
	}

	return fmt.Sprintf("Ты получил %s 1 шт. %s", userGetItem.Item.View, textCountLeft)
}

func itemHpLeft(cell Cell, instrument Instrument) string {
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

func GrowingItem(cell Cell, user User, instrument Instrument) (string, error) {
	var updateItemTime = time.Now()

	if cell.LastGrowing != nil && time.Now().Before(cell.LastGrowing.Add(time.Duration(*cell.Item.IntervalGrowing)*time.Minute)) {
		nextTimeGrowing := cell.LastGrowing.Add(time.Duration(*cell.Item.IntervalGrowing) * time.Minute).Format("15:04:05 02.01.06")

		return fmt.Sprintf("Ты уже использовал %s\nМожно будет повторить %s!", instrument.Good.View, nextTimeGrowing), errors.New("user can not growing")
	}

	if cell.NextStateTime == nil && cell.Item.Growing != nil {
		updateItemTime = updateItemTime.Add(time.Duration(*cell.Item.Growing) * time.Minute)
	} else {
		updateItemTime = *cell.NextStateTime
	}
	updateItemTime = updateItemTime.Add(-time.Duration(*instrument.Good.GrowingUpTime) * time.Minute)

	if isItemGrown(cell, updateItemTime) {
		var result string
		updateUserMoney := *user.Money - *cell.Item.Cost

		if instrument.Result != nil {
			user.UserGetResult(*instrument.Result)
			result = fmt.Sprintf("\nТы получил %s %d шт. %s", instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)
		}

		User{TgId: user.TgId, Money: &updateUserMoney}.UpdateUser()

		cell.UpdateCellAfterGrowing(instrument)

		return fmt.Sprintf("Оно выросло!%s", result), nil

	} else {
		t := time.Now()
		Cell{
			ID:            cell.ID,
			NextStateTime: &updateItemTime,
			LastGrowing:   &t,
		}.UpdateCell(cell.ID)
		return "\U0001F973 Вырастет " + updateItemTime.Format("15:04:05 02.01.06") + "!", nil

	}
}

func DestructItem(cell Cell, user User, instrument Instrument) string {
	var ItemDestructionHp int
	if cell.DestructionHp == nil {
		ItemDestructionHp = *cell.Item.DestructionHp
	} else {
		ItemDestructionHp = *cell.DestructionHp
	}

	ItemDestructionHp = ItemDestructionHp - *instrument.Good.Destruction

	if isItemCrushed(cell, ItemDestructionHp) {
		var result string
		if instrument.Result != nil {
			user.UserGetResult(*instrument.Result)
			result = fmt.Sprintf("\nТы получил %s %d шт. %s", instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)
		} else {
			result = "что то не так"
		}
		updateUserMoney := *user.Money - *cell.Item.Cost

		User{TgId: user.TgId, Money: &updateUserMoney}.UpdateUser()

		cell.UpdateCellAfterDestruction(instrument)

		return result
	} else {
		err := config.Db.
			Where(&Cell{ID: cell.ID}).
			Updates(Cell{ID: cell.ID, DestructionHp: &ItemDestructionHp}).
			Update("next_state_time", nil).
			Update("last_growing", nil).
			Error
		if err != nil {
			panic(err)
		}

		return "Попробуй еще.. (" + itemHpLeft(cell, instrument) + ")"
	}
}

func isItemGrown(cell Cell, updateItemTime time.Time) bool {
	if cell.Item.Growing != nil && updateItemTime.Before(time.Now()) {
		return true
	} else {
		return false
	}
}

func isItemCrushed(cell Cell, ItemHp int) bool {
	if cell.Item.DestructionHp != nil && ItemHp <= 0 {
		return true
	} else {
		return false
	}
}

func UserGetItemUpdateModels(user User, cell Cell, instrumentView string) string {
	var userGetItem UserItem

	err, instrument := checkItemsOnNeededInstrument(cell, instrumentView)
	if err != nil {
		return "Предмет не поддается под таким инструментом"
	}

	if instrument == nil || instrument.ResultId == nil {
		userGetItem = GetOrCreateUserItem(user, *cell.Item)
	} else {
		userGetItem = GetOrCreateUserItem(user, *instrument.Result.Item)
	}

	if isUserHasMaxCountItem(userGetItem) {
		return "У тебя уже есть такой!"
	}

	if !canUserPayItem(user, cell) && cell.NeedPay {
		return "Не хватает деняк!"
	}

	if instrumentView == "👋" {
		return UserGetItemWithHand(cell, user, userGetItem)
	} else if instrumentView == "🤜" {
		return DestructItem(cell, user, *instrument)
	} else if len(cell.Item.Instruments) != 0 {
		return UserGetItemWithInstrument(cell, user, *instrument)
	}

	return "Нельзя взять!"

}

func canUserPayItem(user User, cell Cell) bool {
	return cell.Item.Cost == nil || *user.Money >= *cell.Item.Cost
}

func isUserHasMaxCountItem(item UserItem) bool {
	if item.Item.MaxCountUserHas == nil || *item.Count < *item.Item.MaxCountUserHas {
		return false
	}
	return true
}

func ViewItemInfo(location Location) string {
	cell := Cell{MapsId: *location.MapsId, AxisX: *location.AxisX, AxisY: *location.AxisY}
	cell = cell.GetCell()
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
	if cell.Item.Cost != nil && *cell.Item.Cost != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Стоимость*: `%d`💰\n", *cell.Item.Cost)
	}
	if cell.Item.Destruction != nil && *cell.Item.Destruction != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Сила*: `%d %s`\n", *cell.Item.Destruction, cell.Item.View)
	}
	if cell.DestructionHp != nil && *cell.Item.DestructionHp != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Прочность*: `%d`\n", *cell.DestructionHp)
	} else if cell.Item.DestructionHp != nil && *cell.Item.DestructionHp != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Прочность*: `%d`\n", *cell.Item.DestructionHp)
	}
	if cell.Item.Growing != nil && cell.NextStateTime != nil {
		itemInfo = itemInfo + fmt.Sprintf("*\U0001F973 Вырастет*: %s\n", cell.NextStateTime.Format("15:04:05 02.01.06"))
	} else if cell.Item.Growing != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Время роста*: `%d мин.`\n", *cell.Item.Growing)
	}
	if cell.Item.IntervalGrowing != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Интервал ускорения роста*: `раз в %d мин.`\n", *cell.Item.IntervalGrowing)
	}
	if cell.LastGrowing != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Последнее ускорение:* %s\n", cell.LastGrowing.Format("15:04:05"))
	}
	if len(cell.Item.Instruments) != 0 {
		var itemsInstrument string
		for _, i := range cell.Item.Instruments {
			itemsInstrument = itemsInstrument + fmt.Sprintf("%s - `%s`\n", i.Good.View, i.Good.Name)
		}
		itemInfo = itemInfo + fmt.Sprintf("*Чем можно взаимодествовать*:\n%s", itemsInstrument)
	}

	return itemInfo
}

func swapItem(user User, cell Cell, instrument Instrument, userInstrument Item) string {
	user.UserGetResult(*instrument.Result)
	result := fmt.Sprintf("Ты получил %s %d шт. %s", instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)

	updateUserMoney := *user.Money - *cell.Item.Cost

	User{TgId: user.TgId, Money: &updateUserMoney}.UpdateUser()

	instrumentMsg, err := UpdateUserInstrument(user, userInstrument)
	if err != nil {
		result += instrumentMsg
	}
	cell.UpdateCellAfterDestruction(instrument)

	return result
}
