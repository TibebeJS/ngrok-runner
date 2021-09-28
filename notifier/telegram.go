package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Telegram struct {
	BotToken string
	ChatId   int
}

func (t *Telegram) Notify(message string) error {
	URL := fmt.Sprintf("%s%s%s", "https://api.telegram.org/bot", t.BotToken, "/sendMessage")

	base, err := url.Parse(URL)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("parse_mode", "Markdown")
	params.Add("chat_id", strconv.Itoa(t.ChatId))
	params.Add("text", message)
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
