package notifier

import (
	"fmt"
	"net/http"
	"net/url"
)

type WebHook struct {
	Endpoint string
}

func (b *WebHook) Notify(message string) (bool, error) {
	URL := b.Endpoint

	base, err := url.Parse(URL)
	if err != nil {
		return false, err
	}

	params := url.Values{}
	// params.Add("text", message)
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	return resp.StatusCode >= 200 && resp.StatusCode <= 299, err
}
