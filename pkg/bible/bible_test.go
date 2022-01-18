package bible

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/AidenHadisi/MyDailyBibleBot/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetVerse(t *testing.T) {

	tests := []struct {
		name      string
		verse     string
		result    string
		status    int
		wantErr   bool
		fromCache bool
		cacheErr  error
	}{
		{"successful api call", "John 1:1", "result John 1:1", http.StatusOK, false, false, errors.New("dne")},
		{"unsuccessful api call", "John 1:1", "result John 1:1", http.StatusBadRequest, true, false, errors.New("dne")},
		{"from cache", "John 1:2", "result John 1:2", http.StatusOK, false, true, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockedHttpClient := new(mocks.HttpClient)
			mockedCache := new(mocks.Cache)
			var cacheData string
			if test.fromCache {
				cacheData = fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(test.result, "\n", " "), test.verse)
			} else {
				cacheData = ""
			}

			mockedCache.On("Get", test.verse).Return(cacheData, test.cacheErr)
			mockedCache.On("Set", test.verse, fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(test.result, "\n", " "), test.verse), mock.Anything).Return(nil)

			r := ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(`{"text": "%s"}`, test.result))))
			mockedHttpClient.On("Do", mock.Anything).Return(&http.Response{Body: r, StatusCode: test.status}, nil)

			api := NewBibleAPI(mockedHttpClient, mockedCache)
			result, err := api.GetVerse(test.verse)
			if test.fromCache {
				mockedHttpClient.AssertNotCalled(t, "Do", mock.Anything)
			} else {
				if !test.wantErr {
					mockedCache.AssertNumberOfCalls(t, "Set", 1)

				}
			}
			if test.wantErr {
				assert.Error(t, err)
			} else {

				assert.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(test.result, "\n", " "), test.verse), result)
			}

		})
	}
}
