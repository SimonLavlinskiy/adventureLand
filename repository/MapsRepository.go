package repository

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/config"
	"strings"
	"time"
)

type Map struct {
	ID               uint   `gorm:"primaryKey"`
	Name             string `gorm:"embedded"`
	SizeX            int    `gorm:"embedded"`
	SizeY            int    `gorm:"embedded"`
	StartX           int    `gorm:"embedded"`
	StartY           int    `gorm:"embedded"`
	DayType          string `gorm:"embedded"`
	EmptySpaceSymbol string `gorm:"embedded"`
}

var displaySize = 6

func (m Map) CreateMap() Map {
	err := config.Db.
		Create(&m).
		Error

	if err != nil {
		fmt.Println(err)
	}

	return m
}

func GetMaps() []Map {
	var result []Map

	err := config.Db.Find(&result).Error

	if err != nil {
		fmt.Println("Нет карт!")
	}

	return result
}

type UserMap struct {
	LeftIndent  int
	RightIndent int
	UpperIndent int
	DownIndent  int
}

type MapButtons struct {
	Up         string
	UpData     string
	Left       string
	LeftData   string
	Right      string
	RightData  string
	Down       string
	DownData   string
	Center     string
	CenterData string
}

func DefaultButtons(center string) MapButtons {
	return MapButtons{
		Up:     "🔼",
		Left:   "◀️️",
		Right:  "▶️",
		Down:   "🔽",
		Center: center,
	}
}

func DefaultUserMap(location Location, displaySize int) UserMap {
	return UserMap{
		LeftIndent:  *location.AxisX - displaySize,
		RightIndent: *location.AxisX + displaySize,
		UpperIndent: *location.AxisY + displaySize,
		DownIndent:  *location.AxisY - displaySize,
	}
}

func GetMyMap(us User) (textMessage string, buttons tg.InlineKeyboardMarkup) {
	loc := GetOrCreateMyLocation(us)
	mapSize := CalculateUserMapBorder(loc, loc.Maps)
	messageMap := fmt.Sprintf("*Карта*: _%s_ *X*: _%d_  *Y*: _%d_", loc.Maps.Name, *loc.AxisX, *loc.AxisY)

	type Point = [2]int
	m := map[Point]Cell{}

	var resultCell []Cell

	err := config.Db.
		Preload("Item").
		Preload("Teleport").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.Result").
		Preload("Item.Instruments.NextStageItem").
		Where(Cell{MapsId: *loc.MapsId}).
		Where("axis_x >= " + ToString(mapSize.LeftIndent)).
		Where("axis_x <= " + ToString(mapSize.RightIndent)).
		Where("axis_y >= " + ToString(mapSize.DownIndent)).
		Where("axis_y <= " + ToString(mapSize.UpperIndent)).
		Order("axis_x ASC").
		Order("axis_y ASC").
		Find(&resultCell).Error
	if err != nil {
		panic(err)
	}

	for _, cell := range resultCell {
		m[Point{cell.AxisX, cell.AxisY}] = cell
	}

	resultLocationsOnlineUser := GetLocationOnlineUser(loc, mapSize)
	for _, location := range resultLocationsOnlineUser {
		if location.User.ID != 0 {
			m[Point{*location.AxisX, *location.AxisY}] = Cell{View: location.User.Avatar, ID: m[Point{*location.AxisX, *location.AxisY}].ID}
		}
	}

	buttons = CalculateButtonMap(loc, us, m)

	m[Point{*loc.AxisX, *loc.AxisY}] = Cell{View: us.Avatar, ID: m[Point{*loc.AxisX, *loc.AxisY}].ID}
	if us.MenuLocation == "Меню" {
		m[Point{*loc.AxisX, *loc.AxisY + 1}] = Cell{View: "⬇️", ID: m[Point{*loc.AxisX, *loc.AxisY}].ID}
	}

	Maps := configurationMap(mapSize, loc.Maps, loc, us, m)

	for i, row := range Maps {
		if i >= mapSize.DownIndent && i <= mapSize.UpperIndent {
			messageMap = fmt.Sprintf("%s\n%s", strings.Join(row, ""), messageMap)
		}
	}

	messageMap = fmt.Sprintf("%s\n%s", GetStatsLine(us), messageMap)

	return messageMap, buttons
}

func CalculateNightMap(user User, l Location, x int, y int) bool {

	if (*l.AxisX == x-2 || *l.AxisX == x+2) && (*l.AxisY >= y-1 && *l.AxisY <= y+1) {
		return true
	}
	if (*l.AxisX >= x-1 && *l.AxisX <= x+1) && (*l.AxisY >= y-2 && *l.AxisY <= y+2) {
		return true
	}

	if user.LeftHandId != nil && user.LeftHand.Type == "light" || user.RightHandId != nil && user.RightHand.Type == "light" {
		if (*l.AxisX == x-4 || *l.AxisX == x+4) && (*l.AxisY >= y-1 && *l.AxisY <= y+1) {
			return true
		}
		if (*l.AxisX == x-3 || *l.AxisX == x+3) && (*l.AxisY >= y-2 && *l.AxisY <= y+2) {
			return true
		}
		if (*l.AxisX == x-2 || *l.AxisX == x+2) && (*l.AxisY >= y-3 && *l.AxisY <= y+3) {
			return true
		}
		if (*l.AxisX >= x-1 && *l.AxisX <= x+1) && (*l.AxisY >= y-4 && *l.AxisY <= y+4) {
			return true
		}
	}

	return false
}

func configurationMap(mapSize UserMap, resMap Map, resLocation Location, user User, m map[[2]int]Cell) [][]string {
	t := time.Now()
	day := "06" <= t.Format("15") && t.Format("15") <= "23"
	type Point [2]int
	Maps := make([][]string, resMap.SizeY+1)

	if (day && resLocation.Maps.DayType == "default") || resLocation.Maps.DayType == "alwaysDay" {
		DayMap(Maps, mapSize, m, resLocation)
	} else {
		NightMap(Maps, mapSize, m, resLocation, user)
	}
	return Maps
}

func DayMap(Maps [][]string, mapSize UserMap, m map[[2]int]Cell, resLocation Location) {
	type Point [2]int

	for y := range Maps {
		for x := mapSize.LeftIndent; x <= mapSize.RightIndent; x++ {
			if m[Point{x, y}].ID != 0 || m[Point{x, y}] == m[Point{*resLocation.AxisX, *resLocation.AxisY}] {
				appendVisibleUserZone(m, x, y, Maps)
			} else {
				Maps[y] = append(Maps[y], resLocation.Maps.EmptySpaceSymbol)
			}
		}
	}
}

func NightMap(Maps [][]string, mapSize UserMap, m map[[2]int]Cell, resLocation Location, user User) {
	type Point [2]int

	for y := range Maps {
		for x := mapSize.LeftIndent; x <= mapSize.RightIndent; x++ {
			if CalculateNightMap(user, resLocation, x, y) && m[Point{x, y}].ID != 0 {
				appendVisibleUserZone(m, x, y, Maps)
			} else if (y+x)%2 == 1 || m[Point{x, y}].ID != 0 && (y+x)%2 == 1 {
				Maps[y] = append(Maps[y], "⬛️")
			} else if (y+x)%2 == 0 || m[Point{x, y}].ID != 0 && (y+x)%2 == 0 {
				Maps[y] = append(Maps[y], "✨")
			} else {
				Maps[y] = append(Maps[y], resLocation.Maps.EmptySpaceSymbol)
			}
		}
	}
}

func appendVisibleUserZone(m map[[2]int]Cell, x int, y int, Maps [][]string) {
	type Point [2]int

	if m[Point{x, y}].IsItem() || m[Point{x, y}].IsWorkbench() || m[Point{x, y}].IsSwap() || m[Point{x, y}].IsQuest() || m[Point{x, y}].IsChat() {
		Maps[y] = append(Maps[y], m[Point{x, y}].Item.View)
	} else {
		Maps[y] = append(Maps[y], m[Point{x, y}].View)
	}
}

func CalculateUserMapBorder(resLocation Location, resMap Map) UserMap {
	d := displaySize
	if d*2 > resMap.SizeX || d*2 > resMap.SizeY {
		if resMap.SizeX-resMap.SizeY > 0 {
			d = resMap.SizeY / 2
		} else {
			d = resMap.SizeX / 2
		}
	}
	mapSize := DefaultUserMap(resLocation, d)

	if *resLocation.AxisX < d {
		mapSize.LeftIndent = 0
		mapSize.RightIndent = d * 2
	}
	if *resLocation.AxisY < d {
		mapSize.DownIndent = 0
		mapSize.UpperIndent = d * 2
	}
	if mapSize.RightIndent >= resMap.SizeX && *resLocation.AxisX > d {
		mapSize.LeftIndent = resMap.SizeX - d*2
		mapSize.RightIndent = resMap.SizeX
	}
	if mapSize.UpperIndent >= resMap.SizeY && *resLocation.AxisY > d {
		mapSize.DownIndent = resMap.SizeY - d*2
		mapSize.UpperIndent = resMap.SizeY
	}

	return mapSize
}

func CalculateButtonMap(resLocation Location, resUser User, m map[[2]int]Cell) tg.InlineKeyboardMarkup {
	type Point = [2]int

	buttons := DefaultButtons(resUser.Avatar)

	var CellsAroundUser []Cell
	CellsAroundUser = append(CellsAroundUser,
		m[Point{*resLocation.AxisX, *resLocation.AxisY + 1}],
		m[Point{*resLocation.AxisX, *resLocation.AxisY - 1}],
		m[Point{*resLocation.AxisX + 1, *resLocation.AxisY}],
		m[Point{*resLocation.AxisX - 1, *resLocation.AxisY}],
	)

	buttons = PutButton(CellsAroundUser, buttons, resUser)

	return CreateMapKeyboard(buttons)
}

func PutButton(CellsAroundUser []Cell, btn MapButtons, resUser User) MapButtons {
	for i, cell := range CellsAroundUser {
		switch true {
		case cell.IsDefaultCell():
			text := cell.View
			switch i {
			case 0:
				btn.Up, btn.UpData = text, text
			case 1:
				btn.Down, btn.DownData = text, text
			case 2:
				btn.Right, btn.RightData = text, text
			case 3:
				btn.Left, btn.LeftData = text, text
			}
		case cell.IsTeleport() || cell.IsHome():
			button := fmt.Sprintf("%s%s", resUser.Avatar, cell.View)
			data := fmt.Sprintf("move %d", cell.ID)
			switch i {
			case 0:
				btn.Up = button
				btn.UpData = data
			case 1:
				btn.Down = button
				btn.DownData = data
			case 2:
				btn.Right = button
				btn.RightData = data
			case 3:
				btn.Left = button
				btn.LeftData = data
			}
		case cell.IsWorkbench() || cell.IsQuest() || cell.IsChat():
			var data string
			if cell.IsWorkbench() {
				data = fmt.Sprintf("wrench %d", cell.ID)
			} else if cell.IsQuest() {
				data = fmt.Sprintf("quests %d", cell.ID)
			} else if cell.IsChat() {
				data = fmt.Sprintf("chat %d", cell.ID)
			}
			switch i {
			case 0:
				btn.Up += cell.Item.View
				btn.UpData = data
			case 1:
				btn.Down += cell.Item.View
				btn.DownData = data
			case 2:
				btn.Right += cell.Item.View
				btn.RightData = data
			case 3:
				btn.Left += cell.Item.View
				btn.LeftData = data
			}
		case cell.IsItem() || cell.IsSwap():
			switch i {
			case 0:
				btn.Up, btn.UpData = cell.ViewItemButton(resUser)
			case 1:
				btn.Down, btn.DownData = cell.ViewItemButton(resUser)
			case 2:
				btn.Right, btn.RightData = cell.ViewItemButton(resUser)
			case 3:
				btn.Left, btn.LeftData = cell.ViewItemButton(resUser)
			}
		case cell.IsWordleGame():
			switch i {
			case 0:
				btn.UpData = fmt.Sprintf("wordle_game up %s", cell.View)
				btn.Up = fmt.Sprintf("%s%s", btn.Up, cell.View)
			case 1:
				btn.DownData = fmt.Sprintf("wordle_game down %s", cell.View)
				btn.Down = fmt.Sprintf("%s%s", btn.Down, cell.View)
			case 2:
				btn.RightData = fmt.Sprintf("wordle_game right %s", cell.View)
				btn.Right = fmt.Sprintf("%s%s", btn.Right, cell.View)
			case 3:
				btn.LeftData = fmt.Sprintf("wordle_game left %s", cell.View)
				btn.Left = fmt.Sprintf("%s%s", btn.Left, cell.View)
			}
		case cell.ID == 0:
			switch i {
			case 0:
				btn.Up, btn.UpData = "🚫", "🚫"
			case 1:
				btn.Down, btn.DownData = "🚫", "🚫"
			case 2:
				btn.Right, btn.RightData = "🚫", "🚫"
			case 3:
				btn.Left, btn.LeftData = "🚫", "🚫"
			}
		default:
			switch i {
			case 0:
				btn.UpData = fmt.Sprintf("move %d", cell.ID)
			case 1:
				btn.DownData = fmt.Sprintf("move %d", cell.ID)
			case 2:
				btn.RightData = fmt.Sprintf("move %d", cell.ID)
			case 3:
				btn.LeftData = fmt.Sprintf("move %d", cell.ID)
			}
		}
	}

	return btn
}

func CreateMapKeyboard(buttons MapButtons) tg.InlineKeyboardMarkup {
	nearUsers := "Это не кнопка"

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Вещи 🧥", "goodsMoving"),
			tg.NewInlineKeyboardButtonData(buttons.Up, buttons.UpData),
			tg.NewInlineKeyboardButtonData("Рюкзак 🎒", "category"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(buttons.Left, buttons.LeftData),
			tg.NewInlineKeyboardButtonData(buttons.Center, buttons.Center),
			tg.NewInlineKeyboardButtonData(buttons.Right, buttons.RightData),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(nearUsers, nearUsers),
			tg.NewInlineKeyboardButtonData(buttons.Down, buttons.DownData),
			tg.NewInlineKeyboardButtonData(v.GetString("user_location.menu"), v.GetString("user_location.menu")),
		),
	)
}

func GetStatsLine(user User) string {
	var hpRes string
	var stRes string
	hp := fmt.Sprintf("%d", user.Health)
	st := fmt.Sprintf("%d", user.Satiety)

	length := len(hp)
	for i := 0; i < length; i++ {
		point := string(hp[i])
		hpRes += fmt.Sprintf("%s⃣", point)
	}
	hpRes += "♥️"

	length = len(st)
	stRes += "\U0001F9C3"
	for i := 0; i < length; i++ {
		point := string(st[i])
		stRes += fmt.Sprintf("%s⃣", point)
	}

	return fmt.Sprintf("%s   🅱️ *E T* 🅰️   %s\n", hpRes, stRes)
}
