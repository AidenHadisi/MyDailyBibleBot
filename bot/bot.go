package bot

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

var verses = []string{"Mark 8:36", "Proverbs 21:21", "2 Corinthians 5:7", "Proverbs 27:19", "Ephesians 5:15-16", "Psalm 31:3", "Proverbs 4:23", "Psalm 25:4", "Ecclesiastes 3:1", "John 6:35", "Proverbs 13:3", "John 7:38",
	"Proverbs 19:8", "Proverbs 10:17", "Matthew 16:25", "Matthew 6:34", "1 John 4:9", "Psalm 118:24", "1 Corinthians 15:22", "Luke 11:28", "Matthew 10:39", "John 14:6", "Philippians 1:21", "Matthew 5:14", "Galatians 5:25",
	"Romans 14:8", "Psalm 119:93", "Revelation 3:19", "1 John 5:12", "Psalm 54:4", "Proverbs 14:12", "Romans 8:6", "John 10:10", "James 2:17", "2 Corinthians 10:3", "Acts 17:28", "Psalm 119:1", "1 Corinthians 10:31",
	"Psalm 34:22", "Romans 12:18", "Psalm 24:1", "Amos 5:4", "Matthew 3:8", "1 Peter 2:16", "Luke 9:24", "John 1:3", "Leviticus 20:26", "2 Timothy 3:12", "John 6:57"}

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
	gocron.Every(1).Hour().From(gocron.NextTick()).Do(bot.hourlyPost)
	<-gocron.Start()
	return bot, nil
}

func (bot *Bot) hourlyPost() {
	randomVerse := verses[rand.Intn(len(verses))]
	response, err := bot.GetVerse(randomVerse)
	if err != nil {
		log.Println(err)
		return
	}

	reply := fmt.Sprintf("\"%s\" - %s", response.Text, randomVerse)

	bot.TwitterClient.Statuses.Update(reply, nil)
}
