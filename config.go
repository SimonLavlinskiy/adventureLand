package main

type Configurations struct {
	MsgSeparator string
	UserLocation
	CallbackChar
	Message
}

type UserLocation struct {
	menu    string
	maps    string
	profile string
}

type CallbackChar struct {
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

type Message struct {
	Doing
	Emoji
}

type Doing struct {
	up    string
	down  string
	left  string
	right string
}

type Emoji struct {
	water            string
	clock            string
	casino           string
	forbidden        string
	wrench           string
	offline          string
	online           string
	stop_use         string
	exclamation_mark string
	hand             string
	foot             string
}
