package clockify

import (
	"io"
	"net/http"
	"net/url"
	"path"
)

type API struct {
	APIKey   string
	Endpoint string

	url *url.URL
}

func (a *API) call(method string, requestPath string, body io.Reader) (*http.Response, error) {
	if a.url == nil {
		var err error
		a.url, err = url.Parse(a.Endpoint)
		if err != nil {
			return nil, err
		}
	}

	// Copy the cached URL to modify it.
	reqURL := *a.url
	reqURL.Path = path.Join(reqURL.Path, requestPath)
	req, err := http.NewRequest(method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", a.APIKey)

	client := &http.Client{}
	return client.Do(req)
}
