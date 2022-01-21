package bot

import (
	"fmt"
	"log"

	"github.com/AidenHadisi/MyDailyBibleBot/configs"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/bible"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/cron"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/image"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/parser"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/twitter"

	twt "github.com/dghubble/go-twitter/twitter"
)

//Bot defines MyDailyBibleBot structure
type Bot struct {
	config  *configs.Config
	twitter twitter.ITwitter
	bible   *bible.BibleAPI
	cron    cron.Cron
	image   *image.ImageProcessor
}

func NewBot(cfg *configs.Config, twitter twitter.ITwitter, b *bible.BibleAPI, cron cron.Cron, image *image.ImageProcessor) *Bot {
	return &Bot{
		twitter: twitter,
		bible:   b,
		config:  cfg,
		cron:    cron,
		image:   image,
	}
}

func (b *Bot) Init() error {
	//init the api client
	err := b.bible.Init()
	if err != nil {
		return err
	}

	//start listening to twitter
	c, err := b.twitter.ListenToMentions(b.config.UserName)
	if err != nil {
		return err
	}
	go b.handleMessages(c)

	//start the cron
	err = b.cron.CreateJob("0 */5 * * *", b.randomPost)
	if err != nil {
		return err
	}

	err = b.cron.StartCrons()
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleMessages(messages <-chan interface{}) {
	for message := range messages {
		if msg, ok := message.(*twt.Tweet); ok {
			b.messageHandler(msg)
		}
	}
}

func (b *Bot) messageHandler(tweet *twt.Tweet) {
	if tweet.User.ScreenName == b.config.UserName {
		return
	}

	verseRequest := parser.NewParser()
	err := verseRequest.Parse(tweet.Text)
	if err != nil {
		return
	}

	go b.reply(tweet, verseRequest)
}

func (b *Bot) reply(tweet *twt.Tweet, parsed *parser.Parser) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("reply recovered -- %v", err)
		}
	}()

	text, err := b.bible.GetVerse(parsed.GetPath())
	if err != nil {
		log.Println(err)
		return
	}

	if parsed.HasImage() {
		by, err := b.image.Process(parsed.Img, text, parsed.Size)
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

func (b *Bot) randomPost() {
	resp, err := b.bible.GetRandomVerse()
	if err != nil {
		return
	}
	image, err := b.image.Process("https://picsum.photos/1200/625", resp, 40)
	if err != nil {
		return
	}
	b.twitter.Tweet("", 0, [][]byte{image})

}

func (b *Bot) Shutdown() {
	b.twitter.Stop()
	b.cron.StopCrons()
}
