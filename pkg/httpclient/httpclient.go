package httpclient

import (
	"net/http"
)

type HttpClient interface {
	//Do executes a request and returns response.
	Do(req *http.Request) (*http.Response, error)
}
