package migrations

import (
	"fmt"
	"project0/config"
	"project0/src/models"
)

func Migrate() bool {
	err := config.Db.AutoMigrate(
		models.User{},
		models.Location{},
		models.Teleport{},
		models.Item{},
		models.Cell{},
		models.Map{},
		models.UserItem{},
		models.Instrument{},
		models.Receipt{},
		models.Quest{},
		models.UserQuest{},
		models.QuestTask{},
		models.Result{},
		models.Chat{},
		models.ChatUser{},
		models.Word{},
		models.WordleGameProcess{},
		models.UserWords{},
		models.UserBox{},
		models.News{},
		models.UserSleep{},
	)
	if err != nil {
		fmt.Println("Migration failed")
		panic(err)
		return false
	}
	//err = config.Db.SetupJoinTable(&models.Item{}, "Instruments", &models.InstrumentItem{})

	return true
}
