package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

type verse struct {
	BookID   string `json:"book_id"`
	BookName string `json:"book_name"`
	Chapter  int
	Verse    int
	Text     string
}

type stringQuery struct {
	Translation string `url:"translation"`
}

//Response defines the structure for verses returned from the API
type Response struct {
	StatusCode int
	Header     http.Header
	Error      string `json:"error"`
	Reference  string
	Verses     []*verse
	Text       string
}

//GetVerse gets the requested verses from the API
func (bot *Bot) GetVerse(verse string) (*Response, error) {
	return bot.get(verse, &stringQuery{Translation: "kjv"})
}

func (bot *Bot) get(path string, reqData interface{}) (*Response, error) {
	return bot.sendRequest(http.MethodGet, path, reqData)
}

func (bot *Bot) sendRequest(method, path string, reqData interface{}) (*Response, error) {
	resp := &Response{}

	req, err := http.NewRequest(method, baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	v, err := query.Values(reqData)
	if err == nil {
		req.URL.RawQuery = v.Encode()

	}

	err = bot.doRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (bot *Bot) doRequest(req *http.Request, res *Response) error {

	response, err := bot.HTTPClient.Do(req)

	if err != nil {
		return fmt.Errorf("API request error: %s", err.Error())
	}
	defer response.Body.Close()

	res.Header = response.Header
	res.StatusCode = response.StatusCode

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &res)

	if err != nil {
		return fmt.Errorf("failed to decode API response: %s", err.Error())
	}

	return nil
}
