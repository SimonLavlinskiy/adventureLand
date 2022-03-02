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
		repository.Cellule{},
		repository.Map{},
		repository.UserItem{},
		repository.Instrument{},
		repository.Receipt{},
	)
	if err != nil {
		fmt.Println("Migration failed")
		panic(err)
		return false
	}
	//err = config.Db.SetupJoinTable(&repository.Item{}, "Instruments", &repository.InstrumentItem{})

	return true
}
