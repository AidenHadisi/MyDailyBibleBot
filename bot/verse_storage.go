package bot

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

var topics []string
var verses map[string][]string

func LoadVerses() error {
	byteValue, err := ioutil.ReadFile("topics.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &verses)
	if err != nil {
		return err
	}

	for k := range verses {
		topics = append(topics, k)
	}

	return nil
}

func GetRandomVerse() string {
	randomTopic := topics[rand.Intn(len(topics))]
	randomVerse := verses[randomTopic][rand.Intn(len(verses[randomTopic]))]
	return randomVerse
}
