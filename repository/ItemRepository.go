package repository

import (
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

func CheckUserHasInstrument(user User, instrument Instrument) (string, Item) {
	if instrument.Type == "hand" {
		return "Ok", *instrument.Good
	}
	if user.LeftHandId != nil && *user.LeftHandId == *instrument.GoodId {
		return "Ok", *user.LeftHand
	}
	if user.RightHandId != nil && *user.RightHandId == *instrument.GoodId {
		return "Ok", *user.RightHand
	}
	return "User dont have instrument", Item{}
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

	status, userInstrument := CheckUserHasInstrument(user, instrument)
	if status != "Ok" {
		return "Нет инструмента в руках"
	}

	switch instrument.Type {
	case "destruction":
		itemMsg := DesctructionItem(update, cellule, user, userGetItem, instrument)
		instrumentMsg := UpdateUserInstrument(update, user, userInstrument)
		result = itemMsg + instrumentMsg
	case "hand":
		result = DesctructionItem(update, cellule, user, userGetItem, instrument)
	case "growing":

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

func DesctructionItem(update tgbotapi.Update, cellule Cellule, user User, userGetItem UserItem, instrument Instrument) string {
	ItemDestructionHp := *cellule.DestructionHp - *instrument.Good.Destruction

	if isItemCrushed(cellule, ItemDestructionHp) {
		sumCountItem := *userGetItem.Count + *instrument.CountResultItem
		updateUserMoney := *user.Money - *cellule.Item.Cost

		UpdateUser(update, User{Money: &updateUserMoney})
		UpdateUserItem(
			User{ID: user.ID},
			UserItem{
				ID:           userGetItem.ID,
				Count:        &sumCountItem,
				CountUseLeft: userGetItem.Item.CountUse,
			})

		updateCell := updateModelCellule(cellule, instrument)
		UpdateCellule(updateCell.ID, updateCell)

		return "Ты получил " + instrument.ItemsResult.View + " " + ToString(*instrument.CountResultItem) + " шт."

	} else {
		UpdateCellule(cellule.ID,
			Cellule{
				ID:            cellule.ID,
				DestructionHp: &ItemDestructionHp,
			})

		return "Попробуй еще.. (" + itemHpLeft(cellule, instrument) + ")"
	}

}

func isItemCrushed(cellule Cellule, ItemHp int) bool {
	if cellule.Item.DestructionHp != nil && ItemHp <= 0 {
		return true
	} else {
		return false
	}
}

func updateModelCellule(cellule Cellule, instrument Instrument) Cellule {
	cellule.DestructionHp = cellule.Item.DestructionHp

	if *cellule.CountItem > 1 {
		*cellule.CountItem = *cellule.CountItem - 1
	} else {
		*cellule.CountItem = 0
	}

	if instrument.NextStageItem != nil {
		cellule.ItemID = instrument.NextStageItemId
	}

	if instrument.CountNextStageItem != nil {
		cellule.CountItem = instrument.CountNextStageItem
	}

	if instrument.NextStageItem != nil && instrument.NextStageItem.Growing != nil {
		*cellule.NextStateTime = time.Now().Add(time.Duration(*instrument.NextStageItem.Growing) * time.Minute)
	}

	return cellule
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
		countAfterUserGetItem := *cellule.CountItem - 1
		updateUserMoney := *user.Money - *cellule.Item.Cost

		UpdateUserItem(User{ID: user.ID}, UserItem{ID: userGetItem.ID, Count: &sumCountItem, CountUseLeft: userGetItem.Item.CountUse})
		UpdateUser(update, User{Money: &updateUserMoney})
		UpdateCellule(cellule.ID, Cellule{CountItem: &countAfterUserGetItem})

		return "Ты получил " + userGetItem.Item.View + " 1 шт."
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
