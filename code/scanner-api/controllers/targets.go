package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/CSPF-Founder/shieldsup/scanner-api/helpers"
	"github.com/CSPF-Founder/shieldsup/scanner-api/internal/scanner"

	"github.com/go-chi/chi/v5"
)

type InputItem struct {
	Target string `json:"target"`
	Force  bool   `json:"force"`
}

type TargetAddResponse struct {
	Success bool `json:"success"`
}

type targetController struct {
	*App
}

func newTargetController(app *App) *targetController {
	return &targetController{
		App: app,
	}
}

func (c *targetController) registerRoutes() http.Handler {
	router := chi.NewRouter()
	router.Post("/", c.AddScan)
	router.Post("/results", c.GetResult)
	return router
}

func (c *targetController) AddScan(w http.ResponseWriter, r *http.Request) {
	sScanStatus := os.Getenv("n_scanstatus")

	if sScanStatus != "NONE" && sScanStatus != "" {
		helpers.SendError(w, http.StatusBadRequest, "Another scan is already running")
		return
	}

	var inputItem InputItem
	decodeErr := json.NewDecoder(r.Body).Decode(&inputItem)
	if decodeErr != nil {
		helpers.SendError(w, http.StatusUnprocessableEntity, "Unprocessable Entity")
		return
	}

	go func() {
		module := scanner.NewScannerModule(c.config, c.logger, inputItem.Target)
		module.StartScan()
	}()

	res := TargetAddResponse{
		Success: true,
	}
	helpers.SendJSON(w, http.StatusOK, res)
}

func (c *targetController) GetResult(w http.ResponseWriter, r *http.Request) {
	var inputItem InputItem

	err := json.NewDecoder(r.Body).Decode(&inputItem)
	if err != nil {
		helpers.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	module := scanner.NewScannerModule(c.config, c.logger, inputItem.Target)
	output := module.RetrieveResults(inputItem.Force)
	// if _, ok := output.(map[string]any); ok {
	helpers.SendJSON(w, http.StatusOK, output)
	// }

	//TODO: Check the logic again
	// res := map[string]any{
	// 	"Success":    false,
	// 	"ScanStatus": enums.ScanStatusNotCompleted,
	// 	"Message":    []string{"Scan not completed"},
	// }

	// helpers.SendJSON(w, http.StatusOK, res)
}
