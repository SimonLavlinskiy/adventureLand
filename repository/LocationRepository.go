package repository

import (
	"errors"
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	"time"
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
	UpdateAt time.Time `gorm:"autoUpdateTime"`
}

func GetOrCreateMyLocation(user User) Location {
	AsX := 155
	AsY := 54
	startMap := 50

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

	var resultCell Cell

	err = config.Db.
		Preload("Item").
		First(&resultCell, &Cell{MapsId: *locStruct.MapsId, AxisX: *locStruct.AxisX, AxisY: *locStruct.AxisY}).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return "\nСюда никак не пройти(", errors.New("can't get through")
		}
	}

	if !resultCell.CanStep || resultCell.Item != nil && resultCell.ItemCount != nil && *resultCell.ItemCount > 0 && !resultCell.Item.CanStep {
		return "\nСюда никак не пройти(", errors.New("can't get through")
	}

	err = config.Db.
		Where(&Location{UserTgId: user.TgId}).
		Updates(locStruct).
		Error
	if err != nil {
		panic(err)
	}

	user.UserStepCounter()

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

func GetLocationOnlineUser(userLocation Location, mapSize UserMap) []Location {
	var resultLocationsOnlineUser []Location

	err := config.Db.
		Preload("User").
		Where(Cell{MapsId: *userLocation.MapsId}).
		Where("axis_x >= " + ToString(mapSize.LeftIndent)).
		Where("axis_x <= " + ToString(mapSize.RightIndent)).
		Where("axis_y >= " + ToString(mapSize.DownIndent)).
		Where("axis_y <= " + ToString(mapSize.UpperIndent)).
		Where(fmt.Sprintf("update_at >= '%s'", time.Now().Add(time.Duration(v.GetInt("main_info.last_step_user_online_min"))*time.Minute).Format("2006-01-02 15:04:05"))).
		Order("axis_x ASC").
		Order("axis_y ASC").
		Find(&resultLocationsOnlineUser).Error

	if err != nil {
		panic(err)
	}

	return resultLocationsOnlineUser
}
