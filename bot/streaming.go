package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/patrickmn/go-cache"
)

// Start starts the bot and begins listening to Tweets
func (bot *Bot) Start() {
	params := &twitter.StreamFilterParams{
		Track:         []string{"@MyDailyBibleBot"},
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
	parsed, err := ParseText(tweet.Text)

	if err != nil || !parsed.IsValid() {
		log.Println(err)
		return
	}

	verse := fmt.Sprintf("%s %s:%s", strings.Title(strings.ToLower(parsed.Book)), parsed.Chapter, parsed.Start)

	if parsed.IsMultiVerse() {
		verse = fmt.Sprintf("%s-%s", verse, parsed.End)
	}

	go bot.fetchAndReply(tweet, verse)
}

func (bot *Bot) fetchAndReply(tweet *twitter.Tweet, verse string) {
	cached, found := bot.cache.Get(verse)

	if found {
		reply := fmt.Sprintf("@%s \"%s\" - %s", tweet.User.ScreenName, cached, verse)

		reponseParts := breakupString(reply, 280)

		t := tweet
		var err error
		for _, part := range reponseParts {
			t, _, err = bot.TwitterClient.Statuses.Update(part, &twitter.StatusUpdateParams{
				InReplyToStatusID: t.ID,
			})

			if err != nil {
				return
			}

		}
		bot.TwitterClient.Statuses.Update(reply, &twitter.StatusUpdateParams{
			InReplyToStatusID: tweet.ID,
		})
		return
	}

	response, err := bot.GetVerse(verse)
	if err != nil {
		log.Println(err)
		return
	}

	result := fmt.Sprintf("@%s \"%s\" - %s", tweet.User.ScreenName, strings.TrimSuffix(response.Text, "\n"), verse)

	reponseParts := breakupString(result, 280)

	if len(reponseParts) > 3 {
		return
	}

	t := tweet
	for _, part := range reponseParts {
		t, _, err = bot.TwitterClient.Statuses.Update(part, &twitter.StatusUpdateParams{
			InReplyToStatusID: t.ID,
		})

		if err != nil {
			return
		}

	}

	bot.cache.Set(verse, strings.TrimSuffix(response.Text, "\n"), cache.DefaultExpiration)
}
