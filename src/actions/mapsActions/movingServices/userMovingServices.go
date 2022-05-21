package movingServices

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/config"
	"project0/src/controllers/userItemController"
	"project0/src/controllers/userMapController"
	"project0/src/models"
)

func UserMoving(user models.User, cell models.Cell) (msg string, buttons tg.InlineKeyboardMarkup) {
	locMsg, err := UpdateLocation(user, cell)
	msgMap, buttons := userMapController.GetMyMap(user)

	if err != nil {
		if err.Error() == "user has not home" {
			buttons = BuyHomeKeyboard()
			msg = locMsg
		} else {
			msg = fmt.Sprintf("%s%s%s", msgMap, v.GetString("msg_separator"), locMsg)
		}
		return msg, buttons
	}

	lighterMsg, err := CheckUserHasLighter(user)
	if err != nil {
		msg = fmt.Sprintf("%s%s", v.GetString("msg_separator"), lighterMsg)
	}
	msg = fmt.Sprintf("%s%s", msgMap, msg)

	return msg, buttons
}

func BuyHomeKeyboard() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(fmt.Sprintf("🏘 Купить дом! 🏘 (%d 💰)", v.GetInt("main_info.cost_of_house")), "buyHome"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Отмена", "cancel"),
		),
	)
}

func CheckUserHasLighter(u models.User) (string, error) {
	if u.Clothes.LeftHandId != nil && u.Clothes.LeftHand.Type == "light" {

		res, err := userItemController.UpdateUserInstrument(u, *u.Clothes.LeftHand)
		if err != nil {
			return res, errors.New("lighter is updated")
		}

	}

	if u.Clothes.RightHandId != nil && u.Clothes.RightHand.Type == "light" {

		res, err := userItemController.UpdateUserInstrument(u, *u.Clothes.RightHand)
		if err != nil {
			return res, errors.New("lighter is updated")
		}

	}
	return "Ok", nil
}
func UpdateLocation(user models.User, cell models.Cell) (string, error) {
	var err error

	cell = isCellTeleport(cell)

	if cell, err = isCellHome(cell, user); err != nil {
		return "\nУ тебя еще нет дома, очень жаль...", errors.New("user has not home")
	}

	if !cell.CanStep || cell.Item != nil && cell.ItemCount != nil && *cell.ItemCount > 0 && !cell.Item.CanStep {
		return "\nСюда никак не пройти(", errors.New("can't get through")
	}

	userLocation := models.Location{MapsId: &cell.MapsId, UserID: user.ID, AxisX: &cell.AxisX, AxisY: &cell.AxisY}

	err = config.Db.
		Where(&models.Location{UserTgId: user.TgId}).
		Updates(&userLocation).
		Error
	if err != nil {
		panic(err)
	}

	user.UserStepCounter()

	return "Ok", nil
}

func isCellTeleport(cell models.Cell) models.Cell {
	if *cell.Type == "teleport" && cell.TeleportID != nil {
		return models.Cell{
			AxisX:  cell.Teleport.StartX,
			AxisY:  cell.Teleport.StartY,
			MapsId: cell.Teleport.MapId,
		}.GetCell()
	}
	return cell
}

func isCellHome(cell models.Cell, user models.User) (models.Cell, error) {
	if *cell.Type == "home" && user.HomeId != nil {
		mapId := int(user.Home.ID)
		return models.Cell{
			AxisX:  user.Home.StartX,
			AxisY:  user.Home.StartY,
			MapsId: mapId,
		}.GetCell(), nil
	} else if *cell.Type == "home" && user.HomeId == nil {
		return cell, errors.New("user has not home")
	}
	return cell, nil
}
