package itemCellController

import (
	"errors"
	"fmt"
	"project0/src/controllers/resultController"
	"project0/src/models"
	"project0/src/repositories"
	"time"
)

func GrowingItem(itemCell models.ItemCell, user models.User, instrument models.Instrument) (result string, err error) {
	if result, err = canUserGrowIt(itemCell, instrument); err != nil {
		return result, errors.New("not update userItem")
	}

	t := time.Now()
	updateItemTime := GetNewItemTime(itemCell, instrument)
	itemCell.BrokenTime = GetNewBreakingItemTime(itemCell, instrument)

	if isItemGrown(itemCell, updateItemTime) {
		UpdateItemCellAfterGrowing(itemCell, instrument)

		result = getResultAfterItemGrown(user, itemCell, instrument)
	} else {
		itemCell.GrowingTime = &updateItemTime
		itemCell.LastGrowing = &t
		itemCell.UpdateItemCell()

		growingTime := updateItemTime.Sub(t)
		h := growingTime.Truncate(time.Hour).Hours()
		m := growingTime.Truncate(time.Minute).Minutes() - growingTime.Truncate(time.Hour).Minutes()

		result = fmt.Sprintf("\U0001F973 Вырастет через %vч %vм !", h, m)
	}

	return result, nil
}

func GetNewItemTime(itemCell models.ItemCell, instrument models.Instrument) (updateItemTime time.Time) {
	if itemCell.GrowingTime == nil && itemCell.Item.Growing != nil {
		updateItemTime = time.Now()
		updateItemTime = updateItemTime.Add(time.Duration(*itemCell.Item.Growing) * time.Minute)
	} else {
		updateItemTime = *itemCell.GrowingTime
	}

	updateItemTime = updateItemTime.Add(-time.Duration(*instrument.Good.GrowingUpTime) * time.Minute)

	return updateItemTime
}

func GetNewBreakingItemTime(itemCell models.ItemCell, instrument models.Instrument) (brokenTime *time.Time) {
	bt := time.Time{}
	if itemCell.BrokenTime != nil {
		bt = *itemCell.BrokenTime
	} else if itemCell.Item.Breaking != nil {
		bt = time.Now().Add(time.Duration(*itemCell.Item.Breaking) * time.Minute)
	} else {
		return nil
	}

	bt = bt.Add(time.Duration(*instrument.Good.GrowingUpTime) * time.Minute)
	brokenTime = &bt

	return brokenTime
}

func canUserGrowIt(itemCell models.ItemCell, instrument models.Instrument) (result string, err error) {
	if itemCell.ContainedItemId != nil && itemCell.ContainedItemCount != nil && *itemCell.ContainedItemCount > 0 {
		result = "Необходимо сначала собрать урожай 👩‍🌾"
		err = errors.New("user can not growing")

		return result, err
	}

	if itemCell.LastGrowing != nil && time.Now().Before(itemCell.LastGrowing.Add(time.Duration(*itemCell.Item.IntervalGrowing)*time.Minute)) {
		nextTimeGrowing := itemCell.LastGrowing.
			Add(time.Duration(*itemCell.Item.IntervalGrowing) * time.Minute).Sub(time.Now())

		h := nextTimeGrowing.Truncate(time.Hour).Hours()
		m := nextTimeGrowing.Truncate(time.Minute).Minutes()

		growingTime := itemCell.GrowingTime.Sub(time.Now())
		hg := growingTime.Truncate(time.Hour).Hours()
		mg := growingTime.Truncate(time.Minute).Minutes() - growingTime.Truncate(time.Hour).Minutes()

		result = fmt.Sprintf("Ты уже использовал %s\nМожно будет повторить через *%vч %vм*! ⏰\n\n\U0001F973 Вырастет через %vч %vм !", instrument.Good.View, h, m, hg, mg)
		err = errors.New("user can not growing")
	}
	return result, err
}

func getResultAfterItemGrown(user models.User, itemCell models.ItemCell, instrument models.Instrument) (result string) {
	result = "Оно выросло!"

	if instrument.Result != nil {
		resultController.UserGetResult(user, *instrument.Result)
		result = fmt.Sprintf("%s\nТы получил %s %d шт. %s", result, instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)
	}

	updateUserMoney := *user.Money - *itemCell.Item.Cost
	repositories.UpdateUser(models.User{TgId: user.TgId, Money: &updateUserMoney})

	return result
}

func isItemGrown(itemCell models.ItemCell, updateItemTime time.Time) bool {
	if itemCell.Item.Growing != nil && updateItemTime.Before(time.Now()) {
		return true
	} else {
		return false
	}
}
