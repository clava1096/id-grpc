package initialization

import (
	"errors"
	"fmt"
	"id-backend-grpc/internal/app/config"
	"id-backend-grpc/internal/controllers"
	"id-backend-grpc/internal/migration"
	"id-backend-grpc/internal/repositories"
	"id-backend-grpc/internal/services"
)

type Initializer struct {
	Is *controllers.IdentityService

	// TODO здесь будет всякий кал по типу конфига, IdentityService и прочего
}

func NewInitializer(cfg *config.Config) (*Initializer, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	// Инициализация репозиториев
	tokenRepo, err := repositories.NewTokenRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create token repository: %w", err)
	}

	userRepo, err := repositories.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	originRepo, err := repositories.NewOriginRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create origin repository: %w", err)
	}

	// Инициализация сервисов
	tokenService := services.NewTokenService(tokenRepo)
	userService := services.NewUserService(userRepo)
	originService := services.NewOriginService(originRepo)
	identityHandler := controllers.NewIdentityService(tokenService, userService, originService)
	return &Initializer{
		Is: identityHandler,
	}, nil
}

func (i *Initializer) StartDatabase() error {
	err := migration.Migrate()
	if err != nil {
		return err
	}
	return nil
}

func (i *Initializer) Close() error { // TODO для остановки
	return nil
}
