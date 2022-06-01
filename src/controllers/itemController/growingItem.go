package itemController

import (
	"errors"
	"fmt"
	"project0/src/controllers/resultController"
	"project0/src/models"
	"project0/src/repositories"
	"time"
)

func GrowingItem(cell models.Cell, user models.User, instrument models.Instrument) (result string, err error) {
	if result, err := canUserGrowIt(cell, instrument); err != nil {
		return result, errors.New("not update userItem")
	}

	updateItemTime := GetNewItemTime(cell, instrument)

	if isItemGrown(cell, updateItemTime) {
		cell.UpdateCellAfterGrowing(instrument)
		result = getResultAfterItemGrown(user, cell, instrument)
		return result, nil
	}

	t := time.Now()
	cell.NextStateTime = &updateItemTime
	cell.LastGrowing = &t
	cell.UpdateCell()

	growingTime := updateItemTime.Sub(t)
	h := growingTime.Truncate(time.Hour).Hours()
	m := growingTime.Truncate(time.Minute).Minutes() - growingTime.Truncate(time.Hour).Minutes()

	result = fmt.Sprintf("\U0001F973 Вырастет через %vч %vм !", h, m)
	return result, nil
}

func GetNewItemTime(cell models.Cell, instrument models.Instrument) (updateItemTime time.Time) {
	if cell.NextStateTime == nil && cell.Item.Growing != nil {
		updateItemTime = time.Now()
		updateItemTime = updateItemTime.Add(time.Duration(*cell.Item.Growing) * time.Minute)
	} else {
		updateItemTime = *cell.NextStateTime
	}

	updateItemTime = updateItemTime.Add(-time.Duration(*instrument.Good.GrowingUpTime) * time.Minute)

	return updateItemTime
}

func canUserGrowIt(cell models.Cell, instrument models.Instrument) (result string, err error) {
	if cell.LastGrowing != nil && time.Now().Before(cell.LastGrowing.Add(time.Duration(*cell.Item.IntervalGrowing)*time.Minute)) {
		nextTimeGrowing := cell.LastGrowing.
			Add(time.Duration(*cell.Item.IntervalGrowing) * time.Minute).Sub(time.Now())

		h := nextTimeGrowing.Truncate(time.Hour).Hours()
		m := nextTimeGrowing.Truncate(time.Minute).Minutes()

		result = fmt.Sprintf("Ты уже использовал %s\n\nМожно будет повторить через %vч %vм!", instrument.Good.View, h, m)
		err = errors.New("user can not growing")
	}
	return result, err
}

func getResultAfterItemGrown(user models.User, cell models.Cell, instrument models.Instrument) (result string) {
	result = "Оно выросло!"

	if instrument.Result != nil {
		resultController.UserGetResult(user, *instrument.Result)
		result = fmt.Sprintf("%s\nТы получил %s %d шт. %s", result, instrument.Result.Item.View, *instrument.Result.CountItem, instrument.Result.Item.Name)
	}

	updateUserMoney := *user.Money - *cell.Item.Cost
	repositories.UpdateUser(models.User{TgId: user.TgId, Money: &updateUserMoney})

	return result
}
