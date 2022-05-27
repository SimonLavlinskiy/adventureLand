package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func IsDictionaryHasWord(word string) bool {
	key, _ := os.LookupEnv("YANDEX_DICTIONARY_KEY")
	resp, _ := http.Get(
		fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=ru-ru&text=%s", key, word),
	)

	r := Response{}
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &r)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if len(r.Def) > 0 {
		return true
	}

	return false
}

type Response struct {
	Head struct{}
	Def  []struct {
		Text string
		Pos  string
		Tr   []struct {
			Text string
			Pos  string
			Fr   int
			Syn  []struct {
				Text string
				Pos  string
				Fr   int
			}
		}
	}
}
