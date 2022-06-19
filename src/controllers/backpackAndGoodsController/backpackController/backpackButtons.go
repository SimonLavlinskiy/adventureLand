package backpackController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/cellController"
	"project0/src/models"
	str "strings"
)

func BackpackInlineKeyboard(user models.User, items []models.UserItem, i int, backpackType string) tg.InlineKeyboardMarkup {
	cell, _ := cellController.GetCellUnderUser(user)

	switch backpackType {
	case "food":
		return FoodListBackpackInlineKeyboard(cell, items, i)
	case "sprout":
		return SproutListBackpackInlineKeyboard(cell, items, i, backpackType)
	case "resource":
		return ResourceListBackpackInlineKeyboard(cell, items, i, backpackType)
	default:
		return DefaultListBackpackInlineKeyboard(cell, items, i, backpackType)
	}
}

func FoodListBackpackInlineKeyboard(cell models.Cell, items []models.UserItem, i int) tg.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт.   +%d ♥️️   +%d\U0001F9C3", items[i].Item.View, *items[i].Count, *items[i].Item.Healing, *items[i].Item.Satiety),
				fmt.Sprintf("%s %d %d food", v.GetString("callback_char.description"), items[i].ID, i)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("🍽 1шт", fmt.Sprintf("%s %d %d", v.GetString("callback_char.eat_food"), items[i].ID, i)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d food", v.GetString("callback_char.backpack_moving"), i-1)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d food", v.GetString("callback_char.count_of_throw_out"), items[i].ID, i)),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d food", v.GetString("callback_char.backpack_moving"), i+1)),
		),
		tg.NewInlineKeyboardRow(
			ButtonDestroyOrSellItem(cell, items, i, "food"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
		),
	)
}

func SproutListBackpackInlineKeyboard(cell models.Cell, items []models.UserItem, i int, itemType string) tg.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт. - %s", items[i].Item.View, *items[i].Count, *items[i].Item.Description),
				fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.description"), items[i].ID, i, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋\U0001F9A0🗺", fmt.Sprintf("%s %d %d 1 %s", v.GetString("callback_char.throw_out_item"), items[i].ID, i, itemType)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i-1, itemType)),
		),
		tg.NewInlineKeyboardRow(
			ButtonDestroyOrSellItem(cell, items, i, itemType),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i+1, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
		),
	)
}

func ResourceListBackpackInlineKeyboard(cell models.Cell, items []models.UserItem, i int, itemType string) tg.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт.", items[i].Item.View, *items[i].Count),
				fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.description"), items[i].ID, i, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.count_of_throw_out"), items[i].ID, i, itemType)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i-1, itemType)),
		),
		tg.NewInlineKeyboardRow(
			ButtonDestroyOrSellItem(cell, items, i, itemType),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i+1, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
		),
	)
}

func DefaultListBackpackInlineKeyboard(cell models.Cell, items []models.UserItem, i int, itemType string) tg.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %dшт.", items[i].Item.View, *items[i].Count),
				fmt.Sprintf("%s %d %d %s", v.GetString("callback_char.description"), items[i].ID, i, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("👋🗑🗺", fmt.Sprintf("%s %d %d 1 %s", v.GetString("callback_char.throw_out_item"), items[i].ID, i, itemType)),
			tg.NewInlineKeyboardButtonData("🔺", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i-1, itemType)),
		),
		tg.NewInlineKeyboardRow(
			ButtonDestroyOrSellItem(cell, items, i, itemType),
			tg.NewInlineKeyboardButtonData("🔻", fmt.Sprintf("%s %d %s", v.GetString("callback_char.backpack_moving"), i+1, itemType)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
		),
	)
}

func BackpackCategoryKeyboard() (msgText string, buttons tg.InlineKeyboardMarkup) {
	categories := str.Fields(v.GetString("user_location.item_categories.categories"))

	var rows [][]tg.InlineKeyboardButton

	for _, category := range categories {
		rows = append(rows, tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s", v.GetString(fmt.Sprintf("user_location.item_categories.%s", category))),
			fmt.Sprintf("%s", fmt.Sprintf("category %s", category)),
		)))
	}

	cancel := tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData("Отмена", "cancel"))
	rows = append(rows, cancel)

	msgText = fmt.Sprintf("🎒 *Рюкзачок*\n%s", v.GetString("user_location.item_categories.category_title"))
	buttons = tg.NewInlineKeyboardMarkup(
		rows...,
	)

	return msgText, buttons
}

func ButtonDestroyOrSellItem(cell models.Cell, items []models.UserItem, i int, itemType string) (button tg.InlineKeyboardButton) {
	if cell.CanSell && items[i].Item.Cost != nil && *items[i].Item.Cost > 0 {
		var sellCost int
		if *items[i].Item.Cost == 1 {
			sellCost = *items[i].Item.Cost
		} else {
			sellCost = *items[i].Item.Cost / 2
		}
		button = tg.NewInlineKeyboardButtonData(fmt.Sprintf("💰 %d 💰", sellCost), fmt.Sprintf("%s %d %d %s false", v.GetString("callback_char.select_count_sell_item"), items[i].ID, i, itemType))
	} else {
		button = tg.NewInlineKeyboardButtonData("💥🗑💥", fmt.Sprintf("%s %d %d %s false", v.GetString("callback_char.delete_item"), items[i].ID, i, itemType))
	}

	return button
}
