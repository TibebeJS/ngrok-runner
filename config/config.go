package config

type Configurations struct {
	General  GeneralConfigurations
	Telegram TelegramConfigurations
	WebHook  WebHookConfigurations
}

type GeneralConfigurations struct {
	Cron string
}

type TelegramConfigurations struct {
	ChatId   int
	BotToken string
}

type WebHookConfigurations struct {
	Endpoint string
}
