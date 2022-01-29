package migrations

import (
	"fmt"
	"project0/config"
	"project0/repository"
)

func Migrate() bool {
	err := config.Db.AutoMigrate(repository.User{})
	if err != nil {
		fmt.Println("User Migration failed")
		panic(err)
		return false
	}
	err = config.Db.AutoMigrate(repository.Location{})
	if err != nil {
		fmt.Println("Location Migration failed")
		panic(err)
		return false
	}
	return true
}
