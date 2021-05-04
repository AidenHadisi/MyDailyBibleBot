package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AidenHadisi/go-text-splitter/splitter"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/patrickmn/go-cache"
)

// Start starts the bot and begins listening to Tweets
func (bot *Bot) Start() {
	params := &twitter.StreamFilterParams{
		Track:         []string{fmt.Sprintf("@%s", botUsername)},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := bot.TwitterClient.Streams.Filter(params)
	if err != nil {
		log.Fatal(err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = bot.handleMessage
	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	stream.Stop()
}

func (bot *Bot) handleMessage(tweet *twitter.Tweet) {
	if tweet.User.ScreenName == botUsername {
		return
	}

	parsed, err := ParseText(tweet.Text)

	if err != nil || !parsed.IsValid() {
		log.Println(err)
		return
	}

	verse := fmt.Sprintf("%s %s:%s", strings.Title(strings.ToLower(parsed.Book)), parsed.Chapter, parsed.Start)

	if parsed.IsMultiVerse() {
		verse = fmt.Sprintf("%s-%s", verse, parsed.End)
	}

	go bot.fetch(tweet, verse)
}

func (bot *Bot) fetch(tweet *twitter.Tweet, verse string) {
	cached, found := bot.cache.Get(verse)

	if found {
		reply := fmt.Sprintf("@%s \"%s\" - %s", tweet.User.ScreenName, cached, verse)
		bot.textToTweet(reply, tweet)
		return
	}

	response, err := bot.GetVerse(verse)
	if err != nil {
		log.Println(err)
		return
	}

	reply := fmt.Sprintf("@%s \"%s\" - %s", tweet.User.ScreenName, strings.TrimSuffix(response.Text, "\n"), verse)

	reponseParts := splitter.Split(reply, 260)

	if len(reponseParts) > 3 {
		return
	}

	bot.textToTweet(reply, tweet)

	bot.cache.Set(verse, strings.TrimSuffix(response.Text, "\n"), cache.DefaultExpiration)
}
