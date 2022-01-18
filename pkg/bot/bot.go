package bot

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AidenHadisi/MyDailyBibleBot/configs"
	"github.com/AidenHadisi/MyDailyBibleBot/internal/bible"
	"github.com/AidenHadisi/MyDailyBibleBot/internal/twitter"

	"github.com/AidenHadisi/go-text-splitter/splitter"
	"github.com/jasonlvhit/gocron"
)

//Bot defines MyDailyBibleBot structure
type Bot struct {
	config  *configs.Config
	twitter twitter.ITwitter
	bible   bible.BibleAPI
	done    chan bool
}

func NewBot(cfg *configs.Config, twitter twitter.ITwitter, b bible.BibleAPI) *Bot {
	return &Bot{
		twitter: twitter,
		bible:   b,
		config:  cfg,
		done:    make(chan bool, 1),
	}
}

func (b *Bot) Start() error {
	err := LoadVerses()
	if err != nil {
		return err
	}
	c, err := b.twitter.ListenToMentions(b.config.UserName)
	if err != nil {
		return err
	}
	go b.handleMessages(c)

	err = gocron.Every(5).Hour().From(gocron.NextTick()).Do(b.randomPost)
	if err != nil {
		return err
	}
	gocron.Start()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	go func() {
		<-ch
		b.done <- true
	}()
	<-b.done
	b.shutdown()
	return nil
}

func (b *Bot) handleMessages(messages <-chan interface{}) {
	for message := range messages {
		if msg, ok := message.(*twitter.Tweet); ok {
			b.messageHandler(msg)
		}
	}
}

func (b *Bot) messageHandler(tweet *twitter.Tweet) {
	if tweet.User.ScreenName == b.config.UserName {
		return
	}

	verseRequest := NewVerseRequest()
	err := verseRequest.Parse(tweet.Text)
	if err != nil || !verseRequest.IsValid() {
		return
	}

	go b.reply(tweet, verseRequest)
}

func (b *Bot) reply(tweet *twitter.Tweet, req *VerseRequest) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("reply recovered -- %v", err)
		}
	}()

	text, err := b.bible.GetVerse(req.GetPath())
	if err != nil {
		log.Println(err)
		return
	}

	if req.HasImage() {
		processor := NewImageProcessor(&http.Client{Timeout: time.Minute})
		by, err := processor.Process(req.Img, text)
		if err != nil {
			message := fmt.Sprintf("@%s %s", tweet.User.ScreenName, "Sorry we weren't able to process that image.")
			b.twitter.Tweet(message, tweet.ID, nil)
			return
		}
		message := fmt.Sprintf("@%s %s", tweet.User.ScreenName, "")
		b.twitter.Tweet(message, tweet.ID, [][]byte{by})
	} else {
		message := fmt.Sprintf("@%s %s", tweet.User.ScreenName, text)
		if len([]rune(message)) > 280 {
			message = fmt.Sprintf("@%s %s", tweet.User.ScreenName, "Sorry the requested verse is too long for a single tweet.")
		}
		b.twitter.Tweet(message, tweet.ID, nil)
	}

}

func (b *Bot) sendText(tweet *twitter.Tweet, verse string) {
	reponseParts := splitter.Split(verse, 260)
	t := tweet
	var err error
	for index, part := range reponseParts {
		var message string
		//If intial tweet @ the sender, otherwise reply to previous bot message to create a thread
		if index == 0 {
			message = fmt.Sprintf("@%s %s", t.User.ScreenName, part)
		} else {
			message = fmt.Sprintf("@%s %s", b.config.UserName, part)
		}
		err = b.twitter.Tweet(message, t.ID, nil)
		if err != nil {
			return
		}
	}
}

func (b *Bot) randomPost() {
	randomVerse := GetRandomVerse()
	resp, err := b.bible.GetVerse(randomVerse)
	if err != nil {
		return
	}
	reply := fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(resp, "\n", " "), randomVerse)
	b.twitter.Tweet(reply, 0, nil)
}

func (b *Bot) shutdown() {
	b.twitter.Stop()
}
