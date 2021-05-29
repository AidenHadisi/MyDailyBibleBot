package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/google/go-querystring/query"
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

	byteValue, err := ioutil.ReadFile("bot/verses.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &verses)
	if err != nil {
		return nil, err
	}

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

func (bot *Bot) get(path string, reqData, success, failure interface{}) (*Response, error) {
	return bot.sendRequest(http.MethodGet, path, reqData, success, failure)
}

func (bot *Bot) sendRequest(method, path string, reqData, success, failaure interface{}) (*Response, error) {

	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}

	v, err := query.Values(reqData)
	if err == nil {
		req.URL.RawQuery = v.Encode()
	}

	resp, err := bot.do(req, success, failaure)

	return resp, err
}

func (bot *Bot) do(req *http.Request, success, failure interface{}) (*Response, error) {

	response := &Response{}
	resp, err := bot.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// if resp.StatusCode == http.StatusNoContent || resp.ContentLength == 0 {
	// 	return response, nil
	// }

	response.StatusCode = resp.StatusCode
	response.Header = resp.Header

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = decodeResp(response, bodyBytes, success, failure)

	if err != nil {
		err = fmt.Errorf("failed to decode API response: %s", err.Error())
	}

	return response, err
}

func decodeResp(resp *Response, body []byte, success, failure interface{}) error {
	if status := resp.StatusCode; 200 <= status && status <= 299 {
		if success != nil {
			resp.Success = success

			return json.Unmarshal(body, &resp.Success)
		}

	} else {
		if failure != nil {
			resp.Failure = failure
			return json.Unmarshal(body, &resp.Failure)
		}
	}
	return nil

}
