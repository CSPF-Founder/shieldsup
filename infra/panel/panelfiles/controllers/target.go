package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	ctx "github.com/CSPF-Founder/shieldsup/onpremise/panel/context"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/repositories/datarepos"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/services"
	mid "github.com/CSPF-Founder/shieldsup/onpremise/panel/middlewares"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models/datamodels"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/utils"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/views"

	"github.com/go-chi/chi/v5"
)

type targetController struct {
	*App
	targetRepo     datarepos.TargetRepository
	scanResultRepo datarepos.ScanResultRepository
}

func newTargetController(app *App, targetRepo datarepos.TargetRepository, scanResultRepo datarepos.ScanResultRepository) *targetController {
	return &targetController{
		App:            app,
		targetRepo:     targetRepo,
		scanResultRepo: scanResultRepo,
	}
}

func (c *targetController) registerRoutes() http.Handler {
	router := chi.NewRouter()

	// Authenticated Routes
	router.Group(func(r chi.Router) {
		r.Use(mid.RequireLogin)

		r.Get("/list", c.List)                            // List Targets table
		r.Post("/list", c.List)                           // Get All Targets List
		r.Get("/add", c.DisplayAdd)                       // Display add Target form
		r.Post("/add", c.AddHandler)                      // Handle add target form
		r.Post("/check-multi-status", c.CheckMultiStatus) // Check multiple target status

		// eg: DELETE /scans/1
		r.Route("/{targetID:[0-9a-fA-F]{24}}", func(r chi.Router) {
			r.Delete("/", c.Delete) // Handle Delete action
			r.Get("/report", c.DownloadReport)

			scanResult := newScanResultController(c.App,
				c.scanResultRepo,
				c.targetRepo,
			)
			r.Mount("/scan-results", scanResult.registerRoutes())
		})

	})

	return router
}

// Show Add Target Page
func (c *targetController) DisplayAdd(w http.ResponseWriter, r *http.Request) {
	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "Add Scan"
	if err := views.RenderTemplate(w, "target/add", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}

// Add Target Handler
func (c *targetController) AddHandler(w http.ResponseWriter, r *http.Request) {
	targetAddress := r.PostFormValue("target_address")

	targetType := enums.ParseTargetType(targetAddress)
	if targetType == "" || targetType == enums.TargetTypeInvalid {
		c.SendJSONError(w, "Invalid Target")
		return
	}

	user := ctx.Get(r, "user").(models.User)
	target := &datamodels.Target{
		TargetAddress: targetAddress,
		Flag:          enums.ScanFlagWaitingToStart,
		ScanStatus:    enums.TargetStatusYetToStart,
		TargetType:    targetType,
		CustomerName:  user.Username,
	}

	_, err := c.targetRepo.SaveTarget(r.Context(), target)
	if err != nil {
		c.SendJSONError(w, "Unable to add the scan")
		return
	}
	c.SendJSONSuccess(w, "Successfully added the scan")
}

// Show Target List
func (c *targetController) List(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		data := map[string]int{
			"YET_TO_START":     int(enums.TargetStatusYetToStart),
			"SCAN_STARTED":     int(enums.TargetStatusScanStarted),
			"REPORT_GENERATED": int(enums.TargetStatusReportGenerated),
			"SCAN_FAILED":      int(enums.TargetStatusScanFailed),
		}
		templateData := views.NewTemplateData(c.config, c.session, r)
		templateData.Title = "View Scans"
		templateData.Data = data

		if err := views.RenderTemplate(w, "target/list", templateData); err != nil {
			c.logger.Error("Error rendering template: ", err)
		}
	}
	if r.Method == http.MethodPost {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "Unable to fetch the data", http.StatusUnprocessableEntity)
			}
		}()

		// Parse request parameters
		draw, err := strconv.Atoi(r.FormValue("draw"))
		if err != nil {
			http.Error(w, "Invalid draw parameter", http.StatusBadRequest)
			return
		}

		start, err := strconv.ParseInt(r.FormValue("start"), 10, 64)
		if err != nil || start < 0 {
			start = 0
		}

		length, err := strconv.ParseInt(r.FormValue("length"), 10, 64)
		if err != nil || length < 0 {
			length = 0
		}

		username := ctx.Get(r, "user").(models.User).Username // Fetch username from Auth::user()->getUsername()
		totalData, err := c.targetRepo.CountByCustomerUsername(r.Context(), username)
		if err != nil {
			c.logger.Error("Error fetching the data: ", err)
			http.Error(w, "Unable to fetch the data", http.StatusUnprocessableEntity)
			return
		}

		totalFiltered := totalData
		targetList, err := c.targetRepo.ListByCustomerUsername(r.Context(), username, start, length)
		if err != nil {
			c.logger.Error("Error fetching the data: ", err)
			http.Error(w, "Unable to fetch the data", http.StatusUnprocessableEntity)
			return
		}

		data := make([]map[string]any, 0)

		for _, target := range targetList {

			nestedData := make(map[string]any)
			nestedData["id"] = fmt.Sprintf("%v", target.ID.Hex())
			nestedData["target_address"] = target.TargetAddress

			if target.ScanStatus == enums.TargetStatusScanStarted {
				nestedData["scan_status_text"] = `<span class="spinner-border spinner-border-sm text-primary" aria-hidden="true"></span> <span role="status">Scanning...</span>`
			} else {
				nestedData["scan_status_text"], err = enums.TargetStatusMap.GetText(target.ScanStatus)
				if err != nil {
					nestedData["scan_status_text"] = "Unknown"
				}
			}

			nestedData["scan_status"] = target.ScanStatus
			scanStartedText := ""
			scanCompletedText := ""
			if target.ScanStartedTime == nil {
				scanStartedText = "-"
			} else {
				scanStartedText = target.ScanStartedTimeStr()
			}
			if target.ScanCompletedTime == nil {
				scanCompletedText = "-"
			} else {
				scanCompletedText = target.ScanCompletedTimeStr()
			}
			nestedData["scan_started_time"] = scanStartedText
			nestedData["scan_completed_time"] = scanCompletedText

			action := ""
			if target.ScanStatus == enums.TargetStatusReportGenerated {
				action += fmt.Sprintf(`<a class="btn btn-sm btn-primary m-1 report-button" href="/targets/%v/report">Report</a>`, target.ID.Hex())
				action += fmt.Sprintf(`<a class="btn btn-sm btn-primary m-1 alerts-button" href="/targets/%v/scan-results">Alerts</a>`, target.ID.Hex())
			} else {
				action += fmt.Sprintf(`<a class="btn btn-sm btn-dark m-1 disabled report-button" disabled href="/targets/%v/report">Report</a>`, target.ID.Hex())
				action += fmt.Sprintf(`<a class="btn btn-sm btn-dark m-1 disabled alerts-button" disabled href="/targets/%v/scan-results">Alerts</a>`, target.ID.Hex())
			}

			if target.ScanStatus == enums.TargetStatusScanStarted {
				action += fmt.Sprintf(`<button data-id="%v" class="btn btn-sm btn-dark text-white m-1 delete-target disabled" disabled>Delete</button>`, target.ID.Hex())
			} else {
				action += fmt.Sprintf(`<button data-id="%v" class="btn btn-sm btn-danger text-white m-1 delete-target">Delete</button>`, target.ID.Hex())
			}

			nestedData["action"] = action

			data = append(data, nestedData)
		}

		jsonOutput := map[string]any{
			"draw":            draw,
			"recordsTotal":    totalData,
			"recordsFiltered": totalFiltered,
			"records":         data,
		}

		w.Header().Set("Content-Type", "application/json")
		c.App.JSONResponse(w, jsonOutput, 200)
	}
}

func (c *targetController) Delete(w http.ResponseWriter, r *http.Request) {
	targetID := chi.URLParam(r, "targetID")

	user := ctx.Get(r, "user").(models.User)

	// Fetching the target from the database
	target, err := c.targetRepo.ByIdAndCustomerUsername(r.Context(), targetID, user.Username)
	if err != nil {
		c.SendJSONError(w, "Invalid request")
		return
	}

	if services.DeleteTarget(r.Context(), *target, c.targetRepo) {
		c.SendJSONSuccess(w, "Successfully deleted the scan", http.StatusOK)
	} else {
		c.SendJSONError(w, "Unable to delete the scan", http.StatusUnprocessableEntity)
	}
}

func (c *targetController) DownloadReport(w http.ResponseWriter, r *http.Request) {
	targetID := chi.URLParam(r, "targetID")

	user := ctx.Get(r, "user").(models.User)
	target, err := c.targetRepo.ByIdAndCustomerUsername(r.Context(), targetID, user.Username)
	if err != nil {
		c.FlashAndGoBack(w, r, enums.FlashWarning, "Invalid Request")
		return
	}

	err = services.DownloadReport(*target, w, r)
	if err != nil {
		c.logger.Error("Error downloading report: ", err)
		c.FlashAndGoBack(w, r, enums.FlashWarning, "Unable to download the report")
	}
}

func (c *targetController) CheckMultiStatus(w http.ResponseWriter, r *http.Request) {
	requiredParams := []string{"target_ids[]"}

	if !utils.CheckAllParamsExist(r, requiredParams) {
		c.SendJSONError(w, "Please fill all the inputs", http.StatusUnprocessableEntity)
		return
	}

	targetIds := r.PostForm["target_ids[]"]

	user := ctx.Get(r, "user").(models.User)
	targetList, err := c.targetRepo.ListByIdsAndCustomerUsername(r.Context(), targetIds, user.Username)
	if err != nil {
		c.SendJSONError(w, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	if targetList == nil {
		c.SendJSONError(w, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	var entries []map[string]any

	for _, target := range targetList {
		scanStatusText, err := enums.TargetStatusMap.GetText(target.ScanStatus)
		if err != nil {
			scanStatusText = "Unknown"
		}
		entry := map[string]any{
			"id":                  target.ID.Hex(),
			"scan_status":         target.ScanStatus,
			"scan_status_text":    scanStatusText,
			"scan_started_time":   target.ScanStartedTimeStr(),
			"scan_completed_time": target.ScanCompletedTimeStr(),
		}
		entries = append(entries, entry)
	}

	c.JSONResponse(w, map[string]any{"data": entries}, http.StatusOK)

}
