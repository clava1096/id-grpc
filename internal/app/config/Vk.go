package config

type VKConfig struct {
	AppID       string `yaml:"app_id"`
	RedirectURI string `yaml:"redirect_uri"`
}

func LoadVKConfig() (*VKConfig, error) {
	var vkConfig VKConfig
	if err := getSettingFromVault("vk", &vkConfig); err != nil {
		return nil, err
	}

	return &vkConfig, nil
}
