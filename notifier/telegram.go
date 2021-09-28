package notifier

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Telegram struct {
	BotToken string
	ChatId   int
}

func (b *Telegram) Notify(message string) (bool, error) {
	URL := fmt.Sprintf("%s%s%s", "https://api.telegram.org/bot", b.BotToken, "/sendMessage")

	base, err := url.Parse(URL)
	if err != nil {
		return false, err
	}

	params := url.Values{}
	params.Add("parse_mode", "Markdown")
	params.Add("chat_id", strconv.Itoa(b.ChatId))
	params.Add("text", message)
	base.RawQuery = params.Encode()

	fmt.Printf("Encoded URL is %q\n", base.String())

	resp, err := http.Get(base.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	return resp.StatusCode >= 200 && resp.StatusCode <= 299, err
}
