package bible

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AidenHadisi/MyDailyBibleBot/pkg/cache"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/httpclient"
)

type BibleAPI struct {
	client httpclient.HttpClient
	cache  cache.Cache
}

type BibleApiResult struct {
	Text string
}

func NewBibleAPI(client httpclient.HttpClient, cache cache.Cache) *BibleAPI {
	return &BibleAPI{
		client: client,
		cache:  cache,
	}
}

//GetVerse gets the requested verses from the API
func (bot *BibleAPI) GetVerse(verse string) (string, error) {
	cached, err := bot.cache.Get(verse)
	if err == nil {
		return cached, nil
	}
	url := fmt.Sprintf("https://bible-api.com/%s?translation=kjv", verse)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := bot.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api returned code %d", resp.StatusCode)
	}

	bibleApiResult := &BibleApiResult{}
	err = json.NewDecoder(resp.Body).Decode(bibleApiResult)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(bibleApiResult.Text, "\n", " "), verse)

	bot.cache.Set(verse, text, time.Hour*2)
	return text, nil
}
