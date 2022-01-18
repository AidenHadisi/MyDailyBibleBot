package twitter

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
)

type ITwitter interface {
	//ListenToMentions starts a twitter stream and returns a channel where tweets can be listened on.
	ListenToMentions(username string) (<-chan interface{}, error)

	//Tweet tweets a message. It can include a tweet id for reply and many images.
	Tweet(text string, inReplyTo int64, images [][]byte) error

	//Stop stops the twitter stream and closes channels.
	Stop()
}

//TwitterApi is a facade for 3rd party twitter api client
type TwitterApi struct {
	client     *twitter.Client
	streams    []*twitter.Stream
	streamLock sync.Mutex
}

func NewTwitterApi(client *twitter.Client) *TwitterApi {
	return &TwitterApi{
		client:  client,
		streams: make([]*twitter.Stream, 0),
	}
}

//ListenToMentions starts listening to a twitter mentions of given username.
//It returns a channel that can be listened on.
func (t *TwitterApi) ListenToMentions(username string) (<-chan interface{}, error) {
	t.streamLock.Lock()
	defer t.streamLock.Unlock()
	params := &twitter.StreamFilterParams{
		Track:         []string{fmt.Sprintf("@%s", username)},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := t.client.Streams.Filter(params)
	if err != nil {
		return nil, err
	}

	t.streams = append(t.streams, stream)

	return stream.Messages, nil
}

//Tweet tweets a response at a user
func (t *TwitterApi) Tweet(text string, inReplyTo int64, images [][]byte) error {
	var mediaIds []int64
	var err error
	if len(images) > 0 {
		mediaIds, err = t.uploadImagesToTwitter(images)
		if err != nil {
			return err
		}
	}
	_, _, err = t.client.Statuses.Update(text, &twitter.StatusUpdateParams{
		MediaIds:          mediaIds,
		InReplyToStatusID: inReplyTo,
	})

	return err
}

func (t *TwitterApi) uploadImagesToTwitter(images [][]byte) ([]int64, error) {
	var mediaIds []int64
	for _, img := range images {
		res, _, err := t.client.Media.Upload(img, "png")
		if err != nil || res.MediaID <= 0 {
			return nil, errors.New("failed to upload images to twitter")
		}

		mediaIds = append(mediaIds, res.MediaID)
	}

	return mediaIds, nil
}

func (t *TwitterApi) Stop() {
	t.streamLock.Lock()
	defer t.streamLock.Unlock()
	for _, s := range t.streams {
		s.Stop()
	}
}
