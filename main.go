package main

import (
	"log"
	"os"

	"github.com/AidenHadisi/MyDailyBibleBot/bot"
	"github.com/dghubble/go-twitter/twitter"
)

func main() {
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

	// res, err := bibleBot.GetVerse("/romans 1000:")
	// if err != ni {
	// 	log.Printf("Error %s", er)
	// }

	// fmt.Printf("%s", res.Errr)

	// _, _, err = bot.TwitterClient.Statuses.Update("just setting up my twttr", nil)
	// if err != nil {
	// 	log.Printf("%s\n", err)
	// }

}
