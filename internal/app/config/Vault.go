package config

import (
	"context"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
	"github.com/mitchellh/mapstructure"
	"os"
)

var vaultClient *vault.Client

func getSettingFromVault[T any](path string, result *T) error {
	config := vault.DefaultConfig() // modify for more granular configuration

	vaultClient, err := vault.NewClient(config)
	if err != nil {
		return err
	}

	host, err := os.ReadFile("/etc/secrets/host")
	if err != nil {
		return err
	}

	err = vaultClient.SetAddress(string(host))
	if err != nil {
		return err
	}

	roleID, err := os.ReadFile("/etc/secrets/role-id")
	if err != nil {
		return err
	}

	secretID := &auth.SecretID{FromFile: "/etc/secrets/secret-id"}

	appRoleAuth, err := auth.NewAppRoleAuth(
		string(roleID),
		secretID,
	)
	if err != nil {
		return err
	}

	if _, err = vaultClient.Auth().Login(context.Background(), appRoleAuth); err != nil {
		return err
	}

	mountPoint, err := os.ReadFile("/etc/secrets/mount-point")
	if err != nil {
		return err
	}

	kvSecret, err := vaultClient.KVv2(string(mountPoint)).Get(context.Background(), path)
	if err != nil {
		return err
	}

	// Декодируем данные из kvSecret.Data в result с использованием mapstructure
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  result,
		TagName: "json", // или "mapstructure" — зависит от тегов структуры
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(kvSecret.Data); err != nil {
		return err
	}

	return nil
}
