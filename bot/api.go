package bot

import (
	"fmt"
	"net/http"
)

const (
	baseURL = "https://bible-api.com/"
)

type ResponseData struct {
	StatusCode int
	Header     http.Header
}

type Response struct {
	ResponseData
	Success interface{}
	Failure interface{}
}

type verse struct {
	BookID   string `json:"book_id"`
	BookName string `json:"book_name"`
	Chapter  int
	Verse    int
	Text     string
}

type BibleOptions struct {
	Translation string `url:"translation"`
}

type Verse struct {
	Reference string
	Verses    []*verse
	Text      string
}

//Response defines the structure for error returned from the API
type ErrResponse struct {
	ErrorMessage string `json:"error"`
}

func (e ErrResponse) Error() string {

	return fmt.Sprintf("API request failed: %s", e.ErrorMessage)

}

//GetVerse gets the requested verses from the API
func (bot *Bot) GetVerse(verse string, bibleOptions *BibleOptions) (*Verse, error) {

	resp, err := bot.get(baseURL+verse, bibleOptions, &Verse{}, &ErrResponse{})

	if err != nil {
		return nil, err
	}

	if resp.Failure != nil {
		return nil, resp.Failure.(*ErrResponse)
	}

	return resp.Success.(*Verse), nil

}
