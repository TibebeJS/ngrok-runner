package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	c "github.com/TibebeJS/ngrok-runner/config"
	"github.com/TibebeJS/ngrok-runner/models"
	"github.com/TibebeJS/ngrok-runner/notifier"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func FetchTunnels() (models.TunnelsResponse, error) {
	URL := "http://127.0.0.1:4040/api/tunnels"

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var cResp models.TunnelsResponse

	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Fatal(err)
	}
	return cResp, nil
}

func LoadConfig() (c.Configurations, error) {
	viper.SetConfigName("config")

	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	var configuration c.Configurations

	if err := viper.ReadInConfig(); err != nil {
		return configuration, err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, err
}

func main() {

	configuration, err := LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(configuration)

	var telegramNotifier notifier.Notifier = &notifier.Telegram{
		BotToken: configuration.Telegram.BotToken,
		ChatId:   configuration.Telegram.ChatId,
	}

	var webHookNotifier notifier.Notifier = &notifier.WebHook{
		Endpoint: configuration.WebHook.Endpoint,
	}

	c := cron.New(cron.WithSeconds())

	c.AddFunc(configuration.General.Cron, func() {
		tunnels, err := FetchTunnels()

		if err != nil {
			log.Fatalln(err)
		}

		message := "Tunnels:\n========\n{tun[0]}"

		for _, tunnel := range tunnels.Tunnels {
			message += fmt.Sprintln(tunnel.Protocol, "\t-", tunnel.PublicURL)
		}

		telegramNotifier.Notify(message)
		webHookNotifier.Notify(message)
	})

	go c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	// c.Stop()
}
