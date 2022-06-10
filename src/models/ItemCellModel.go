package models

import (
	"fmt"
	"project0/config"
	"time"
)

type ItemCell struct {
	ID uint `gorm:"primaryKey"`
	ItemInfo
	Contain
}

type ItemInfo struct {
	Item          *Item
	ItemId        *int `gorm:"embedded"`
	ItemCount     *int `gorm:"embedded"`
	DestructionHp *int `gorm:"embedded"`
	BrokenTime    *time.Time
	GrowingTime   *time.Time
	LastGrowing   *time.Time
}

type Contain struct {
	ContainedItemId         *int `gorm:"embedded"`
	ContainedItem           *Item
	ContainedItemCount      *int `gorm:"embedded"`
	ContainedItemBrokenTime *time.Time
}

func (itemCell ItemCell) UpdateItemCell() {
	err := config.Db.
		Model(ItemCell{}).
		Where(ItemCell{ID: itemCell.ID}).
		Updates(itemCell).
		Error

	if err != nil {
		fmt.Println("Update ItemCell error:", err)
	}
}

func (itemCell ItemCell) UpdateContainedItemCellAfterBreaking() {

	if itemCell.ContainedItem.ItemAfterBreakingId == nil {
		itemCell.ContainedItemBrokenTime = nil
		itemCell.ContainedItemId = nil
		itemCell.ContainedItemCount = nil
	} else {
		itemCell.ContainedItemId = itemCell.ContainedItem.ItemAfterBreakingId
		itemCell.ContainedItemBrokenTime = nil

		if itemCell.ContainedItem.ItemAfterBreaking.Breaking != nil {
			bt := time.Now().Add(time.Duration(*itemCell.ContainedItem.ItemAfterBreaking.Breaking) * time.Minute)
			itemCell.ContainedItemBrokenTime = &bt
		}
	}

	err := config.Db.Model(ItemCell{}).
		Where(&ItemCell{ID: itemCell.ID}).
		Update("contained_item_id", itemCell.ContainedItemId).
		Update("contained_item_broken_time", itemCell.ContainedItemBrokenTime).
		Update("contained_item_count", itemCell.ContainedItemCount).
		Error
	if err != nil {
		panic(err)
	}

	fmt.Println("UpdateContainedItemCellAfterBreaking: itemCell=", itemCell.ID)
}

func (itemCell ItemCell) SpecialUpdateItemCell() {
	err := config.Db.Model(ItemCell{}).
		Where(&ItemCell{ID: itemCell.ID}).
		Update("item_id", itemCell.ItemId).
		Update("item_count", itemCell.ItemCount).
		Update("destruction_hp", itemCell.DestructionHp).
		Update("growing_time", itemCell.GrowingTime).
		Update("last_growing", itemCell.LastGrowing).
		Update("contained_item_id", itemCell.ContainedItemId).
		Update("contained_item_count", itemCell.ContainedItemCount).
		Update("contained_item_broken_time", itemCell.ContainedItemBrokenTime).
		Update("broken_time", itemCell.BrokenTime).
		Error
	if err != nil {
		panic(err)
	}
}

func (itemCell ItemCell) BecomeToContainedItem() ItemCell {
	itemCell.ItemId = itemCell.ContainedItemId
	itemCell.ItemCount = itemCell.ContainedItemCount
	itemCell.ContainedItemCount = nil
	itemCell.DestructionHp = nil
	itemCell.BrokenTime = nil
	itemCell.ContainedItemBrokenTime = nil
	itemCell.ContainedItemId = nil
	itemCell.GrowingTime = nil
	itemCell.LastGrowing = nil

	if itemCell.ContainedItem.DestructionHp != nil {
		itemCell.DestructionHp = nil
	}

	if itemCell.ContainedItem.Breaking != nil {
		tb := time.Now().Add(time.Duration(*itemCell.ContainedItem.Breaking) * time.Minute)
		itemCell.BrokenTime = &tb
	}

	return itemCell
}

func (itemCell ItemCell) BecomeToNextStageItem(instrument Instrument) ItemCell {
	itemCell.ItemId = instrument.NextStageItemId
	itemCell.DestructionHp = nil
	itemCell.BrokenTime = nil
	itemCell.GrowingTime = nil
	itemCell.LastGrowing = nil
	itemCell.ContainedItemId = nil
	itemCell.ContainedItemCount = nil
	itemCell.ContainedItemBrokenTime = nil

	if instrument.NextStageItem.DestructionHp != nil {
		itemCell.DestructionHp = instrument.NextStageItem.DestructionHp
	}

	if instrument.NextStageItem.Breaking != nil {
		brokenTime := time.Now().Add(time.Duration(*instrument.NextStageItem.Breaking) * time.Minute)
		itemCell.BrokenTime = &brokenTime
	}

	return itemCell
}

func (itemCell ItemCell) AddContainedItem(instrument Instrument) ItemCell {
	if itemCell.Item.Breaking == nil {
		count := *itemCell.ItemCount - 1
		itemCell.ItemCount = &count
	}

	itemCell.DestructionHp = nil
	itemCell.GrowingTime = nil
	itemCell.LastGrowing = nil
	itemCell.ContainedItemId = instrument.GrowingItemId
	itemCell.ContainedItem = instrument.GrowingItem

	if instrument.GrowingItem.Breaking != nil {
		nextStateTime := time.Now().Add(time.Duration(*instrument.GrowingItem.Breaking) * time.Minute)
		itemCell.ContainedItemBrokenTime = &nextStateTime
	}

	if itemCell.ContainedItemCount != nil {
		*itemCell.ContainedItemCount = *instrument.CountGrowingItem + *itemCell.ContainedItemCount
	} else {
		itemCell.ContainedItemCount = instrument.CountGrowingItem
	}

	return itemCell
}
