package config

type RabbitMQConfig *struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	VHost    string `json:"v_host"`
}

func LoadRabbitMQConfig() (*RabbitMQConfig, error) {
	var rabbitMQConfig RabbitMQConfig
	// Секреты из Vault
	if err := getSettingFromVault("rabbitmq", &rabbitMQConfig); err != nil {
		return nil, err
	}

	return &rabbitMQConfig, nil
}
