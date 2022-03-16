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

type MapButtons struct {
	Up     string
	Left   string
	Right  string
	Down   string
	Center string
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

type UserMap struct {
	leftIndent  int
	rightIndent int
	upperIndent int
	downIndent  int
}

var displaySize = 6

func DefaultUserMap(location Location, displaySize int) UserMap {
	return UserMap{
		leftIndent:  *location.AxisX - displaySize,
		rightIndent: *location.AxisX + displaySize,
		upperIndent: *location.AxisY + displaySize,
		downIndent:  *location.AxisY - displaySize,
	}
}

func GetUserMap(update tg.Update) Map {
	resLocation := GetOrCreateMyLocation(update)
	result := Map{}

	err := config.Db.Where(&Map{ID: uint(*resLocation.MapsId)}).First(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetMaps() []Map {
	var result []Map

	err := config.Db.Find(&result).Error

	if err != nil {
		fmt.Println("Нет карт!")
	}

	return result
}

func GetMyMap(update tg.Update) (textMessage string, buttons tg.ReplyKeyboardMarkup) {
	userTgId := GetUserTgId(update)
	us := GetUser(User{TgId: userTgId})
	loc := GetOrCreateMyLocation(update)
	resMap := GetUserMap(update)
	mapSize := CalculateUserMapBorder(loc, resMap)
	messageMap := fmt.Sprintf("*Карта*: _%s_ *X*: _%d_  *Y*: _%d_", loc.Maps.Name, *loc.AxisX, *loc.AxisY)

	type Point = [2]int
	m := map[Point]Cell{}

	var resultCell []Cell
	UpdateCellWithNextStateTime()

	err := config.Db.
		Preload("Item").
		Preload("Teleport").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.Result").
		Preload("Item.Instruments.NextStageItem").
		Where(Cell{MapsId: *loc.MapsId}).
		Where("axis_x >= " + ToString(mapSize.leftIndent)).
		Where("axis_x <= " + ToString(mapSize.rightIndent)).
		Where("axis_y >= " + ToString(mapSize.downIndent)).
		Where("axis_y <= " + ToString(mapSize.upperIndent)).
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
	for _, user := range resultLocationsOnlineUser {
		if user.User.ID != 0 && *user.User.OnlineMap {
			m[Point{*user.AxisX, *user.AxisY}] = Cell{View: user.User.Avatar, ID: m[Point{*user.AxisX, *user.AxisY}].ID}
		}
	}

	buttons = CalculateButtonMap(loc, us, m)

	m[Point{*loc.AxisX, *loc.AxisY}] = Cell{View: us.Avatar, ID: m[Point{*loc.AxisX, *loc.AxisY}].ID}

	Maps := configurationMap(mapSize, resMap, loc, us, m)

	for i, row := range Maps {
		if i >= mapSize.downIndent && i <= mapSize.upperIndent {
			messageMap = fmt.Sprintf("%s\n%s", strings.Join(row, ""), messageMap)
		}
	}

	return messageMap, buttons
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
		mapSize.leftIndent = 0
		mapSize.rightIndent = d * 2
	}
	if *resLocation.AxisY < d {
		mapSize.downIndent = 0
		mapSize.upperIndent = d * 2
	}
	if mapSize.rightIndent >= resMap.SizeX && *resLocation.AxisX > d {
		mapSize.leftIndent = resMap.SizeX - d*2
		mapSize.rightIndent = resMap.SizeX
	}
	if mapSize.upperIndent >= resMap.SizeY && *resLocation.AxisY > d {
		mapSize.downIndent = resMap.SizeY - d*2
		mapSize.upperIndent = resMap.SizeY
	}

	return mapSize
}

func CalculateButtonMap(resLocation Location, resUser User, m map[[2]int]Cell) tg.ReplyKeyboardMarkup {
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

func CreateMapKeyboard(buttons MapButtons) tg.ReplyKeyboardMarkup {
	// Цэ кнопка была включения онлайна и офлайна
	//onlineButton := v.GetString("message.user.online")
	//if *user.OnlineMap {
	//	onlineButton = v.GetString("message.user.offline")
	//}

	nearUsers := "👥"

	return tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Вещи 🧥"),
			tg.NewKeyboardButton(buttons.Up),
			tg.NewKeyboardButton("Рюкзак 🎒"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton(buttons.Left),
			tg.NewKeyboardButton(buttons.Center),
			tg.NewKeyboardButton(buttons.Right),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton(nearUsers), //onlineButton),
			tg.NewKeyboardButton(buttons.Down),
			tg.NewKeyboardButton(v.GetString("user_location.menu")),
		),
	)
}

func configurationMap(mapSize UserMap, resMap Map, resLocation Location, user User, m map[[2]int]Cell) [][]string {
	currentTime := time.Now()
	day := "06" <= currentTime.Format("15") && currentTime.Format("15") <= "23"
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
		for x := mapSize.leftIndent; x <= mapSize.rightIndent; x++ {
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
		for x := mapSize.leftIndent; x <= mapSize.rightIndent; x++ {
			if calculateNightMap(user, resLocation, x, y) && m[Point{x, y}].ID != 0 {
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

func calculateNightMap(user User, l Location, x int, y int) bool {

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

func appendVisibleUserZone(m map[[2]int]Cell, x int, y int, Maps [][]string) {
	type Point [2]int

	if m[Point{x, y}].IsItem() || m[Point{x, y}].IsWorkbench() || m[Point{x, y}].IsSwap() || m[Point{x, y}].IsQuest() {
		Maps[y] = append(Maps[y], m[Point{x, y}].Item.View)
	} else {
		Maps[y] = append(Maps[y], m[Point{x, y}].View)
	}
}

func PutButton(CellsAroundUser []Cell, btn MapButtons, resUser User) MapButtons {

	for i, cell := range CellsAroundUser {
		switch true {
		case cell.IsDefaultCell():
			switch i {
			case 0:
				btn.Up = cell.View
			case 1:
				btn.Down = cell.View
			case 2:
				btn.Right = cell.View
			case 3:
				btn.Left = cell.View
			}
		case cell.IsTeleport():
			button := fmt.Sprintf(" %s %s", resUser.Avatar, cell.View)
			switch i {
			case 0:
				btn.Up += button
			case 1:
				btn.Down += button
			case 2:
				btn.Right += button
			case 3:
				btn.Left += button
			}
		case cell.IsWorkbench():
			switch i {
			case 0:
				btn.Up = fmt.Sprintf("🔧 %s %s", btn.Up, cell.Item.View)
			case 1:
				btn.Down = fmt.Sprintf("🔧 %s %s", btn.Down, cell.Item.View)
			case 2:
				btn.Right = fmt.Sprintf("🔧 %s %s", btn.Right, cell.Item.View)
			case 3:
				btn.Left = fmt.Sprintf("🔧 %s %s", btn.Left, cell.Item.View)
			}
		case cell.IsItem() || cell.IsSwap():
			switch i {
			case 0:
				btn.Up = cell.isItemCost(btn.Up, resUser)
			case 1:
				btn.Down = cell.isItemCost(btn.Down, resUser)
			case 2:
				btn.Right = cell.isItemCost(btn.Right, resUser)
			case 3:
				btn.Left = cell.isItemCost(btn.Left, resUser)
			}
		case cell.IsQuest():
			switch i {
			case 0:
				btn.Up = fmt.Sprintf("📟 %s %s", btn.Up, cell.Item.View)
			case 1:
				btn.Down = fmt.Sprintf("📟 %s %s", btn.Down, cell.Item.View)
			case 2:
				btn.Right = fmt.Sprintf("📟 %s %s", btn.Right, cell.Item.View)
			case 3:
				btn.Left = fmt.Sprintf("📟 %s %s", btn.Left, cell.Item.View)
			}
		case cell.ID == 0:
			switch i {
			case 0:
				btn.Up = "🚫"
			case 1:
				btn.Down = "🚫"
			case 2:
				btn.Right = "🚫"
			case 3:
				btn.Left = "🚫"
			}
		}
	}

	return btn
}

func (c Cell) IsDefaultCell() bool {
	if c.Type != nil && *c.Type == "cell" && !c.CanStep {
		return true
	}
	return false
}

func (c Cell) IsWorkbench() bool {
	if c.Type != nil && *c.Type == "workbench" && c.ItemID != nil {
		return true
	}
	return false
}

func (c Cell) IsTeleport() bool {
	if c.Type != nil && *c.Type == "teleport" && c.TeleportID != nil {
		return true
	}
	return false
}

func (c Cell) IsItem() bool {
	if c.Type != nil && *c.Type == "item" && c.ItemID != nil && *c.ItemCount > 0 {
		return true
	}
	return false
}

func (c Cell) IsSwap() bool {
	if c.Type != nil && *c.Type == "swap" && c.ItemID != nil {
		return true
	}
	return false
}

func (c Cell) IsQuest() bool {
	if c.Type != nil && *c.Type == "quest" && c.ItemID != nil {
		return true
	}
	return false
}

func (c Cell) IsSpecialItem(user User) string {
	instrumentsUserCanUse := GetInstrumentsUserCanUse(user, c)

	if len(instrumentsUserCanUse) > 1 {
		return "❗ 🛠 ❓"
	} else if len(instrumentsUserCanUse) == 1 {
		for i := range instrumentsUserCanUse {
			return i
		}
	}

	return "🚷"
}

func (c Cell) isItemCost(button string, resUser User) string {
	var firstElem = c.IsSpecialItem(resUser)

	button = firstElem + " " + button + " " + c.Item.View

	if c.Item.Cost != nil && *c.Item.Cost > 0 && firstElem != "❗ 🛠 ❓" {
		button = button + " ( " + ToString(*c.Item.Cost) + "💰 )"
	}

	return button
}
