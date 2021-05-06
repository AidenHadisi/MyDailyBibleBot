package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/jasonlvhit/gocron"
	"github.com/patrickmn/go-cache"
)

const (
	baseURL     = "https://bible-api.com/"
	botUsername = "MyDailyBibleBot"
)

//Auth defines a struct twitter auth tokens
type Auth struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

//Bot defines MyDailyBibleBot structure
type Bot struct {
	TwitterClient *twitter.Client
	Auth          *Auth
	HTTPClient    httpClient
	cache         *cache.Cache
	verses        []string
}

//CreateBot created a new instance of MyDailyBibleBot
func CreateBot(auth *Auth) (*Bot, error) {
	config := oauth1.NewConfig(auth.ConsumerKey, auth.ConsumerSecret)
	token := oauth1.NewToken(auth.AccessToken, auth.AccessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	var bot *Bot = new(Bot)
	bot.Auth = auth

	bot.TwitterClient = twitter.NewClient(httpClient)
	bot.HTTPClient = &http.Client{
		Timeout: time.Second,
	}
	bot.cache = cache.New(time.Hour, 15*time.Minute)

	jsonFile, err := os.Open("./verses.json")
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &bot.verses)
	if err != nil {
		return nil, err
	}

	gocron.Every(1).Hour().From(gocron.NextTick()).Do(bot.hourlyPost)
	gocron.Start()

	return bot, nil
}

func (bot *Bot) hourlyPost() {
	randomVerse := bot.verses[rand.Intn(len(bot.verses))]
	response, err := bot.GetVerse(randomVerse)
	if err != nil {
		log.Println(err)
		return
	}

	reply := fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(response.Text, "\n", " "), randomVerse)

	_, _, err = bot.TwitterClient.Statuses.Update(reply, nil)
	if err != nil {
		log.Println(err)
	}
}
