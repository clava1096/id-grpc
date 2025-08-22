package config

type TelegramConfig struct {
	Token    string `yaml:"token"`
	ThreadID int    `yaml:"thread_id"`
	ChatID   string `yaml:"chat_id"`
}

func LoadTelegramConfig() (*TelegramConfig, error) {
	var telegramConfig TelegramConfig
	if err := getSettingFromVault("telegram", &telegramConfig); err != nil {
		return nil, err
	}

	return &telegramConfig, nil
}
