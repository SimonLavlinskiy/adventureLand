package services

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	r "project0/repository"
	"strings"
	"time"
)

func StringToArrayCharacters(word string) []rune {
	chars := []rune(word)
	return chars
}

func CheckWordEquals(userWord string) string {
	var resultString string
	activeWord, _ := r.GetActiveWord()

	userArray := StringToArrayCharacters(userWord)
	secretArray := StringToArrayCharacters(activeWord.SecretWord)
	for i, char := range secretArray {
		if char != userArray[i] {
			if strings.ContainsRune(activeWord.SecretWord, userArray[i]) {
				resultString += "🟨"
			} else {
				resultString += "⬜️"

			}
		} else {
			resultString += "✅"
		}
	}

	return resultString
}

func FormattedUserWord(userWord string) string {

	res := strings.ToLower(userWord)
	res = strings.TrimSpace(res)

	return res

}

func WordleMenuButtons(game r.WordleGameProcess) tg.InlineKeyboardMarkup {
	var button tg.InlineKeyboardButton

	switch game.Status {
	case "new":
		button = tg.NewInlineKeyboardButtonData("Написать слово", "awaitWord")
	case "lose", "win":
		button = tg.NewInlineKeyboardButtonData("Статистика", "wordleUserStatistic")
	}

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			button,
		), tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Правила", "wordleRegulations"),
			tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
		),
	)
}

func WordleMap(user r.User) []tg.MessageConfig {
	var msgs []tg.MessageConfig
	countTries := 6

	game := r.GetOrCreateWordleGameProcess(user)
	words := r.GetUserWords(user, time.Now())

	var msgText string

	msgText += "\U0001F9EE *Игра Вуордле!*\U0001F9EE\n"

	for i, word := range words {
		if i < countTries {
			row := CheckWordEquals(word.Word)
			msgText += fmt.Sprintf("\n%s - *%s*", row, strings.ToUpper(word.Word))
		}
	}

	if len(words) < 6 {
		for x := len(words); x < countTries; x++ {
			msgText += "\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6"
			if x == len(words) && game.Status == "new" {
				msgText += " - _Напиши слово!_"
			}
		}
	}

	var lastText string
	var msg tg.MessageConfig
	if game.Status == "new" && game.CountTries < countTries {
		lastText = "Только 5 букв! 👉🤚 Ни больше, ни меньше! 👌"
	} else if game.Status == "win" {
		lastText = "🏆 Поздравляю, ты выйграл! 🏆"
	} else if game.Status == "lose" {
		lastText = "☠️ Ты проиграл :C Ну ничего, попробуй завтра еще раз! 👋"
	}

	msg.ReplyMarkup = WordleMenuButtons(game)
	msg.Text = fmt.Sprintf("%s%s_%s_", msgText, v.GetString("msg_separator"), lastText)
	msgs = append(msgs, msg)

	return msgs
}

func CheckUserWordFormat(user r.User, userWord string) (tg.MessageConfig, error) {
	var msg tg.MessageConfig
	if len(strings.Fields(userWord)) != 1 {
		msg.Text = "Некорректное количество слов"
		return msg, errors.New("too many words")
	}

	userWord = FormattedUserWord(userWord)

	if chars := []rune(userWord); len(chars) > 5 {
		msg.Text = "Слишком много букв"
		return msg, errors.New("too many chars")
	} else if len(chars) < 5 {
		msg.Text = "Слишком мало букв"
		return msg, errors.New("not enough chars")
	}

	words := r.GetUserWords(user, time.Now())
	for _, word := range words {
		if word.Word == userWord {
			msg.Text = "Такое слово уже было"
			return msg, errors.New("word duplicate")
		}
	}

	return msg, nil
}

func UserSendNextWord(user r.User, newMessage string) []tg.MessageConfig {
	var msgs []tg.MessageConfig

	msg, err := CheckUserWordFormat(user, newMessage)
	if err != nil {
		msgs = WordleMap(user)
		msgs = append(msgs, msg)
		return msgs
	}

	game := r.GetOrCreateWordleGameProcess(user)

	if game.Status != "new" {
		msg.Text = "\U0001FAC2 Успокойся, игра уже окончена! Приходи завтра) 🤝"
		msgs = append(msgs, msg)
		return msgs
	}

	word := FormattedUserWord(newMessage)
	r.CreateUserWord(user, word)
	activeWord, _ := r.GetActiveWord()

	if word == activeWord.SecretWord {
		game.Status = "win"
	}

	game.UpdateWordleGameProcess(user)

	msgs = WordleMap(user)

	return msgs
}
