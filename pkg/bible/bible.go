package bible

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/AidenHadisi/MyDailyBibleBot/assets"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/cache"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/httpclient"
)

type BibleAPI struct {
	client httpclient.HttpClient
	cache  cache.Cache
	topics []string
	verses map[string][]string
}

type BibleApiResult struct {
	Text string
}

func NewBibleAPI(client httpclient.HttpClient, cache cache.Cache) *BibleAPI {
	return &BibleAPI{
		client: client,
		cache:  cache,
		topics: make([]string, 0),
		verses: make(map[string][]string),
	}
}

//Init initializes the api client.
func (b *BibleAPI) Init() error {
	err := json.Unmarshal(assets.Topics, &b.verses)
	if err != nil {
		return err
	}

	for k := range b.verses {
		b.topics = append(b.topics, k)
	}

	return nil
}

//GetRandomVerse gets a random verse from API.
func (b *BibleAPI) GetRandomVerse() (string, error) {
	randomTopic := b.topics[rand.Intn(len(b.topics))]
	randomVerse := b.verses[randomTopic][rand.Intn(len(b.verses[randomTopic]))]
	resp, err := b.GetVerse(randomVerse)
	if err != nil {
		return "", err
	}
	return resp, nil
}

//GetVerse gets the requested verses from the API.
func (b *BibleAPI) GetVerse(verse string) (string, error) {
	cached, err := b.cache.Get(verse)
	if err == nil {
		return cached, nil
	}
	url := fmt.Sprintf("https://bible-api.com/%s?translation=kjv", verse)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := b.client.Do(req)
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

	b.cache.Set(verse, text, time.Hour*2)
	return text, nil
}
