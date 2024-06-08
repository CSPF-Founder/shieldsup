package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/CSPF-Founder/shieldsup/scanner-api/config"
	"github.com/CSPF-Founder/shieldsup/scanner-api/controllers"
	"github.com/CSPF-Founder/shieldsup/scanner-api/logger"
)

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

	app, err := baseSetup(conf, appLogger)
	if err != nil {
		appLogger.Fatal("Error setting up app", err)
		return
	}

	//Start Server
	go app.StartServer()

	// Handle shutdown here
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	appLogger.Info("CTRL+C Received... shutting down servers")
	defer func() {
		_ = app.Shutdown()
	}()
}

func baseSetup(conf *config.Config, appLogger *logger.Logger) (*controllers.App, error) {

	app := controllers.NewApp(conf, appLogger)
	app.SetupRoutes()

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
