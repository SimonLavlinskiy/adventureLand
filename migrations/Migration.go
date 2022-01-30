package migrations

import (
	"fmt"
	"project0/config"
	"project0/repository"
)

func Migrate() bool {
	err := config.Db.AutoMigrate(repository.User{}, repository.Location{})
	if err != nil {
		fmt.Println("User Migration failed")
		panic(err)
		return false
	}
	return true
}
