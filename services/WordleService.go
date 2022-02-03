package services

import (
	"fmt"
	"strings"
)

var resultArray []string

func StringToArrayCharacters(word string) []rune {

	chars := []rune(word)
	for i := 0; i < len(chars); i++ {
		char := string(chars[i])
		fmt.Println(char)
	}
	return chars
}

func CheckWordEquals(userWord string, secretWord string) []string {

	userArray := StringToArrayCharacters(userWord)
	secretArray := StringToArrayCharacters(secretWord)
	for i, v := range secretArray {
		if v != userArray[i] {
			if strings.ContainsRune(secretWord, userArray[i]) {
				resultArray[i] = "🟨"
			} else {
				resultArray[i] = "\U0001F7E5"

			}
		} else {
			resultArray[i] = "\U0001F7E9"
		}
	}

	return resultArray
}

func FormattedUserWord(userWord string) string {

	res := strings.ToLower(userWord)
	res = strings.TrimSpace(res)

	return res

}
