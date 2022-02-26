package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func GetOrCreateMyLocation(update tgbotapi.Update) Location {
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

func UpdateLocation(update tgbotapi.Update, LocationStruct Location) (Location, string) {
	var char []string
	if update.Message != nil {
		char = strings.Fields(update.Message.Text)
	} else {
		char = strings.Fields(update.CallbackQuery.Data)
	}
	userTgId := GetUserTgId(update)
	myLocation := GetOrCreateMyLocation(update)
	resCellule := GetCellule(Cellule{MapsId: *LocationStruct.MapsId, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY})

	if len(char) != 1 && *resCellule.Type == "teleport" && resCellule.TeleportID != nil {
		LocationStruct = Location{
			AxisX:  &resCellule.Teleport.StartX,
			AxisY:  &resCellule.Teleport.StartY,
			MapsId: &resCellule.Teleport.MapId,
		}
	}

	var result Cellule
	var err error

	err = config.Db.Preload("Item").First(&result, &Cellule{MapsId: *LocationStruct.MapsId, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY}).Error
	if err != nil {
		if err.Error() == "record not found" {
			return myLocation, "\nСюда никак не пройти("
		}
		panic(err)
	}

	if !result.CanStep || result.Item != nil && *result.ItemCount > 0 && !result.Item.CanStep {
		return myLocation, "\nСюда никак не пройти("
	} else {
		err = config.Db.Where(&Location{UserTgId: userTgId}).Updates(LocationStruct).Error
		if err != nil {
			panic(err)
		}
	}

	myLocation = GetOrCreateMyLocation(update)
	return myLocation, "Ok"
}

func GetLocationOnlineUser(userlocation Location, mapSize UserMap) []Location {
	var resultLocationsOnlineUser []Location

	err := config.Db.
		Preload("User", "online_map", true).
		Where(Cellule{MapsId: *userlocation.MapsId}).
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
