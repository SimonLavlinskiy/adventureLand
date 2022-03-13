package migrations

import (
	"fmt"
	"project0/config"
	"project0/repository"
)

func Migrate() bool {
	err := config.Db.AutoMigrate(
		repository.User{},
		repository.Location{},
		repository.Teleport{},
		repository.Item{},
		repository.Cell{},
		repository.Map{},
		repository.UserItem{},
		repository.Instrument{},
		repository.Receipt{},
		repository.Quest{},
		repository.UserQuest{},
		repository.QuestTask{},
		repository.QuestResult{},
	)
	if err != nil {
		fmt.Println("Migration failed")
		panic(err)
		return false
	}
	//err = config.Db.SetupJoinTable(&repository.Item{}, "Instruments", &repository.InstrumentItem{})

	return true
}
