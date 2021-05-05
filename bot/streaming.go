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

	verse := fmt.Sprintf("%s %s:%s", parsed.Book, parsed.Chapter, parsed.Start)

	if parsed.IsMultiVerse() {
		verse = fmt.Sprintf("%s-%s", verse, parsed.End)
	}

	go bot.reply(tweet, verse)
}

func (bot *Bot) reply(tweet *twitter.Tweet, verse string) {
	text, err := bot.fetch(verse)

	if err != nil {
		log.Println(err)
		return
	}

	reponseParts := splitter.Split(text, 260)

	t := tweet
	for index, part := range reponseParts {
		var message string
		//If intial tweet @ the sender, otherwise reply to previous bot message to create a thread
		if index == 0 {
			message = fmt.Sprintf("@%s %s", t.User.ScreenName, part)
		} else {
			message = fmt.Sprintf("@%s %s", botUsername, part)
		}
		t, _, err = bot.TwitterClient.Statuses.Update(message, &twitter.StatusUpdateParams{
			InReplyToStatusID: t.ID,
		})

		if err != nil {
			return
		}

	}
}

func (bot *Bot) fetch(verse string) (string, error) {
	cached, found := bot.cache.Get(verse)

	if found {
		return cached.(string), nil
	}

	response, err := bot.GetVerse(verse)
	if err != nil {
		return "", err
	}

	text := fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(response.Text, "\n", ""), verse)

	bot.cache.Set(verse, text, cache.DefaultExpiration)

	return text, nil

}
