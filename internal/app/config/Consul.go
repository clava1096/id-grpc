package config

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	IsReady           atomic.Value
	IsHealthy         atomic.Value
	DBConnected       atomic.Value
	VaultConnected    atomic.Value
	RabbitMQConnected atomic.Value
	MinioConnected    atomic.Value
)

const (
	Version = 20250718.03
	AppName = "id-api"
)

func Init(ctx context.Context) error {
	IsReady.Store(false)
	IsHealthy.Store(true)
	DBConnected.Store(false)
	VaultConnected.Store(false)
	RabbitMQConnected.Store(false)
	MinioConnected.Store(false)
	err := os.Setenv(api.HTTPTokenFileEnvName, "/etc/consul/token")
	if err != nil {
		return err
	}

	err = os.Setenv(api.HTTPSSLEnvName, strconv.FormatBool(true))
	if err != nil {
		return err
	}

	host, err := os.ReadFile("/etc/consul/host")
	if err != nil {
		return err
	}

	port, err := os.ReadFile("/etc/consul/port")
	if err != nil {
		return err
	}

	err = os.Setenv(api.HTTPAddrEnvName, fmt.Sprintf("%s:%s", string(host), string(port)))
	if err != nil {
		return err
	}

	c, err := GetConfig()
	if err != nil {
		return err
	}

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	_, _, err = client.Catalog().Service(AppName, c.Env, &api.QueryOptions{})
	if err != nil {
		return err
	}

	go registerConsul(ctx)

	return nil
}

func registerConsul(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	var client *api.Client
	var err error

	// Создаем клиент один раз, если конфиг не меняется
	client, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		zap.Error(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			// Получаем конфиг и регистрируем сервис
			c, err := GetConfig()
			if err != nil {
				zap.Error(err)
				continue
			}

			registration := &api.AgentServiceRegistration{
				Kind: api.ServiceKindTypical,
				Name: AppName,
				Meta: map[string]string{
					"version": strconv.FormatFloat(Version, 'f', 2, 64),
				},
				Port:              443,
				Address:           c.Host,
				Tags:              []string{c.Env},
				EnableTagOverride: true,
				Check: &api.AgentServiceCheck{
					HTTP:     fmt.Sprintf("https://%s/health", c.Host),
					Interval: "10s",
					Timeout:  "5s",
				},
			}

			serviceID := c.Host
			registration.ID = serviceID

			if err := client.Agent().ServiceRegisterOpts(registration, api.ServiceRegisterOpts{ReplaceExistingChecks: true}); err != nil {
				zap.Error(err)
			}
		}
	}
}
