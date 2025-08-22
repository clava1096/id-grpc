package config

type MinioConfig *struct {
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket,omitempty"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

func LoadMinIOConfig() (*MinioConfig, error) {
	var minIOConfig MinioConfig
	// Секреты из Vault
	if err := getSettingFromVault("minio", &minIOConfig); err != nil {
		return nil, err
	}

	return &minIOConfig, nil
}
