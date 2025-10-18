package main

import (
	"github.com/ichigo7diabol/go-test-wallet/api/handlers"
	"github.com/ichigo7diabol/go-test-wallet/api/openapi"
	"github.com/ichigo7diabol/go-test-wallet/internal/app"
	"github.com/ichigo7diabol/go-test-wallet/internal/config"
	"github.com/ichigo7diabol/go-test-wallet/internal/models"
	"github.com/labstack/echo/v4"
	"go.infratographer.com/x/echox/echozap"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	z, _ := zap.NewProduction()
	defer z.Sync()

	z.Info("Initialziaing")

	e := echo.New()
	e.Logger = echozap.NewLogger(z)

	z.Info("Loading configs")
	config := config.Load()

	z.Info("Configfiguration",
		zap.String("Port", config.Port),
		zap.String("Dsn", config.Dsn),
	)
	z.Info("Connecting to database")
	db, cleanup, err := NewDatabaseConnection(config.Dsn)
	if err != nil {
		z.Sugar().Fatal(err)
	}
	defer cleanup()
	if err := db.AutoMigrate(&models.WalletModel{}); err != nil {
		z.Sugar().Fatal(err)
	}
	repository := app.NewRepository(db)
	walletService := app.NewWalletService(repository)
	h := handlers.NewWalletHandler(walletService)
	openapi.RegisterHandlers(e, h)

	z.Info("Setup middleware")
	e.Use(echozap.Middleware(z))

	z.Info("Starting server")
	e.Logger.Fatal(e.Start(":" + config.Port))
}

func NewDatabaseConnection(dsn string) (*gorm.DB, func() error, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() error { return sqlDB.Close() }
	return db, cleanup, nil
}
