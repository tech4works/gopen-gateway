package interfaces

import "net/http"

type RestTemplate interface {
	MakeRequest(httpRequest *http.Request) (*http.Response, error)
}
