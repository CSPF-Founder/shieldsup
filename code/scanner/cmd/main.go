package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CSPF-Founder/shieldsup/scanner/config"
	"github.com/CSPF-Founder/shieldsup/scanner/db"
	"github.com/CSPF-Founder/shieldsup/scanner/internal/repositories"
	"github.com/CSPF-Founder/shieldsup/scanner/internal/scan"
	"github.com/CSPF-Founder/shieldsup/scanner/logger"
	"github.com/CSPF-Founder/shieldsup/scanner/utils"
)

type Application struct {
	Config *config.Config
	DB     *repositories.Repository
	Logger *logger.Logger
}

func NewApplication(conf *config.Config, appLogger *logger.Logger) *Application {
	return &Application{
		Config: conf,
		Logger: appLogger,
	}
}

type CLIInput struct {
	Module   string
	TargetID string
}

func main() {
	// Load the config
	conf, err := config.LoadConfig()
	// Just warn if a contact address hasn't been configured
	if err != nil {
		log.Fatal("Error loading config", err)
	}

	appLogger, err := logger.NewLogger(conf.Logging)
	if err != nil {
		log.Fatal("Error setting up logging", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Set up a signal channel to capture interrupt and termination signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a goroutine
	go func() {
		// Wait for the interrupt signal
		<-interrupt

		// Perform cleanup operations before exiting (if needed)
		appLogger.Info("Scanner is stopping...")

		// Cancel the context to signal a graceful shutdown
		cancel()
	}()

	app, err := baseSetup(ctx, conf, appLogger)
	if err != nil {
		appLogger.Fatal("Error setting up app", err)
		return
	}

	inputTargetID, err := handleCLI()
	if err != nil {
		appLogger.Fatal("Error handling CLI", err)
		return
	}

	select {
	case <-ctx.Done():
		// Context has timed out
		app.Logger.Info("Scanner has stopped")
	default:
		err = app.RunScan(ctx, inputTargetID)
		if err != nil {
			appLogger.Fatal("Error running scan", err)
			return
		}

	}
}

func baseSetup(ctx context.Context, conf *config.Config, appLogger *logger.Logger) (*Application, error) {
	app := NewApplication(conf, appLogger)

	dbRepos, err := db.SetupDatabase(ctx, conf)
	if err != nil {
		return nil, err
	}
	app.DB = dbRepos

	// if _, err := os.Stat(conf.WorkDir); os.IsNotExist(err) {
	// 	err = os.MkdirAll(conf.WorkDir, 0755)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// if _, err := os.Stat(conf.TempUploadsDir); os.IsNotExist(err) {
	// 	err = os.MkdirAll(conf.TempUploadsDir, 0755)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return app, nil
}

func handleCLI() (string, error) {
	var targetID string

	flag.StringVar(&targetID, "t", "", "target ID")

	flag.Parse()

	if !utils.IsValidObjectId(targetID) {
		return "", fmt.Errorf("Invalid object id received")
	}
	return targetID, nil
}

func (app *Application) RunScan(ctx context.Context, targetID string) error {
	scanner, err := scan.NewScannerModule(app.Logger, app.Config, app.DB, targetID)
	if err != nil {
		return err
	}
	scanner.Run(ctx)

	return nil
}
