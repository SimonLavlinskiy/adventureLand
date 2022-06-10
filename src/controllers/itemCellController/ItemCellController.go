package itemCellController

import (
	"fmt"
	"project0/config"
	"project0/src/models"
	"time"
)

func GetOrCreateItemCell(cell models.Cell) (result models.ItemCell) {
	if cell.ItemCellId == nil {
		err := config.Db.Create(&result).Error
		if err != nil {
			fmt.Println("GetOrCreateItemCell (create itemCell) error: ", err)
		}

		err = config.Db.
			Model(models.Cell{}).
			Where(models.Cell{ID: cell.ID}).
			Update("item_cell_id", result.ID).
			Error
		if err != nil {
			fmt.Println("GetOrCreateItemCell (if) error: ", err)
		}
	} else {
		err := config.Db.
			Where(models.ItemCell{ID: *cell.ItemCellId}).
			First(&result).
			Error
		if err != nil {
			fmt.Println("GetOrCreateItemCell (else) error: ", err)
		}
	}

	return result
}

func UpdateItemCellUnderUser(cell models.Cell, userItem models.UserItem, count int) models.ItemCell {
	itemCell := GetOrCreateItemCell(cell)
	cell.ItemCellId = &itemCell.ID

	var nextStateTime *time.Time
	var breakTime *time.Time

	if userItem.Item.Growing != nil {
		nt := time.Now().Add(time.Duration(*userItem.Item.Growing) * time.Minute)
		nextStateTime = &nt
	}
	if userItem.Item.Breaking != nil {
		bt := time.Now().Add(time.Duration(*userItem.Item.Breaking) * time.Minute)
		breakTime = &bt
	}

	err := config.Db.
		Model(models.ItemCell{}).
		Where(&models.ItemCell{ID: itemCell.ID}).
		Update("item_id", userItem.ItemId).
		Update("item_count", count).
		Update("destruction_hp", nil).
		Update("growing_time", nextStateTime).
		Update("broken_time", breakTime).
		Update("last_growing", nil).
		Error
	if err != nil {
		fmt.Println(err)
	}

	return GetOrCreateItemCell(cell)
}

func UpdateItemCellAfterGrowing(itemCell models.ItemCell, instrument models.Instrument) {
	if instrument.GrowingItem != nil {
		itemCell = itemCell.AddContainedItem(instrument)

		if *itemCell.ItemCount == 0 {
			itemCell = itemCell.BecomeToContainedItem()
		}

	} else if instrument.NextStageItem != nil {
		itemCell = itemCell.BecomeToNextStageItem(instrument)
	} else {
		*itemCell.ItemCount = *itemCell.ItemCount - 1
		if *itemCell.ItemCount == 0 && itemCell.ContainedItem != nil {
			itemCell = itemCell.BecomeToContainedItem()
		}
	}

	itemCell.LastGrowing = nil
	itemCell.GrowingTime = nil

	itemCell.SpecialUpdateItemCell()
}

func UpdateItemCellAfterDestruction(itemCell models.ItemCell, instrument models.Instrument) {
	if itemCell.ItemCount == nil {
		return
	}
	if instrument.NextStageItem != nil {
		UpdatedItemCellOnNextItem(itemCell, instrument)
		return
	}

	itemCell.LastGrowing = nil
	itemCell.GrowingTime = nil
	*itemCell.ItemCount = *itemCell.ItemCount - 1

	if *itemCell.ItemCount <= 0 && itemCell.ContainedItemId != nil {
		itemCell = itemCell.BecomeToContainedItem()
	} else if *itemCell.ItemCount >= 1 && itemCell.Item.DestructionHp != nil {
		itemCell.DestructionHp = itemCell.Item.DestructionHp
	}
	itemCell.SpecialUpdateItemCell()
	return
}

func UpdatedItemCellOnNextItem(itemCell models.ItemCell, instrument models.Instrument) {
	itemCell.ItemId = instrument.NextStageItemId
	itemCell.DestructionHp = nil
	itemCell.ContainedItemCount = nil
	itemCell.ContainedItem = nil
	itemCell.ContainedItemBrokenTime = nil

	if instrument.NextStageItem.Breaking != nil {
		bt := time.Now().Add(time.Duration(*instrument.NextStageItem.Breaking) * time.Minute)
		itemCell.BrokenTime = &bt
	}

	if instrument.NextStageItem.DestructionHp != nil {
		itemCell.DestructionHp = instrument.NextStageItem.DestructionHp
	}

	if instrument.CountNextStageItem != nil {
		itemCell.ItemCount = instrument.CountNextStageItem
	}

	itemCell.SpecialUpdateItemCell()
}

func UpdateItemCellAfterBreaking(itemCell models.ItemCell) {

	if itemCell.ItemCount != nil && *itemCell.ItemCount > 1 && *itemCell.Item.BreakingType == "one" {
		count := *itemCell.ItemCount - 1
		itemCell.ItemCount = &count

		itemCell.BrokenTime = nil
		if itemCell.Item.Breaking != nil {
			bt := time.Now().Add(time.Duration(*itemCell.Item.Breaking) * time.Minute)
			itemCell.BrokenTime = &bt
		}

		if itemCell.Item.DestructionHp != nil {
			itemCell.DestructionHp = itemCell.Item.DestructionHp
		}

	} else if itemCell.ItemCount != nil && *itemCell.ItemCount <= 1 || *itemCell.Item.BreakingType == "all" {

		if itemCell.Item.ItemAfterBreakingId != nil {
			itemCell.ItemId = itemCell.Item.ItemAfterBreakingId

			itemCell.DestructionHp = nil
			if itemCell.Item.ItemAfterBreaking.DestructionHp != nil {
				itemCell.DestructionHp = itemCell.Item.ItemAfterBreaking.DestructionHp
			}

			itemCell.BrokenTime = nil
			if itemCell.Item.ItemAfterBreaking.Breaking != nil {
				bt := time.Now().Add(time.Duration(*itemCell.Item.ItemAfterBreaking.Breaking) * time.Minute)
				itemCell.BrokenTime = &bt
			}

		} else {

			if itemCell.ContainedItemId != nil {
				itemCell = itemCell.BecomeToContainedItem()
			} else {
				itemCell.ItemId = nil
				itemCell.ItemCount = nil
				itemCell.DestructionHp = nil
				itemCell.BrokenTime = nil
				itemCell.ContainedItemCount = nil
				itemCell.ContainedItemBrokenTime = nil
			}
		}
	}

	itemCell.LastGrowing = nil
	itemCell.GrowingTime = nil

	itemCell.SpecialUpdateItemCell()

	fmt.Println("UpdateItemCellAfterBreaking: itemCell=", itemCell.ID)
}
