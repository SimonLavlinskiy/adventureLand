package repositories

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	"project0/src/models"
	"time"
)

type StartMapStruct struct {
	AsX      int
	AsY      int
	startMap int
}

var StartLearningMap = StartMapStruct{
	AsX:      2,
	AsY:      0,
	startMap: 4,
}

func GetOrCreateMyLocation(user models.User) models.Location {

	result := models.Location{
		UserTgId: user.TgId,
		AxisX:    &StartLearningMap.AsX,
		AxisY:    &StartLearningMap.AsY,
		MapsId:   &StartLearningMap.startMap,
	}

	err := config.Db.
		Preload("Maps").
		Where(&models.Location{UserID: user.ID}).
		FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func GetLocationOnlineUser(userLocation models.Location, mapSize models.UserMap) []models.Location {
	var resultLocationsOnlineUser []models.Location

	err := config.Db.
		Preload("User").
		Where(models.Cell{MapsId: *userLocation.MapsId}).
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

	var results []models.Location
	for _, loc := range resultLocationsOnlineUser {
		if loc.User.MenuLocation == "sleep" {
			results = append(results, loc)
		}
	}

	return results
}

func UpdateLocation(location models.Location) {
	err := config.Db.
		Where(&models.Location{UserTgId: location.UserTgId}).
		Updates(&location).
		Error

	if err != nil {
		fmt.Println(err)
	}
}
