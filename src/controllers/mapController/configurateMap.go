package mapController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"project0/src/models"
	"project0/src/repositories"
	"strings"
	"time"
)

var displaySize = 6

func DefaultButtons(center string) models.MapButtons {
	return models.MapButtons{
		Up:     "🔼",
		Left:   "◀️️",
		Right:  "▶️",
		Down:   "🔽",
		Center: center,
	}
}

func DefaultUserMap(location models.Location, displaySize int) models.UserMap {
	return models.UserMap{
		LeftIndent:  *location.AxisX - displaySize,
		RightIndent: *location.AxisX + displaySize,
		UpperIndent: *location.AxisY + displaySize,
		DownIndent:  *location.AxisY - displaySize,
	}
}

func GetMyMap(us models.User) (textMessage string, buttons tg.InlineKeyboardMarkup) {
	loc := repositories.GetOrCreateMyLocation(us)
	mapSize := CalculateUserMapBorder(loc, loc.Maps)
	messageMap := fmt.Sprintf("*Карта*: _%s_ *X*: _%d_  *Y*: _%d_", loc.Maps.Name, *loc.AxisX, *loc.AxisY)

	type Point = [2]int
	m := map[Point]models.Cell{}

	resultCell := repositories.GetCellsUserMap(mapSize, loc)

	for _, cell := range resultCell {
		m[Point{cell.AxisX, cell.AxisY}] = cell
	}

	resultLocationsOnlineUser := repositories.GetLocationOnlineUser(loc, mapSize)
	for _, usersLoc := range resultLocationsOnlineUser {
		if usersLoc.User.ID != 0 {
			m[Point{*usersLoc.AxisX, *usersLoc.AxisY}] = models.Cell{View: usersLoc.User.Avatar, ID: m[Point{*usersLoc.AxisX, *usersLoc.AxisY}].ID}
		}
	}

	buttons = CalculateButtonMap(loc, us, m)

	m[Point{*loc.AxisX, *loc.AxisY}] = models.Cell{View: us.Avatar, ID: m[Point{*loc.AxisX, *loc.AxisY}].ID}
	if us.MenuLocation == "Меню" {
		m[Point{*loc.AxisX, *loc.AxisY + 1}] = models.Cell{View: "⬇️", ID: m[Point{*loc.AxisX, *loc.AxisY}].ID}
	}

	Maps := configurationMap(mapSize, loc.Maps, loc, us, m)

	for i, row := range Maps {
		if i >= mapSize.DownIndent && i <= mapSize.UpperIndent {
			messageMap = fmt.Sprintf("%s\n%s", strings.Join(row, ""), messageMap)
		}
	}

	updatedUser := repositories.GetUser(models.User{ID: us.ID})
	messageMap = fmt.Sprintf("%s\n%s", GetStatsLine(updatedUser), messageMap)

	return messageMap, buttons
}

func configurationMap(mapSize models.UserMap, resMap models.Map, resLocation models.Location, user models.User, m map[[2]int]models.Cell) [][]string {
	t := time.Now()
	day := "06" <= t.Format("15") && t.Format("15") <= "23"
	Maps := make([][]string, resMap.SizeY+1)

	switch true {
	case day && resLocation.Maps.DayType == "default", resLocation.Maps.DayType == "alwaysDay":
		Maps = DayMap(Maps, mapSize, m, resLocation, user)
	case resLocation.Maps.DayType == "alwaysSuperNight":
		Maps = SuperNightMap(Maps, mapSize, m, resLocation, user)
	default: // !day or alwaysNight
		Maps = NightMap(Maps, mapSize, m, resLocation, user)
	}

	return Maps
}

func DayMap(Maps [][]string, mapSize models.UserMap, m map[[2]int]models.Cell, resLocation models.Location, user models.User) [][]string {
	type Point [2]int

	for y := range Maps {
		for x := mapSize.LeftIndent; x <= mapSize.RightIndent; x++ {
			if m[Point{x, y}].ID != 0 || m[Point{x, y}] == m[Point{*resLocation.AxisX, *resLocation.AxisY}] {
				appendVisibleUserZone(m, x, y, Maps, resLocation, user)
			} else {
				Maps[y] = append(Maps[y], resLocation.Maps.EmptySpaceSymbol)
			}
		}
	}

	return Maps
}

func NightMap(Maps [][]string, mapSize models.UserMap, m map[[2]int]models.Cell, resLocation models.Location, user models.User) [][]string {
	type Point [2]int

	for y := range Maps {
		for x := mapSize.LeftIndent; x <= mapSize.RightIndent; x++ {
			if CalculateNightMap(user, resLocation, x, y) && m[Point{x, y}].ID != 0 {
				appendVisibleUserZone(m, x, y, Maps, resLocation, user)
			} else if (y+x)%2 == 1 || m[Point{x, y}].ID != 0 && (y+x)%2 == 1 {
				Maps[y] = append(Maps[y], "⬛️")
			} else if (y+x)%2 == 0 || m[Point{x, y}].ID != 0 && (y+x)%2 == 0 {
				Maps[y] = append(Maps[y], "✨")
			} else {
				Maps[y] = append(Maps[y], resLocation.Maps.EmptySpaceSymbol)
			}
		}
	}

	return Maps
}

func SuperNightMap(Maps [][]string, mapSize models.UserMap, m map[[2]int]models.Cell, resLocation models.Location, user models.User) [][]string {
	type Point [2]int

	for y := range Maps {
		for x := mapSize.LeftIndent; x <= mapSize.RightIndent; x++ {
			if CalculateSuperNightMap(user, resLocation, x, y) && m[Point{x, y}].ID != 0 {
				appendVisibleUserZone(m, x, y, Maps, resLocation, user)
			} else if m[Point{x, y}].ID != 0 {
				Maps[y] = append(Maps[y], RandNightCell())
			} else {
				Maps[y] = append(Maps[y], resLocation.Maps.EmptySpaceSymbol)
			}
		}
	}

	return Maps
}

func RandNightCell() string {
	randomInt := rand.Intn(100)
	if randomInt < 5 {
		return "💀"
	} else if randomInt > 5 && randomInt < 10 {
		return "👹"
	}
	return "⬛️"
}
