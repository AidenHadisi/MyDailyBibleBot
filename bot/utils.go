package bot

import (
	"fmt"

	"github.com/AidenHadisi/go-text-splitter/splitter"
	"github.com/dghubble/go-twitter/twitter"
)

func BreakupString(s string, size int) []string {
	if size >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, size)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == size {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}

func (bot *Bot) textToTweet(text string, tweet *twitter.Tweet) {
	reponseParts := splitter.Split(text, 260)

	t := tweet
	var err error
	for index, part := range reponseParts {
		var message string
		//If intial tweet @ the sender, otherwise reply to previous bot message to create a thread
		if index > 0 {
			message = fmt.Sprintf("@%s %s", botUsername, part)
		} else {
			message = fmt.Sprintf("@%s %s", t.User.ScreenName, part)

		}
		t, _, err = bot.TwitterClient.Statuses.Update(message, &twitter.StatusUpdateParams{
			InReplyToStatusID: t.ID,
		})

		if err != nil {
			return
		}

	}

}
