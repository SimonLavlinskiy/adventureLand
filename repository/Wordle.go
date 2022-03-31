package repository

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	"time"
)

type Word struct {
	ID         uint   `gorm:"primaryKey"`
	SecretWord string `gorm:"embedded"`
	Date       string `gorm:"embedded"`
}

type WordleGameProcess struct {
	ID         uint `gorm:"primaryKey"`
	UserId     uint `gorm:"embedded"`
	User       User
	CountTries int       `gorm:"embedded"`
	Status     string    `gorm:"embedded"`
	Date       time.Time `gorm:"autoCreateTime"`
}

func GetActiveWord() (*Word, error) {

	currentDate := time.Now().Format("2006-01-02")
	result := Word{}
	err := config.Db.Where(Word{Date: currentDate}).First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetOrCreateWordleGameProcess(user User) WordleGameProcess {
	today := time.Now().Format("2006-01-02")

	result := WordleGameProcess{
		UserId:     user.ID,
		CountTries: 0,
		Status:     "new",
	}

	config.Db.
		Where(&WordleGameProcess{UserId: user.ID}).
		Where(fmt.Sprintf("date like '%s%s'", today, "%")).
		FirstOrCreate(&result)

	return result
}

func GetWordleGameProcessOfUser(user User, status *string) []WordleGameProcess {
	var result []WordleGameProcess
	if status != nil {
		config.Db.
			Where(&WordleGameProcess{UserId: user.ID, Status: *status}).
			Find(&result)
	} else {
		config.Db.
			Where(&WordleGameProcess{UserId: user.ID}).
			Find(&result)
	}

	return result
}

func (w WordleGameProcess) UpdateWordleGameProcess(user User) {
	today := time.Now().Format("2006-01-02")

	if w.CountTries < 5 && w.Status == "new" {
		w.CountTries++
	} else if w.CountTries == 5 && w.Status == "new" {
		w.CountTries++
		w.Status = "lose"
	}

	config.Db.
		Where(&WordleGameProcess{UserId: user.ID}).
		Where(fmt.Sprintf("date like '%s%s'", today, "%")).
		Updates(WordleGameProcess{Status: w.Status, CountTries: w.CountTries})
}

func GetWordleUserStatistic(user User) string {
	statusWin := "win"
	statusLose := "lose"
	statusNew := "new"

	games := GetWordleGameProcessOfUser(user, nil)
	wonGames := GetWordleGameProcessOfUser(user, &statusWin)
	loseGames := GetWordleGameProcessOfUser(user, &statusLose)
	newGames := GetWordleGameProcessOfUser(user, &statusNew)

	title := "📊 Статистика 📊"
	totalGames := fmt.Sprintf("*\U0001F9E9 Всего сыграно игр*: _%d игр_", len(games))
	totalWonGames := fmt.Sprintf("*🥇 Всего выиграно игр*: %d", len(wonGames))
	totalLoseGames := fmt.Sprintf("*\U0001F97A Всего проиграно игр*: %d", len(loseGames))
	totalNewGames := fmt.Sprintf("*🏳️ Всего неокончено игр*: %d", len(newGames))

	return fmt.Sprintf("%s%s%s\n%s\n%s\n%s", title, v.GetString("msg_separator"), totalGames, totalWonGames, totalLoseGames, totalNewGames)
}
