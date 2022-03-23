package main

type Messages struct {
	MsgSeparator string
	UserLocation struct {
		menu             string
		maps             string
		profile          string
		tasksMenuMessage string
		ItemCategories   struct {
			categoryTitle   string
			categories      string
			otherCategories string
			food            string
			material        string
			sprout          string
			furniture       string
			other           string
		}
	}
	CallbackChar struct {
		cancel          string
		category        string
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
		buyHome         string
		userDoneQuest   string
		userGetQuest    string
		quest           string
		quests          string
		joinToChat      string
	}
	Elem struct {
		colors        string
		black         string
		brown         string
		red           string
		purple        string
		orange        string
		yellow        string
		green         string
		blue          string
		white         string
		door          string
		dayWindow     string
		eveningWindow string
		nightWindow   string
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
			chat            string
		}
	}
	Errors struct {
		noQuestItem          string
		userNotHasItemInHand string
		userHasOtherItem     string
		userDidNotTask       string
	}
	QuestStatuses struct {
		completed      string
		processed      string
		new            string
		completedEmoji string
		processedEmoji string
		newEmoji       string
	}
	MainInfo struct {
		costOfHouse int
	}
}
