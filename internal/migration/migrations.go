package migration

import (
	"id-backend-grpc/internal/app/config/connections"
	"id-backend-grpc/internal/models"
)

func Migrate() error {
	db, err := connections.New()
	if err != nil {

	}
	defer connections.Close(db)
	if err = db.AutoMigrate(&models.User{},
		&models.Token{},
		models.Origin{}); err != nil {
		return err
	}
	return nil
}
