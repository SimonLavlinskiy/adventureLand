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

func UpdateLocation(user User, cell Cell) (string, error) {
	var err error

	cell = isCellTeleport(cell)

	if cell, err = isCellHome(cell, user); err != nil {
		return "\nУ тебя еще нет дома, очень жаль...", errors.New("user has not home")
	}

	if !cell.CanStep || cell.Item != nil && cell.ItemCount != nil && *cell.ItemCount > 0 && !cell.Item.CanStep {
		return "\nСюда никак не пройти(", errors.New("can't get through")
	}

	err = config.Db.
		Where(&Location{UserTgId: user.TgId}).
		Updates(&Location{MapsId: &cell.MapsId, AxisX: &cell.AxisX, AxisY: &cell.AxisY}).
		Error
	if err != nil {
		panic(err)
	}

	user.UserStepCounter()

	return "Ok", nil
}

func isCellTeleport(cell Cell) Cell {
	if *cell.Type == "teleport" && cell.TeleportID != nil {
		return Cell{
			AxisX:  cell.Teleport.StartX,
			AxisY:  cell.Teleport.StartY,
			MapsId: cell.Teleport.MapId,
		}.GetCell()
	}
	return cell
}

func isCellHome(cell Cell, user User) (Cell, error) {
	if *cell.Type == "home" && user.HomeId != nil {
		mapId := int(user.Home.ID)
		return Cell{
			AxisX:  user.Home.StartX,
			AxisY:  user.Home.StartY,
			MapsId: mapId,
		}.GetCell(), nil
	} else if *cell.Type == "home" && user.HomeId == nil {
		return cell, errors.New("user has not home")
	}
	return cell, nil
}

func GetLocationOnlineUser(userLocation Location, mapSize UserMap) []Location {
	var resultLocationsOnlineUser []Location

	err := config.Db.
		Preload("User").
		Where(Cell{MapsId: *userLocation.MapsId}).
		Where("axis_x >= ?", mapSize.LeftIndent).
		Where("axis_x <= ?", mapSize.RightIndent).
		Where("axis_y >= ?", mapSize.DownIndent).
		Where("axis_y <= ?", mapSize.UpperIndent).
		Where(fmt.Sprintf("update_at >= '%s'", time.Now().Add(time.Duration(v.GetInt("main_info.last_step_user_online_min"))*time.Minute).Format("2006-01-02 15:04:05"))).
		Order("axis_x ASC").
		Order("axis_y ASC").
		Find(&resultLocationsOnlineUser).Error

	if err != nil {
		panic(err)
	}

	return resultLocationsOnlineUser
}
