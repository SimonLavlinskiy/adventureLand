package instrumentServices

import (
	"project0/src/models"
)

func GetInstrumentsUserCanUse(user models.User, cell models.Cell) map[string]models.Item {
	instrumentsUserCanUse := map[string]models.Item{}

	if cell.Item == nil {
		return instrumentsUserCanUse
	}
	instruments := cell.Item.Instruments

	for _, instrument := range instruments {
		if instrument.GoodId != nil && user.LeftHandId != nil && user.LeftHand.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.LeftHand.View] = *user.LeftHand
		}
		if instrument.GoodId != nil && user.RightHandId != nil && user.RightHand.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.RightHand.View] = *user.RightHand
		}
		if instrument.GoodId != nil && user.HeadId != nil && user.Head.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.Head.View] = *user.Head
		}
		if instrument.GoodId != nil && instrument.Good.Type == "fist" {
			instrumentsUserCanUse["🤜"] = *instrument.Good
		}
		if *cell.Type == "item" && instrument.Type == "goToSleep" {
			instrumentsUserCanUse["💤"] = models.Item{Type: "sleep"}
		}
	}

	if cell.Item.CanTake {
		instrumentsUserCanUse["👋"] = models.Item{Type: "hand"}
	}
	if cell.Item.CanStep && *cell.Type != "swap" {
		instrumentsUserCanUse["👣"] = models.Item{Type: "step"}
	}

	return instrumentsUserCanUse
}
