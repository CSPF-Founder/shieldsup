package controllers

import (
	"net/http"

	"github.com/CSPF-Founder/shieldsup/scanner-api/helpers"
	"github.com/CSPF-Founder/shieldsup/scanner-api/internal/cleanup"

	"github.com/go-chi/chi/v5"
)

type CleanUpResponse struct {
	Success bool `json:"success"`
}

type cleanUpController struct {
	*App
}

func newCleanUpController(app *App) *cleanUpController {
	return &cleanUpController{
		App: app,
	}
}

func (c *cleanUpController) registerRoutes() http.Handler {
	router := chi.NewRouter()
	router.Post("/", c.Cleanup)
	return router
}

func (c *cleanUpController) Cleanup(w http.ResponseWriter, r *http.Request) {
	_, err := cleanup.CleanupAll(c.logger, c.config)

	if err != nil {
		helpers.SendError(w, http.StatusBadRequest, "Invalid Data")
		return
	}
	res := CleanUpResponse{
		Success: true,
	}
	helpers.SendJSON(w, http.StatusOK, res)
}
