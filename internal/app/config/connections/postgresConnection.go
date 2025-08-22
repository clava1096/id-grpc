package connections

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"id-backend-grpc/internal/app/config"
	"time"
)

func New() (*gorm.DB, error) {
	conf, err := config.LoadPostgresConfig()
	if err != nil {
		return nil, err
	}
	dsn := conf.ConnectionStringDsn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get generic database: %w", err)
	}

	return sqlDB.Close()
}
