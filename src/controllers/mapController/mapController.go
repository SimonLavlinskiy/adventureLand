package mapController

import (
	"fmt"
	"project0/src/models"
)

func CalculateNightMap(user models.User, l models.Location, x int, y int) bool {

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

func CalculateSuperNightMap(user models.User, l models.Location, x int, y int) bool {

	if (*l.AxisX == x-1 || *l.AxisX == x+1) && (*l.AxisY >= y-1 && *l.AxisY <= y+1) {
		return true
	}
	if (*l.AxisX >= x-1 && *l.AxisX <= x+1) && (*l.AxisY >= y-1 && *l.AxisY <= y+1) {
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

func appendVisibleUserZone(m map[[2]int]models.Cell, x int, y int, Maps [][]string, loc models.Location, user models.User) [][]string {
	type Point [2]int
	var box models.UserBox

	if m[Point{x, y}].IsBox(user) {
		box = models.UserBox{UserId: loc.UserID, BoxId: m[Point{x, y}].Item.ID}
		if !box.IsUserGotBoxToday() {
			Maps[y] = append(Maps[y], m[Point{x, y}].Item.View)
		} else {
			Maps[y] = append(Maps[y], m[Point{x, y}].View)
		}
	} else if m[Point{x, y}].IsItem() || m[Point{x, y}].IsWorkbench() || m[Point{x, y}].IsSwap() || m[Point{x, y}].IsQuest() || m[Point{x, y}].IsChat() {
		Maps[y] = append(Maps[y], m[Point{x, y}].Item.View)
	} else {
		Maps[y] = append(Maps[y], m[Point{x, y}].View)
	}

	return Maps
}

func GetStatsLine(user models.User) string {
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
	if user.Health < 10 {
		hpRes = "❗️" + hpRes
	}

	length = len(st)
	stRes += "\U0001F9C3"
	for i := 0; i < length; i++ {
		point := string(st[i])
		stRes += fmt.Sprintf("%s⃣", point)
	}

	if user.Satiety < 10 {
		stRes += "❗️"
	}

	return fmt.Sprintf("%s   🅱️ *E T* 🅰️   %s\n", hpRes, stRes)
}
