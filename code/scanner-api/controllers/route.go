package controllers

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner-api/config"
	"github.com/CSPF-Founder/shieldsup/scanner-api/logger"
	mid "github.com/CSPF-Founder/shieldsup/scanner-api/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	server *http.Server
	config *config.Config
	logger *logger.Logger
	// httpClient httpclient.HttpClient
	ScanLock sync.Mutex
}

// ServerOption is a functional option that is used to configure the
type ServerOption func(*App)

// NewApp returns a new instance of the app with
// provided options applied.
func NewApp(config *config.Config, appLogger *logger.Logger, options ...ServerOption) *App {
	defaultServer := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         config.ServerConf.ServerAddress,
	}

	app := &App{
		server: defaultServer,
		config: config,
		logger: appLogger,
	}
	for _, opt := range options {
		opt(app)
	}

	// app.httpClient = &http.Client{}
	return app
}

// Start launches the server, listening on the configured address.
func (app *App) StartServer() {
	// If TLS isn't configured, just listen on HTTP
	app.logger.Info(fmt.Sprintf("Starting server at http://%s", app.config.ServerConf.ServerAddress))
	err := app.server.ListenAndServe()
	if err != nil {
		app.logger.Fatal("Error starting server: ", err)
	}
}

// Shutdown attempts to gracefully shutdown the server.
func (app *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return app.server.Shutdown(ctx)
}

// apiV1Routes defines the routes for the API
func (app *App) apiV1Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", HealthChecker)

	// Authenticated Routes
	r.Group(func(r chi.Router) {
		r.Use(mid.RequireStaticAPIKey(app.config.ServerConf.APIKey))

		cleanupCtrl := newCleanUpController(app)
		targetCtrl := newTargetController(app)

		// Routes for endpoints & Subroutes inside the endpoints
		r.Mount("/cleanup", cleanupCtrl.registerRoutes())
		r.Mount("/targets", targetCtrl.registerRoutes())

	})

	return r
}

// SetupRoutes initializes the routes for the app
func (app *App) SetupRoutes() {

	router := chi.NewRouter() // Initialize Chi router

	// r.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))
	router.Use(mid.LoggingMiddleware(app.logger))

	router.Mount("/api/v1", app.apiV1Routes())

	routeHandler := mid.Use(
		router.ServeHTTP,
		mid.ApplySecurityHeaders,
	)
	app.server.Handler = routeHandler
}
