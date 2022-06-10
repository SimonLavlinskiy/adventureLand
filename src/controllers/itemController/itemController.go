package itemController

import (
	"fmt"
	"project0/config"
	"project0/src/models"
	"time"
)

func GetItems() (result []models.Item) {
	err := config.Db.
		Find(&result).
		Order("order by type asc").
		Error
	if err != nil {
		panic(err)
	}

	return result
}

func GetItemId(id uint) (result models.Item) {
	result.ID = id
	err := config.Db.
		Preload("Instruments").
		First(&result).
		Error
	if err != nil {
		panic(err)
	}

	return result
}

func ViewItemInfo(cell models.Cell) string {
	var itemInfo string
	var dressType string

	if cell.ItemCell.Item == nil {
		return "Ячейка пустая"
	}

	if cell.ItemCell.Item.DressType != nil {
		switch *cell.ItemCell.Item.DressType {
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

	itemInfo = fmt.Sprintf("%s *%s* _%s_\n", cell.ItemCell.Item.View, cell.ItemCell.Item.Name, dressType)
	if cell.ItemCell.ItemCount != nil {
		itemInfo += fmt.Sprintf("*Кол-во*: _%d шт._\n", *cell.ItemCell.ItemCount)
	}
	itemInfo = itemInfo + fmt.Sprintf("*Описание*: `%s`\n", *cell.ItemCell.Item.Description)

	if cell.ItemCell.ContainedItem != nil && cell.ItemCell.ContainedItemCount != nil && *cell.ItemCell.ContainedItemCount > 0 {
		itemInfo += fmt.Sprintf("*Содержит*: _%s %s - %d шт._\n", cell.ItemCell.ContainedItem.View, cell.ItemCell.ContainedItem.Name, *cell.ItemCell.ContainedItemCount)
	}

	if cell.ItemCell.Item.Healing != nil && *cell.ItemCell.Item.Healing != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Здоровье*: `+%d♥️`\n", *cell.ItemCell.Item.Healing)
	}
	if cell.ItemCell.Item.Damage != nil && *cell.ItemCell.Item.Damage != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Атака*: `+%d`💥️\n", *cell.ItemCell.Item.Damage)
	}
	if cell.ItemCell.Item.Satiety != nil && *cell.ItemCell.Item.Satiety != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Сытость*: `+%d`\U0001F9C3️\n", *cell.ItemCell.Item.Satiety)
	}
	if cell.ItemCell.Item.Cost != nil && *cell.ItemCell.Item.Cost != 0 && cell.NeedPay {
		itemInfo = itemInfo + fmt.Sprintf("*Стоимость*: `%d`💰\n", *cell.ItemCell.Item.Cost)
	}
	if cell.ItemCell.Item.Destruction != nil && *cell.ItemCell.Item.Destruction != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Сила*: `%d %s`\n", *cell.ItemCell.Item.Destruction, cell.ItemCell.Item.View)
	}
	if cell.ItemCell.DestructionHp != nil && cell.ItemCell.Item.DestructionHp != nil && *cell.ItemCell.Item.DestructionHp != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Прочность*: `%d`\n", *cell.ItemCell.DestructionHp)
	} else if cell.ItemCell.Item.DestructionHp != nil && *cell.ItemCell.Item.DestructionHp != 0 {
		itemInfo = itemInfo + fmt.Sprintf("*Прочность*: `%d`\n", *cell.ItemCell.Item.DestructionHp)
	}
	if cell.ItemCell.Item.Growing != nil && cell.ItemCell.GrowingTime != nil {
		t := cell.ItemCell.GrowingTime.Sub(time.Now())
		h := t.Truncate(time.Hour).Hours()
		m := t.Truncate(time.Minute).Minutes() - t.Truncate(time.Hour).Minutes()
		itemInfo = itemInfo + fmt.Sprintf("*\U0001F973 Вырастет через*: _%vч %vм_\n", h, m)
	} else if cell.ItemCell.Item.Growing != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Время роста*: `%d мин.`\n", *cell.ItemCell.Item.Growing)
	}
	if cell.ItemCell.Item.IntervalGrowing != nil {
		itemInfo = itemInfo + fmt.Sprintf("*Интервал ускорения роста*: `раз в %d мин.`\n", *cell.ItemCell.Item.IntervalGrowing)
	}
	if cell.ItemCell.LastGrowing != nil {
		t := time.Now().Sub(*cell.ItemCell.LastGrowing)
		m := t.Truncate(time.Minute).Minutes()
		itemInfo = itemInfo + fmt.Sprintf("*Последнее ускорение:* %vм назад\n", m)
	}
	if len(cell.ItemCell.Item.Instruments) != 0 {
		var itemsInstrument string
		for _, i := range cell.ItemCell.Item.Instruments {
			if i.GoodId != nil {
				itemsInstrument = itemsInstrument + fmt.Sprintf("%s - `%s`\n", i.Good.View, i.Good.Name)
			}
		}
		itemInfo = itemInfo + fmt.Sprintf("*Чем можно взаимодествовать*:\n%s", itemsInstrument)
	}

	return itemInfo
}
