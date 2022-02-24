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
	resultCell := GetCellule(Cellule{MapsId: *LocationStruct.MapsId, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY})

	if resultCell.ItemID != nil {
		res := UserGetItemUpdateModels(update, resultCell, char[0])

		return res
	}

	return "Не получилось..."
}

func checkItemsOnNeededInstrument(instruments []Instrument, msgInstrumentView string) (string, *Instrument) {
	for _, instrument := range instruments {
		if instrument.Good.View == msgInstrumentView {
			return "Ok", &instrument
		}
	}
	if msgInstrumentView == "👋" {
		return "Ok", nil
	}
	return "Not ok", nil
}

func UserGetItemWithInstrument(update tgbotapi.Update, cellule Cellule, user User, instrument Instrument, userGetItem UserItem) string {
	var result string
	var instrumentMsg string

	status, userInstrument := CheckUserHasInstrument(user, instrument)
	if status != "Ok" {
		return "Нет инструмента в руках"
	}

	switch instrument.Type {
	case "destruction":
		itemMsg := DesctructionItem(update, cellule, user, userGetItem, instrument)
		_, instrumentMsg = UpdateUserInstrument(update, user, userInstrument)
		result = itemMsg + instrumentMsg
	case "hand":
		result = UserGetItemWithHand(update, cellule, user, userGetItem)
	case "growing":
		status, itemMsg := GrowingItem(update, cellule, user, userGetItem, instrument)
		if status == "Ok" {
			_, instrumentMsg = UpdateUserInstrument(update, user, userInstrument)
		}
		result = itemMsg + instrumentMsg
	}

	return result
}

func UserGetItemWithHand(update tgbotapi.Update, cellule Cellule, user User, userGetItem UserItem) string {
	sumCountItem := *userGetItem.Count + 1
	countAfterUserGetItem := *cellule.ItemCount - 1
	updateUserMoney := *user.Money - *cellule.Item.Cost
	var countUseLeft = userGetItem.Item.CountUse

	if userGetItem.CountUseLeft != nil {
		countUseLeft = userGetItem.CountUseLeft
	}
	if *userGetItem.Count == 0 && userGetItem.Item.CountUse != nil {
		*countUseLeft = *userGetItem.Item.CountUse
	}

	UpdateUserItem(User{ID: user.ID}, UserItem{ID: userGetItem.ID, Count: &sumCountItem, CountUseLeft: countUseLeft})
	UpdateUser(update, User{Money: &updateUserMoney})

	textCountLeft := ""
	if countAfterUserGetItem != 0 || cellule.PrevItemID == nil {
		UpdateCellule(cellule.ID, Cellule{ItemCount: &countAfterUserGetItem})
	} else {
		UpdateCellOnPrevItem(cellule)
	}

	if countAfterUserGetItem != 0 {
		textCountLeft = fmt.Sprintf("(Осталось лежать еще %s)", ToString(countAfterUserGetItem))
	}
	return fmt.Sprintf("Ты получил %s 1 шт. %s", userGetItem.Item.View, textCountLeft)
}

func itemHpLeft(cellule Cellule, instrument Instrument) string {
	maxCountHit := int(float64(*cellule.Item.DestructionHp / *instrument.Good.Destruction))
	countHitLeft := int(float64(*cellule.DestructionHp / *instrument.Good.Destruction))

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

func GrowingItem(update tgbotapi.Update, cellule Cellule, user User, userGetItem UserItem, instrument Instrument) (string, string) {
	var updateItemTime = time.Now()

	if cellule.LastGrowing != nil && time.Now().Before(cellule.LastGrowing.Add(time.Duration(*cellule.Item.IntervalGrowing)*time.Minute)) {
		return "Not ok", "Ты уже использовал " + instrument.Good.View + "\nМожно будет повторить " + cellule.LastGrowing.Add(time.Duration(*cellule.Item.IntervalGrowing)*time.Minute).Format("15:04:05 02.01.06") + "!"
	}

	if cellule.NextStateTime == nil && cellule.Item.Growing != nil {
		updateItemTime = updateItemTime.Add(time.Duration(*cellule.Item.Growing) * time.Minute)
	} else {
		updateItemTime = *cellule.NextStateTime
	}
	updateItemTime = updateItemTime.Add(-time.Duration(*instrument.Good.GrowingUpTime) * time.Minute)

	if isItemGrowed(cellule, updateItemTime) {
		var result string
		updateUserMoney := *user.Money - *cellule.Item.Cost

		if instrument.CountResultItem != nil {
			*userGetItem.Count = *userGetItem.Count + *instrument.CountResultItem
			result = "\nТы получил " + instrument.ItemsResult.View + " " + ToString(*instrument.CountResultItem) + " шт."
		}

		UpdateUser(update, User{Money: &updateUserMoney})
		UpdateUserItem(
			User{ID: user.ID},
			UserItem{
				ID:           userGetItem.ID,
				Count:        userGetItem.Count,
				CountUseLeft: userGetItem.Item.CountUse,
			})

		UpdateCelluleAfterGrowing(cellule, instrument)

		return "Ok", "Оно выросло!" + result

	} else {
		t := time.Now()
		UpdateCellule(cellule.ID,
			Cellule{
				ID:            cellule.ID,
				NextStateTime: &updateItemTime,
				LastGrowing:   &t,
			})
		return "Ok", "Вырастет " + updateItemTime.Format("15:04:05 02.01.06") + "!"

	}
}

func DesctructionItem(update tgbotapi.Update, cellule Cellule, user User, userGetItem UserItem, instrument Instrument) string {
	var ItemDestructionHp = *cellule.DestructionHp

	ItemDestructionHp = *cellule.DestructionHp - *instrument.Good.Destruction

	if isItemCrushed(cellule, ItemDestructionHp) {
		var result string
		if instrument.CountResultItem != nil {
			*userGetItem.Count = *userGetItem.Count + *instrument.CountResultItem
			result = "Ты получил " + instrument.ItemsResult.View + " " + ToString(*instrument.CountResultItem) + " шт."
		} else {
			result = "что то не так"
		}
		updateUserMoney := *user.Money - *cellule.Item.Cost

		UpdateUser(update, User{Money: &updateUserMoney})
		UpdateUserItem(
			User{ID: user.ID},
			UserItem{
				ID:           userGetItem.ID,
				Count:        userGetItem.Count,
				CountUseLeft: userGetItem.Item.CountUse,
			})

		UpdateCelluleAfterDestruction(cellule, instrument)

		return result
	} else {
		err := config.Db.
			Where(&Cellule{ID: cellule.ID}).
			Updates(Cellule{ID: cellule.ID, DestructionHp: &ItemDestructionHp}).
			Update("next_state_time", nil).
			Update("last_growing", nil).
			Error
		if err != nil {
			panic(err)
		}

		return "Попробуй еще.. (" + itemHpLeft(cellule, instrument) + ")"
	}
}

func isItemGrowed(cellule Cellule, updateItemTime time.Time) bool {

	if cellule.Item.Growing != nil && updateItemTime.Before(time.Now()) {
		return true
	} else {
		return false
	}
}

func isItemCrushed(cellule Cellule, ItemHp int) bool {
	if cellule.Item.DestructionHp != nil && ItemHp <= 0 {
		return true
	} else {
		return false
	}
}

func UserGetItemUpdateModels(update tgbotapi.Update, cellule Cellule, instrumentView string) string {
	userTgid := GetUserTgId(update)
	user := GetUser(User{TgId: userTgid})

	var userGetItem UserItem

	status, instrument := checkItemsOnNeededInstrument(cellule.Item.Instruments, instrumentView)
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

	if instrumentView == "👋" && len(cellule.Item.Instruments) == 0 {
		return UserGetItemWithHand(update, cellule, user, userGetItem)
	} else {
		return UserGetItemWithInstrument(update, cellule, user, *instrument, userGetItem)
	}

}

func canUserPayItem(user User, cellule Cellule) bool {
	return cellule.Item.Cost == nil || *user.Money >= *cellule.Item.Cost
}

func isUserHasMaxCountItem(item UserItem) bool {
	if item.Item.MaxCountUserHas == nil || *item.Count < *item.Item.MaxCountUserHas {
		return false
	}
	return true
}

func ViewItemInfo(location Location) string {
	cell := GetCellule(Cellule{MapsId: *location.MapsId, AxisX: *location.AxisX, AxisY: *location.AxisY})
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

	itemInfo = fmt.Sprintf("%s *%s* (_%s шт._) %s _%s_\n", cell.Item.View, cell.Item.Name, ToString(*cell.ItemCount), cell.Item.View, dressType)
	itemInfo = itemInfo + fmt.Sprintf("*Описание*: `%s`\n", *cell.Item.Description)

	if cell.Item.Healing != nil && *cell.Item.Healing != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Здоровье*: `+%s♥️`\n", ToString(*cell.Item.Healing))
	}
	if cell.Item.Damage != nil && *cell.Item.Damage != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Атака*: `+%s`💥️\n", ToString(*cell.Item.Damage))
	}
	if cell.Item.Satiety != nil && *cell.Item.Satiety != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Сытость*: `+%s`\U0001F9C3️\n", ToString(*cell.Item.Satiety))
	}
	if cell.Item.Cost != nil && *cell.Item.Cost != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Стоимость*: `%s`💰\n", ToString(*cell.Item.Cost))
	}
	if cell.Item.Destruction != nil && *cell.Item.Destruction != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Сила*: `%s %s`\n", ToString(*cell.Item.Destruction), cell.Item.View)
	}
	if cell.Item.DestructionHp != nil && *cell.Item.DestructionHp != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Прочность*: `%s`\n", ToString(*cell.Item.DestructionHp))
	}
	if cell.Item.Growing != nil && cell.NextStateTime != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Вырастет*: %s\n", cell.NextStateTime.Format("15:04:05 02.01.06"))
	} else {
		itemInfo = itemInfo + fmt.Sprintf("*Время роста*: `%s мин.`\n", ToString(*cell.Item.Growing))
	}
	if cell.Item.IntervalGrowing != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Интервал ускорения роста*: `раз в %s мин.`\n", ToString(*cell.Item.IntervalGrowing))
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
