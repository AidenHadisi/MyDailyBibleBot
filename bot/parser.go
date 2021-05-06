package bot

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var regex = regexp.MustCompile(`(?P<book>(\d\s?)?\w+)\s(?P<chapter>\d+)\s*:\s*(?P<begin>\d+)(\s?-\s?(?P<end>\d+))?`)

//ParsedText defines parsed text results
type ParsedText struct {
	Book    string
	Chapter string
	Start   string
	End     string
}

//IsValid validates the parsed text
func (p *ParsedText) IsValid() bool {
	if p.Book == "" || p.Chapter == "" || p.Start == "" {
		return false
	}

	start, err := strconv.Atoi(p.Start)
	if err != nil {
		return false
	}

	if p.End != "" {
		end, err := strconv.Atoi(p.End)
		if err != nil {
			return false
		}

		if start >= end || end-start > 10 {
			return false
		}
	}

	return true
}

//GetPath constructs a API path from the ParsedText
func (p *ParsedText) GetPath() string {
	verse := fmt.Sprintf("%s %s:%s", p.Book, p.Chapter, p.Start)

	if p.End != "" {
		verse = fmt.Sprintf("%s-%s", verse, p.End)
	}

	return verse
}

//ParseText parses a tweet text and returns requested verse's info
func ParseText(text string) (*ParsedText, error) {

	result := regex.FindStringSubmatch(text)

	if result == nil {
		return nil, errors.New("incorrect text provided")
	}
	parsed := &ParsedText{
		Book:    strings.Title(strings.ToLower(result[1])),
		Chapter: result[3],
		Start:   result[4],
		End:     result[6],
	}

	return parsed, nil

}
