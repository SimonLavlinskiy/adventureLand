package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		status, res := UserGetItemUpdateModels(update, resultCell, char[0])

		if status == "Ok" {
			return res
		} else {
			return res
		}
	}

	return "Не получилось..."
}

func CheckUserHasInstrument(user User, instrument Instrument) (string, Item) {
	if user.LeftHandId != nil && *user.LeftHandId == *instrument.GoodId {
		return "Ok", *user.LeftHand
	}
	if user.RightHandId != nil && *user.RightHandId == *instrument.GoodId {
		return "Ok", *user.RightHand
	}
	fmt.Println("User dont have instrument")
	return "User dont have instrument", Item{}
}

func checkItemsOnNeededInstrument(instruments []Instrument, msgInstrumentView string) (string, *Instrument) {
	for _, instrument := range instruments {
		if instrument.Good.View == msgInstrumentView {
			return "Ok", &instrument
		}
	}
	return "Not ok", nil
}

func UserGetItemWithoutInstrument(update tgbotapi.Update, cellule Cellule, user User) (string, string) {
	userItem := GetOrCreateUserItem(update, *cellule.Item)
	sumCountItemResult := *userItem.Count + 1

	if isUserHasMaxCountItem(userItem) {
		return "Not ok", "👋 У тебя уже есть такой!"
	}
	if !canUserPayItem(user, cellule) {
		return "Not ok", "👋 Не хватает деняк!"
	}

	UpdateModelsWhenUserGetItem(update, user, userItem, cellule, nil, sumCountItemResult)

	return "Ok", "Ты получил " + userItem.Item.View + " 1 шт."
}

func UserGetItemWithInstrument(update tgbotapi.Update, cellule Cellule, user User, instrumentView string) (string, string) {
	userGetItem := UserItem{}
	var result string

	status, instrument := checkItemsOnNeededInstrument(cellule.Item.Instruments, instrumentView)
	if status != "Ok" {
		return "Not ok", "Предмет не поддается под таким инструментом"
	}

	status, userInstrument := CheckUserHasInstrument(user, *instrument)
	if status != "Ok" {
		return "Not ok", "Нет инструмента в руках"
	}

	if instrument.ItemsResultId == nil {
		userGetItem = GetOrCreateUserItem(update, *cellule.Item)
	} else {
		userGetItem = GetOrCreateUserItem(update, *instrument.ItemsResult)
	}

	if isUserHasMaxCountItem(userGetItem) {
		return "Not ok", "У тебя уже есть такой!"
	}
	if !canUserPayItem(user, cellule) {
		return "Not ok", "Не хватает деняк!"
	}

	switch instrument.Type {
	case "destruction":
		result = DesctructionItem(update, cellule, user, userGetItem, *instrument)
		result += UpdateUserInstrument(update, user, userInstrument)
	}

	return "Ok", result
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
	updateCell := Cellule{}
	ItemDestructionHp := *cellule.DestructionHp - *instrument.Good.Destruction

	var result string

	if cellule.Item.DestructionHp != nil && ItemDestructionHp > 0 {
		updateCell = Cellule{
			ID:            cellule.ID,
			DestructionHp: &ItemDestructionHp,
		}
		result = "Попробуй еще.. (" + itemHpLeft(cellule, instrument) + ")"
	} else if cellule.Item.DestructionHp != nil && ItemDestructionHp <= 0 {
		var updateCountItem int
		ItemDestructionHp = *cellule.Item.DestructionHp
		sumCountUserItemResult := *userGetItem.Count + *instrument.CountResultItem
		sumCountItemResult := *userGetItem.Count + *instrument.CountResultItem
		itemId := cellule.ItemID

		if *cellule.CountItem > 1 {
			updateCountItem = *cellule.CountItem - 1
		} else if *cellule.CountItem == 1 {
			updateCountItem = 0

			if instrument.NextStageItem != nil {
				itemId = instrument.NextStageItemId
			}

			if instrument.CountNextStageItem != nil {
				updateCountItem = *instrument.CountNextStageItem
			}

		}

		updateCell = Cellule{
			ID:            cellule.ID,
			ItemID:        itemId,
			DestructionHp: &ItemDestructionHp,
			CountItem:     &updateCountItem,
		}

		UpdateUserItem(User{ID: user.ID}, UserItem{ID: userGetItem.ID, Count: &sumCountUserItemResult, CountUseLeft: userGetItem.Item.CountUse})
		UpdateModelsWhenUserGetItem(update, user, userGetItem, cellule, &instrument, sumCountItemResult)

		result = "Ты получил " + instrument.ItemsResult.View + " " + ToString(*instrument.CountResultItem) + " шт."
	}

	UpdateCellule(updateCell.ID, updateCell)

	return result
}

func UserGetItemUpdateModels(update tgbotapi.Update, cellule Cellule, instrumentView string) (string, string) {
	userTgid := GetUserTgId(update)
	user := GetUser(User{TgId: userTgid})

	if instrumentView == "👋" && len(cellule.Item.Instruments) == 0 {
		return UserGetItemWithoutInstrument(update, cellule, user)
	}

	return UserGetItemWithInstrument(update, cellule, user, instrumentView)

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
