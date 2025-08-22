package config

import (
	"fmt"
	"os"
	"strconv"
)

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database,omitempty"`
}

func (c *PostgresConfig) ConnectionStringDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		c.Host,
		c.User,
		c.Password,
		c.Database,
		c.Port,
	)
}

func LoadPostgresConfig() (*PostgresConfig, error) {
	host, err := os.ReadFile("/etc/postgres/host")
	if err != nil {
		return nil, err
	}

	dataBase, err := os.ReadFile("/etc/postgres/database")
	if err != nil {
		return nil, err
	}

	port, err := os.ReadFile("/etc/postgres/port")
	if err != nil {
		return nil, err
	}

	user, err := os.ReadFile("/etc/postgres/username")
	if err != nil {
		return nil, err
	}

	password, err := os.ReadFile("/etc/postgres/password")
	if err != nil {
		return nil, err
	}

	portInt, err := strconv.Atoi(string(port))

	if err != nil {
		return nil, err
	}

	return &PostgresConfig{
		Host:     string(host),
		Port:     uint(portInt),
		Database: string(dataBase),
		User:     string(user),
		Password: string(password),
	}, nil

}
