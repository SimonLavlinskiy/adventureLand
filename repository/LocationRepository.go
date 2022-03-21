package repository

import (
	"errors"
	"project0/config"
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

func GetOrCreateMyLocation(user User) Location {
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

func (u User) GetUserLocation() Location {
	result := Location{}

	err := config.Db.
		Preload("Maps").
		Where(&Location{UserID: u.ID}).
		First(result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func UpdateLocation(char []string, locStruct Location, user User) (string, error) {
	var err error

	cell := Cell{MapsId: *locStruct.MapsId, AxisX: *locStruct.AxisX, AxisY: *locStruct.AxisY}
	cell = cell.GetCell()

	locStruct = isCellTeleport(char, cell, locStruct)
	if locStruct, err = isCellHome(char, cell, locStruct, user); err != nil {
		return "\nУ тебя еще нет дома, очень жаль...", errors.New("user has not home")
	}

	var result Cell

	err = config.Db.
		Preload("Item").
		First(&result, &Cell{MapsId: *locStruct.MapsId, AxisX: *locStruct.AxisX, AxisY: *locStruct.AxisY}).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return "\nСюда никак не пройти(", errors.New("can't get through")
		}
	}

	if !result.CanStep || result.Item != nil && *result.ItemCount > 0 && !result.Item.CanStep {
		return "\nСюда никак не пройти(", errors.New("can't get through")
	}

	err = config.Db.
		Where(&Location{UserTgId: user.TgId}).
		Updates(locStruct).
		Error
	if err != nil {
		panic(err)
	}

	return "Ok", nil
}

func isCellTeleport(char []string, cell Cell, location Location) Location {
	if len(char) != 1 && *cell.Type == "teleport" && cell.TeleportID != nil {
		return Location{
			AxisX:  &cell.Teleport.StartX,
			AxisY:  &cell.Teleport.StartY,
			MapsId: &cell.Teleport.MapId,
		}
	}
	return location
}

func isCellHome(char []string, cell Cell, location Location, user User) (Location, error) {
	if len(char) != 1 && *cell.Type == "home" && user.HomeId != nil {
		masId := int(user.Home.ID)
		return Location{
			AxisX:  &user.Home.StartX,
			AxisY:  &user.Home.StartY,
			MapsId: &masId,
		}, nil
	} else if len(char) != 1 && *cell.Type == "home" && user.HomeId == nil {
		return location, errors.New("user has not home")
	}
	return location, nil
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
