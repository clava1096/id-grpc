package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"id-backend-grpc/internal/app/config"
	"id-backend-grpc/internal/controllers"
	"id-backend-grpc/internal/initialization"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c, err := config.GetConfig()
	if err != nil {
		log.Fatalf("error while init config: %s", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	setup(c, ctx)
}

func setup(cfg *config.Config, ctx context.Context) {
	logger, err := initialization.CreateLogger(config.GetLogginConfig())
	if err != nil {
		logger.Panic("failed to initialize application: %v", zap.Error(err))
	}

	initializer, err := initialization.NewInitializer(cfg)
	if err != nil {
		logger.Panic("failed to initialize application", zap.Error(err))
	}

	if err := config.Init(ctx); err != nil {
		logger.Panic("Failed to init config", zap.Error(err))
	}
	if err := initializer.StartDatabase(); err != nil {
		logger.Panic("failed to migrate database", zap.Error(err))
	}
	config.DBConnected.Store(true)
	config.IsReady.Store(true)
	e := echo.New()
	route(e, initializer.Is)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		<-ctx.Done()
		logger.Info("Shutdown Server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown server gracefully", zap.Error(err))
		}
		if err := initializer.Close(); err != nil { //TODO дописать
			logger.Error("Failed to close resources", zap.Error(err))
		}

	}()

	logger.Info("Starting HTTP server", zap.String("port", "8080"))
	if err := e.StartServer(server); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}

	logger.Fatal("Exit from server")

}

func route(e *echo.Echo, is *controllers.IdentityService) {
	api := e.Group("/api")

	api.GET("/testHandler", controllers.TestHandler)

	authGroup := api.Group("/user")
	{
		authGroup.POST("/registration", is.Registration)
		authGroup.POST("/auth", is.Authorization)
		authGroup.POST("/logout", is.Logout)
		authGroup.GET("/get-me", is.User)
		authGroup.POST("/verify", is.Verify)
	}

	integrationGroup := api.Group("/integrations")
	{
		integrationGroup.POST("/telegram/auth-callback", is.TelegramAuthCallback)
		integrationGroup.POST("/telegram", is.TelegramIntegration)

		integrationGroup.POST("/vk/auth-callback", is.VKAuthCallback)
		integrationGroup.POST("/vk", is.VKIntegration)
	}

	profileGroup := api.Group("/profile")
	{
		profileGroup.POST("/upload-avatar", is.UploadAvatar)
		profileGroup.GET("/avatar", is.LoadAvatar)
	}

	emailGroup := api.Group("/email")
	{
		emailGroup.POST("/confirm", is.ConfirmEmail)
	}

	passwordGroup := api.Group("/password")
	{
		passwordGroup.POST("/reset", is.ResetPassword)
		passwordGroup.POST("/save-reset", is.SaveResetPassword)
		passwordGroup.POST("/change", is.ChangePassword)
	}

	adminGroup := api.Group("/admin")
	{
		adminGroup.GET("/create-origin-uuid", is.CreateOriginUUID)
	}

	api.POST("/validate-origin-uuid", is.ValidateOriginUUID)
}
