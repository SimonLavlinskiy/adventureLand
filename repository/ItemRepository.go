package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		instrumentMsg = UpdateUserInstrument(update, user, userInstrument)
		result = itemMsg + instrumentMsg
	case "hand":
		result = DesctructionItem(update, cellule, user, userGetItem, instrument)
	case "growing":
		status, itemMsg := GrowingItem(update, cellule, user, userGetItem, instrument)
		if status == "Ok" {
			instrumentMsg = UpdateUserInstrument(update, user, userInstrument)
		}
		result = itemMsg + instrumentMsg

	}

	return result
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
		return "Not ok", "Ты уже использовал " + instrument.Good.View + "\nМожно будет повторить " + cellule.LastGrowing.Add(time.Duration(*cellule.Item.IntervalGrowing)*time.Minute).Format("2006.01.02 15:04:05") + "!"
	}

	if cellule.NextStateTime == nil && cellule.Item.Growing != nil {
		updateItemTime = updateItemTime.Add(time.Duration(*cellule.Item.Growing) * time.Minute)
	} else {
		updateItemTime = *cellule.NextStateTime
	}
	updateItemTime = updateItemTime.Add(-time.Duration(*instrument.Good.GrowingUpTime) * time.Minute)
	fmt.Println(updateItemTime, isItemGrowed(cellule, updateItemTime))

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
		return "Ok", "Вырастет " + updateItemTime.Format("2006.01.02 15:04:05") + "!"

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
		UpdateCellule(cellule.ID,
			Cellule{
				ID:            cellule.ID,
				DestructionHp: &ItemDestructionHp,
			})

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
		sumCountItem := *userGetItem.Count + 1
		countAfterUserGetItem := *cellule.ItemCount - 1
		updateUserMoney := *user.Money - *cellule.Item.Cost
		countUseLeft := *userGetItem.CountUseLeft
		if *userGetItem.Count == 0 {
			countUseLeft = *userGetItem.Item.CountUse
		}

		UpdateUserItem(User{ID: user.ID}, UserItem{ID: userGetItem.ID, Count: &sumCountItem, CountUseLeft: &countUseLeft})
		UpdateUser(update, User{Money: &updateUserMoney})
		if countAfterUserGetItem != 0 || cellule.PrevItemID == nil {
			UpdateCellule(cellule.ID, Cellule{ItemCount: &countAfterUserGetItem})
		} else {
			UpdateCellOnPrevItem(cellule)
		}

		return "Ты получил " + userGetItem.Item.View + " 1 шт. (Осталось лежать еще " + ToString(countAfterUserGetItem) + ")"
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
