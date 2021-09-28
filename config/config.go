package config

type Configurations struct {
	General  GeneralConfigurations
	Telegram TelegramConfigurations
}

type GeneralConfigurations struct {
	Cron string
}

type TelegramConfigurations struct {
	ChatId   int
	BotToken string
}
