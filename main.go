package main

import (
	"log"

	"github.com/AidenHadisi/MyDailyBibleBot/bot"
	"github.com/dghubble/go-twitter/twitter"
)

func main() {
	auth := &bot.Auth{
		ConsumerKey:    "2Llovmsx9eWcXyieYO3tlzZxj",
		ConsumerSecret: "aM7odZl49AVtzJR76PRhFdZWWVnVZW6Fqpxnf14Y1nGdpIgbss",
		AccessToken:    "1386538984680026116-Hu1MrJnJYTjlH5J3TrWpi5TLITabT9",
		AccessSecret:   "AcMxtxSb2EY6BgxYqk9u8mVApwbfE59f8ofndVp18nAYS",
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
