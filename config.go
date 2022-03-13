package main

type Messages struct {
	MsgSeparator string
	UserLocation struct {
		menu             string
		maps             string
		profile          string
		tasksMenuMessage string
	}
	CallbackChar struct {
		cancel          string
		backpackMoving  string
		goodsMoving     string
		eatFood         string
		deleteItem      string
		dressGood       string
		takeOffGood     string
		changeLeftHand  string
		changeRightHand string
		changeAvatar    string
		description     string
		workbench       string
		receipt         string
		putItem         string
		putCountItem    string
		makeNewItem     string
		hand            string
		foot            string
		throwOutItem    string
		countOfDelete   string
	}
	Message struct {
		ListOfAvatar string
		Doing        struct {
			up    string
			down  string
			left  string
			right string
		}
		Emoji struct {
			water           string
			clock           string
			casino          string
			forbidden       string
			wrench          string
			offline         string
			online          string
			stopUse         string
			exclamationMark string
			hand            string
			foot            string
			quest           string
		}
	}
	Errors struct {
		noQuestItem          string
		userNotHasItemInHand string
		userHasOtherItem     string
	}
	QuestStatuses struct {
		completed string
		processed string
		new       string
	}
}
