package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"time"
)

type Item struct {
	ID              uint         `gorm:"primaryKey"`
	Name            string       `gorm:"embedded"`
	Description     *string      `gorm:"embedded"`
	View            string       `gorm:"embedded"`
	Type            string       `gorm:"embedded"`
	Cost            *int         `gorm:"embedded"`
	Healing         *int         `gorm:"embedded"`
	Damage          *int         `gorm:"embedded"`
	Satiety         *int         `gorm:"embedded"`
	Destruction     *int         `gorm:"embedded"`
	DestructionHp   *int         `gorm:"embedded"`
	GrowingUpTime   *int         `gorm:"embedded"`
	Growing         *int         `gorm:"embedded"`
	IntervalGrowing *int         `gorm:"embedded"`
	CanTake         bool         `gorm:"embedded"`
	CanStep         bool         `gorm:"embedded"`
	Instruments     []Instrument `gorm:"many2many:instrument_item;"`
	DressType       *string      `gorm:"embedded"`
	IsBackpack      bool         `gorm:"embedded"`
	IsInventory     bool         `gorm:"embedded"`
	MaxCountUserHas *int         `gorm:"embedded"`
	CountUse        *int         `gorm:"embedded"`
}

type InstrumentItem struct {
	ItemID       int `gorm:"primaryKey"`
	InstrumentID int `gorm:"primaryKey"`
}

func UserGetItem(update tgbotapi.Update, LocationStruct Location, char []string) string {
	resultCell := Cell{MapsId: *LocationStruct.MapsId, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY}
	resultCell = resultCell.GetCell()

	if resultCell.ItemID == nil {
		return "Не получилось..."
	}

	return UserGetItemUpdateModels(update, resultCell, char[0])
}

func checkItemsOnNeededInstrument(cell Cell, msgInstrumentView string) (string, *Instrument) {
	for _, instrument := range cell.Item.Instruments {
		if instrument.Good.View == msgInstrumentView {
			return "Ok", &instrument
		}
	}
	if msgInstrumentView == "👋" && cell.Item.CanTake {
		return "Ok", nil
	}
	return "Not ok", nil
}

func UserGetItemWithInstrument(update tgbotapi.Update, cell Cell, user User, instrument Instrument, userGetItem UserItem) string {
	var result string
	var instrumentMsg string

	status, userInstrument := user.CheckUserHasInstrument(instrument)
	if status != "Ok" {
		return "Нет инструмента в руках"
	}

	switch instrument.Type {
	case "destruction":
		itemMsg := DestructItem(update, cell, user, userGetItem, instrument)
		_, instrumentMsg = UpdateUserInstrument(update, user, userInstrument)
		if instrumentMsg != "Ok" {
			result = itemMsg + instrumentMsg
		}
		result = itemMsg
	case "growing":
		status, itemMsg := GrowingItem(update, cell, user, userGetItem, instrument)
		if status == "Ok" {
			_, instrumentMsg = UpdateUserInstrument(update, user, userInstrument)
		}
		if instrumentMsg != "Ok" {
			result = itemMsg + instrumentMsg
		}
		result = itemMsg
	case "swap":
		result = swapItem(update, user, cell, userGetItem, instrument, userInstrument)
	}

	return result
}

func UserGetItemWithHand(update tgbotapi.Update, cell Cell, user User, userGetItem UserItem) string {
	sumCountItem := *userGetItem.Count + 1
	updateUserMoney := *user.Money - *cell.Item.Cost
	var countUseLeft = userGetItem.Item.CountUse

	if userGetItem.CountUseLeft != nil {
		countUseLeft = userGetItem.CountUseLeft
	}
	if *userGetItem.Count == 0 && userGetItem.Item.CountUse != nil {
		*countUseLeft = *userGetItem.Item.CountUse
	}

	UserItem{ID: userGetItem.ID, Count: &sumCountItem, CountUseLeft: countUseLeft}.UpdateUserItem(User{ID: user.ID})
	User{Money: &updateUserMoney}.UpdateUser(update)

	var countAfterUserGetItem *int
	textCountLeft := ""
	if *cell.Type != "swap" && (*countAfterUserGetItem != 0 || cell.PrevItemID == nil) {
		*countAfterUserGetItem = *cell.ItemCount - 1
		Cell{ItemCount: countAfterUserGetItem}.UpdateCell(cell.ID)
		textCountLeft = fmt.Sprintf("(Осталось лежать еще %d)", countAfterUserGetItem)
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

func GrowingItem(update tgbotapi.Update, cell Cell, user User, userGetItem UserItem, instrument Instrument) (string, string) {
	var updateItemTime = time.Now()

	if cell.LastGrowing != nil && time.Now().Before(cell.LastGrowing.Add(time.Duration(*cell.Item.IntervalGrowing)*time.Minute)) {
		nextTimeGrowing := cell.LastGrowing.Add(time.Duration(*cell.Item.IntervalGrowing) * time.Minute).Format("15:04:05 02.01.06")

		return "Not ok", fmt.Sprintf("Ты уже использовал %s\nМожно будет повторить %s!", instrument.Good.View, nextTimeGrowing)
	}

	if cell.NextStateTime == nil && cell.Item.Growing != nil {
		updateItemTime = updateItemTime.Add(time.Duration(*cell.Item.Growing) * time.Minute)
	} else {
		updateItemTime = *cell.NextStateTime
	}
	updateItemTime = updateItemTime.Add(-time.Duration(*instrument.Good.GrowingUpTime) * time.Minute)

	if isItemGrowed(cell, updateItemTime) {
		var result string
		updateUserMoney := *user.Money - *cell.Item.Cost

		if instrument.CountResultItem != nil {
			*userGetItem.Count = *userGetItem.Count + *instrument.CountResultItem
			result = fmt.Sprintf("\nТы получил %s %d шт. %s", instrument.ItemsResult.View, *instrument.CountResultItem, instrument.ItemsResult.Name)
		}

		User{Money: &updateUserMoney}.UpdateUser(update)
		UserItem{
			ID:           userGetItem.ID,
			Count:        userGetItem.Count,
			CountUseLeft: userGetItem.Item.CountUse,
		}.UpdateUserItem(User{ID: user.ID})

		cell.UpdateCellAfterGrowing(instrument)

		return "Ok", fmt.Sprintf("Оно выросло!%s", result)

	} else {
		t := time.Now()
		Cell{
			ID:            cell.ID,
			NextStateTime: &updateItemTime,
			LastGrowing:   &t,
		}.UpdateCell(cell.ID)
		return "Ok", "Вырастет " + updateItemTime.Format("15:04:05 02.01.06") + "!"

	}
}

func DestructItem(update tgbotapi.Update, cellule Cell, user User, userGetItem UserItem, instrument Instrument) string {
	var ItemDestructionHp int
	if cellule.DestructionHp == nil {
		ItemDestructionHp = *cellule.Item.DestructionHp
	} else {
		ItemDestructionHp = *cellule.DestructionHp
	}

	ItemDestructionHp = ItemDestructionHp - *instrument.Good.Destruction

	if isItemCrushed(cellule, ItemDestructionHp) {
		var result string
		if instrument.CountResultItem != nil {
			*userGetItem.Count = *userGetItem.Count + *instrument.CountResultItem
			result = fmt.Sprintf("Ты получил %s %d шт. %s", instrument.ItemsResult.View, *instrument.CountResultItem, instrument.ItemsResult.Name)
		} else {
			result = "что то не так"
		}
		updateUserMoney := *user.Money - *cellule.Item.Cost

		User{Money: &updateUserMoney}.UpdateUser(update)
		UserItem{
			ID:           userGetItem.ID,
			Count:        userGetItem.Count,
			CountUseLeft: userGetItem.Item.CountUse,
		}.UpdateUserItem(User{ID: user.ID})

		cellule.UpdateCellAfterDestruction(instrument)

		return result
	} else {
		err := config.Db.
			Where(&Cell{ID: cellule.ID}).
			Updates(Cell{ID: cellule.ID, DestructionHp: &ItemDestructionHp}).
			Update("next_state_time", nil).
			Update("last_growing", nil).
			Error
		if err != nil {
			panic(err)
		}

		return "Попробуй еще.. (" + itemHpLeft(cellule, instrument) + ")"
	}
}

func isItemGrowed(cellule Cell, updateItemTime time.Time) bool {
	if cellule.Item.Growing != nil && updateItemTime.Before(time.Now()) {
		return true
	} else {
		return false
	}
}

func isItemCrushed(cellule Cell, ItemHp int) bool {
	if cellule.Item.DestructionHp != nil && ItemHp <= 0 {
		return true
	} else {
		return false
	}
}

func UserGetItemUpdateModels(update tgbotapi.Update, cellule Cell, instrumentView string) string {
	userTgid := GetUserTgId(update)
	user := GetUser(User{TgId: userTgid})

	var userGetItem UserItem

	status, instrument := checkItemsOnNeededInstrument(cellule, instrumentView)
	if status != "Ok" {
		return "Предмет не поддается под таким инструментом"
	}

	if instrument == nil || instrument.ItemsResultId == nil {
		userGetItem = GetOrCreateUserItem(update, *cellule.Item)
	} else {
		userGetItem = GetOrCreateUserItem(update, *instrument.ItemsResult)
	}

	if isUserHasMaxCountItem(userGetItem) {
		return "У тебя уже есть такой!"
	}

	if !canUserPayItem(user, cellule) {
		return "Не хватает деняк!"
	}

	if instrumentView == "👋" {
		return UserGetItemWithHand(update, cellule, user, userGetItem)
	} else if instrumentView != "👋" && len(cellule.Item.Instruments) != 0 {
		return UserGetItemWithInstrument(update, cellule, user, *instrument, userGetItem)
	}

	return "Нельзя взять!"

}

func canUserPayItem(user User, cellule Cell) bool {
	return cellule.Item.Cost == nil || *user.Money >= *cellule.Item.Cost
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
		itemInfo = fmt.Sprintf("_%d шт._\n", *cell.ItemCount)
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
		itemInfo = itemInfo + fmt.Sprintf("*Вырастет*: %s\n", cell.NextStateTime.Format("15:04:05 02.01.06"))
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

func swapItem(update tgbotapi.Update, user User, cell Cell, userGetItem UserItem, instrument Instrument, userInstrument Item) string {
	*userGetItem.Count = *userGetItem.Count + *instrument.CountResultItem
	result := fmt.Sprintf("Ты получил %s %d шт. %s", instrument.ItemsResult.View, *instrument.CountResultItem, instrument.ItemsResult.Name)

	updateUserMoney := *user.Money - *cell.Item.Cost

	User{Money: &updateUserMoney}.UpdateUser(update)
	UserItem{
		ID:           userGetItem.ID,
		Count:        userGetItem.Count,
		CountUseLeft: userGetItem.Item.CountUse,
	}.UpdateUserItem(User{ID: user.ID})

	_, instrumentMsg := UpdateUserInstrument(update, user, userInstrument)
	if instrumentMsg != "Ok" {
		result += instrumentMsg
	}
	cell.UpdateCellAfterDestruction(instrument)

	return result
}
