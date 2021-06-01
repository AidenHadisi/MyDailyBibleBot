package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	request "github.com/AidenHadisi/go-simple-request"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/jasonlvhit/gocron"
	"github.com/patrickmn/go-cache"
)

const botUsername = "MyDailyBibleBot"

var verses []string

//Auth defines a struct for Twitter auth tokens
type Auth struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

//Bot defines MyDailyBibleBot structure
type Bot struct {
	TwitterClient *twitter.Client
	Auth          *Auth
	req           *request.Request
	cache         *cache.Cache
}

//CreateBot created a new instance of MyDailyBibleBot
func CreateBot(auth *Auth) (*Bot, error) {

	config := oauth1.NewConfig(auth.ConsumerKey, auth.ConsumerSecret)
	token := oauth1.NewToken(auth.AccessToken, auth.AccessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	var bot *Bot = new(Bot)
	bot.Auth = auth

	bot.TwitterClient = twitter.NewClient(httpClient)
	bot.cache = cache.New(time.Hour, 15*time.Minute)

	byteValue, err := ioutil.ReadFile("bot/verses.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &verses)
	if err != nil {
		return nil, err
	}

	bot.req = request.New().SetFailure(&ErrResponse{}).SetSuccess(&Response{})

	gocron.Every(5).Hour().From(gocron.NextTick()).Do(bot.hourlyPost)
	gocron.Start()

	return bot, nil
}

func (bot *Bot) hourlyPost() {
	randomVerse := verses[rand.Intn(len(verses))]
	resp, err := bot.GetVerse(randomVerse, &BibleOptions{Translation: "kjv"})
	if err != nil {
		log.Println(err)
		return
	}

	reply := fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(resp.Text, "\n", " "), randomVerse)

	_, _, err = bot.TwitterClient.Statuses.Update(reply, nil)
	if err != nil {
		log.Println(err)
	}
}
