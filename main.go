package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

	viper.AddConfigPath("/etc/ngrok-runner/")

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

func FetchTunnelsAndNotify(configuration c.Configurations) func() {
	var telegramNotifier notifier.Notifier = &notifier.Telegram{
		BotToken: configuration.Telegram.BotToken,
		ChatId:   configuration.Telegram.ChatId,
	}

	var webHookNotifier notifier.Notifier = &notifier.WebHook{
		Endpoint: configuration.WebHook.Endpoint,
		Auth: notifier.Auth{
			Endpoint:  configuration.WebHook.Auth.Endpoint,
			Email:     configuration.WebHook.Auth.Email,
			Password:  configuration.WebHook.Auth.Password,
			FieldName: configuration.WebHook.FieldName,
		},
	}

	return func() {

		tunnels, err := FetchTunnels()

		if err != nil {
			log.Fatalln(err)
		}

		var tunnelsDiscovered struct {
			http  string
			https string
		}

		message := "Tunnels:\n========\n"

		for _, tunnel := range tunnels.Tunnels {
			if tunnel.Protocol == "https" {
				tunnelsDiscovered.https = tunnel.PublicURL
			} else if tunnel.Protocol == "http" {
				tunnelsDiscovered.http = tunnel.PublicURL
			}
			message += fmt.Sprintln(tunnel.Protocol, "\t-", tunnel.PublicURL)
		}

		if err := telegramNotifier.Notify(message); err != nil {
			fmt.Println(err)
		}

		if err := webHookNotifier.Notify(tunnelsDiscovered.https); err != nil {
			fmt.Println(err)

			message := "Webhook Failure:\n========\n"

			message += err.Error()

			telegramNotifier.Notify(message)
		}
	}

}

func main() {

	configuration, err := LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	script := "start_ngrok.sh"
	cmd := exec.Command("nohup", "sh", script)

	cmd.Start()

	c := cron.New()

	go FetchTunnelsAndNotify(configuration)()

	c.AddFunc(configuration.General.Cron, FetchTunnelsAndNotify(configuration))

	go c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
