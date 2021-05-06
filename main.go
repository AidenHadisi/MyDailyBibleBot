package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"githu.com/AidenHadisi/MyDailyBibleBot/bot"
	"github.com/dghubble/go-twitter/twitter"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	auth := &bot.Auth{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		AccessToken:    os.Getenv("ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("ACCESS_SECRET"),
	}

	bibleBot, err := bot.CreateBot(auth)
	if err != nil {
		log.Fatalf("Error %s", err)
	}

	_, _, err = bibleBot.TwitterClient.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	})

	if err != nil {
		log.Fatalf("Error %s", err)
	}

	bibleBot.Start()



}
