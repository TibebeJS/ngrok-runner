package notifier

import (
	"net/http"
	"net/url"
)

type WebHook struct {
	Endpoint string
}

func (b *WebHook) Notify(message string) error {
	URL := b.Endpoint

	base, err := url.Parse(URL)
	if err != nil {
		return err
	}

	params := url.Values{}
	// params.Add("text", message)
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
