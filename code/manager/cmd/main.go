package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CSPF-Founder/shieldsup/common/manager/config"
	"github.com/CSPF-Founder/shieldsup/common/manager/db"
	"github.com/CSPF-Founder/shieldsup/common/manager/internal/repositories"
	"github.com/CSPF-Founder/shieldsup/common/manager/internal/scan"
	"github.com/CSPF-Founder/shieldsup/common/manager/logger"
	"github.com/CSPF-Founder/shieldsup/common/manager/utils"
)

type application struct {
	Config config.AppConfig
	DB     *repositories.Repository
	logger *logger.Logger
}

func main() {

	conf := config.LoadConfig()
	appLogger, err := logger.NewLogger(conf.Logging)
	if err != nil {
		logger.GetFallBackLogger().Error("Error setting up logging: ", err)
	}
	appLogger.Info("Initializing ShieldsUp Scanner Manager...")

	ctx, cancel := context.WithCancel(context.Background())
	// Set up a signal channel to capture interrupt and termination signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a goroutine
	go func() {
		// Wait for the interrupt signal
		<-interrupt

		// Perform cleanup operations before exiting (if needed)
		appLogger.Info("Service is stopping...")

		// Cancel the context to signal a graceful shutdown
		cancel()
	}()

	conn, err := db.ConnectDBWithRetry(ctx, conf.DatabaseURI, conf.DatabaseName, 3, 5*time.Second)
	if err != nil {
		appLogger.Fatal("Unable to connect to database", err)
		return
	}
	appLogger.Info("Connected to database")
	defer func() {
		err = conn.Client().Disconnect(ctx)
		if err != nil {
			appLogger.Error("Error disconnecting from database", err)
		}
	}()

	// Wrapper for the SQLC generated models
	app := &application{
		Config: conf,
		DB:     repositories.NewRepository(conn, conf.DeploymentType),
		logger: appLogger,
	}

	scannerInstance := scan.NewScanner(
		conf.ScannerCmd,
		conf.DeploymentType,
		conf.SSHKeyPath,
		app.DB,
		app.logger,
		app.Config.ScanLogsDir,
	)

	appLogger.Info("Service is Running...")

	for {
		select {
		case <-ctx.Done():
			appLogger.Info("Shutting down gracefully...")
			return
		default:
			currentTime := time.Now()
			if currentTime.Hour() == 0 {
				// Do not run scans between 12am and 1am
				if err := utils.SleepContext(ctx, 5*time.Minute); err != nil {
					appLogger.Error("Error sleeping", err)
				}
			} else {
				scannerInstance.Run(ctx)
			}

			if err := utils.SleepContext(ctx, 10*time.Second); err != nil {
				appLogger.Error("Error sleeping", err)
			}
		}
	}
}
