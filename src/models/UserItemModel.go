package models

import (
	"fmt"
	"project0/config"
)

type UserItem struct {
	ID           int  `gorm:"primaryKey"`
	Count        *int `gorm:"embedded"`
	CountUseLeft *int `gorm:"embedded"`
	UserId       int  `gorm:"embedded"`
	User         User
	ItemId       int `gorm:"embedded"`
	Item         Item
}

func (ui UserItem) UserGetUserItem() UserItem {
	zero := 0
	result := UserItem{
		UserId:       int(ui.User.ID),
		ItemId:       ui.ItemId,
		Count:        &zero,
		CountUseLeft: ui.CountUseLeft,
	}
	err := config.Db.
		Preload("Item").
		Preload("User").
		Where(ui).
		FirstOrCreate(&result).
		Error

	if err != nil {
		panic(err)
	}

	return result
}

func (u User) UserGetResultItem(r Result) {
	ui := UserItem{UserId: int(u.ID), ItemId: int(*r.ItemId)}.UserGetUserItem()
	resItemCount := *ui.Count + int(*r.CountItem)

	if ui.CountUseLeft == nil || *ui.CountUseLeft == 0 {
		item := Item{ID: uint(ui.ItemId)}.GetItem()
		ui.CountUseLeft = item.CountUse
	}

	u.UpdateUserItem(UserItem{ID: ui.ID, Count: &resItemCount, CountUseLeft: ui.CountUseLeft})
}

func (u User) UserGetResultExtraItem(r Result) {
	ui := UserItem{UserId: int(u.ID), ItemId: int(*r.SpecialItemId)}.UserGetUserItem()
	resItemCount := *ui.Count + int(*r.SpecialItemCount)

	if ui.CountUseLeft == nil {
		item := Item{ID: uint(ui.ItemId)}.GetItem()
		ui.CountUseLeft = item.CountUse
	}

	u.UpdateUserItem(UserItem{ID: ui.ID, Count: &resItemCount, CountUseLeft: ui.CountUseLeft})
}

func (u User) UpdateUserItem(ui UserItem) {
	err := config.Db.
		Where(UserItem{UserId: int(u.ID), ID: ui.ID}).
		Updates(&ui).
		Error

	if err != nil {
		panic(err)
	}
}

func (ui UserItem) GetFullDescriptionOfUserItem() string {
	userItem := ui.UserGetUserItem()
	var fullDescriptionUserItem string

	switch userItem.Item.Type {
	case "food":
		fullDescriptionUserItem = fmt.Sprintf("%s *%s* - %dшт.\n*Здоровье*: +%d ♥️️\n*Сытость*: +%d  \U0001F9C3\n", userItem.Item.View, userItem.Item.Name, *userItem.Count, *userItem.Item.Healing, *userItem.Item.Satiety)
	case "resource", "sprout", "furniture":
		fullDescriptionUserItem = fmt.Sprintf("%s *%s* - %dшт.\n", userItem.Item.View, userItem.Item.Name, *userItem.Count)
	}

	if userItem.Item.IsInventory == true {
		fullDescriptionUserItem = fmt.Sprintf("%s *%s* - %d шт.\n*Сила*: + %d💥\n", userItem.Item.View, userItem.Item.Name, *userItem.Count, *userItem.Item.Damage)
	}

	itemDescription := "Описания нет("

	if userItem.Item.Description != nil {
		itemDescription = fmt.Sprintf("*Описание*: %s", *userItem.Item.Description)
	}

	return fullDescriptionUserItem + itemDescription
}
