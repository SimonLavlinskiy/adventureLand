package repository

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strings"
)

type Location struct {
	ID       uint `gorm:"primaryKey"`
	UserTgId uint
	User     User
	UserID   uint
	AxisX    *int
	AxisY    *int
	MapsId   *int `gorm:"embedded"`
	Maps     Map
}

func GetOrCreateMyLocation(update tg.Update) Location {
	userTgId := GetUserTgId(update)
	user := GetUser(User{TgId: userTgId})

	AsX := 7
	AsY := 2
	startMap := 1

	result := Location{
		UserTgId: user.TgId,
		AxisX:    &AsX,
		AxisY:    &AsY,
		MapsId:   &startMap,
	}

	err := config.Db.
		Preload("Maps").
		Where(&Location{UserID: user.ID}).
		FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func UpdateLocation(update tg.Update, locStruct Location) (Location, string) {
	var char []string
	if update.Message != nil {
		char = strings.Fields(update.Message.Text)
	} else {
		char = strings.Fields(update.CallbackQuery.Data)
	}
	userTgId := GetUserTgId(update)
	usLoc := GetOrCreateMyLocation(update)
	cell := Cell{MapsId: *locStruct.MapsId, AxisX: *locStruct.AxisX, AxisY: *locStruct.AxisY}
	cell = cell.GetCell()

	if len(char) != 1 && *cell.Type == "teleport" && cell.TeleportID != nil {
		locStruct = Location{
			AxisX:  &cell.Teleport.StartX,
			AxisY:  &cell.Teleport.StartY,
			MapsId: &cell.Teleport.MapId,
		}
	}

	var result Cell
	var err error

	err = config.Db.
		Preload("Item").
		First(&result, &Cell{MapsId: *locStruct.MapsId, AxisX: *locStruct.AxisX, AxisY: *locStruct.AxisY}).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return usLoc, "\nСюда никак не пройти("
		}
		panic(err)
	}

	if !result.CanStep || result.Item != nil && *result.ItemCount > 0 && !result.Item.CanStep {
		return usLoc, "\nСюда никак не пройти("
	}

	err = config.Db.
		Where(&Location{UserTgId: userTgId}).
		Updates(locStruct).
		Error
	if err != nil {
		panic(err)
	}

	usLoc = GetOrCreateMyLocation(update)
	return usLoc, "Ok"
}

func GetLocationOnlineUser(userlocation Location, mapSize UserMap) []Location {
	var resultLocationsOnlineUser []Location

	err := config.Db.
		Preload("User", "online_map", true).
		Where(Cell{MapsId: *userlocation.MapsId}).
		Where("axis_x >= " + ToString(mapSize.leftIndent)).
		Where("axis_x <= " + ToString(mapSize.rightIndent)).
		Where("axis_y >= " + ToString(mapSize.downIndent)).
		Where("axis_y <= " + ToString(mapSize.upperIndent)).
		Order("axis_x ASC").
		Order("axis_y ASC").
		Find(&resultLocationsOnlineUser).Error

	if err != nil {
		panic(err)
	}

	return resultLocationsOnlineUser
}
