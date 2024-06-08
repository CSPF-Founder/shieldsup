package controllers

import (
	"net/http"
	"sort"

	ctx "github.com/CSPF-Founder/shieldsup/onpremise/panel/context"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/repositories/datarepos"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/services"
	mid "github.com/CSPF-Founder/shieldsup/onpremise/panel/middlewares"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models/datamodels"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/utils/iputils"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/views"

	"github.com/go-chi/chi/v5"
)

type ScanResultResponse struct {
	Target                              *datamodels.Target
	Records                             *map[string][]datamodels.ScanResult
	ScanResults                         []datamodels.ScanResult
	VulnerabilityStats                  Vulnerability
	OverallCVSSScore                    float64
	TotalAlerts                         int
	TotalTargets                        int64
	CVSSScoreByHost                     map[string]float64
	DefaultRemediation                  string
	NumberOfTARowsForDefaultRemediation int
}

type Vulnerability struct {
	Critical          int
	High              int
	Medium            int
	Low               int
	Info              int
	NoVulnerabilities int
}

type scanResultController struct {
	*App
	scanResultRepo datarepos.ScanResultRepository
	targetRepo     datarepos.TargetRepository
}

func newScanResultController(app *App,
	scanResultRepo datarepos.ScanResultRepository,
	targetRepo datarepos.TargetRepository,
) *scanResultController {
	return &scanResultController{
		App:            app,
		scanResultRepo: scanResultRepo,
		targetRepo:     targetRepo,
	}
}

func (c *scanResultController) registerRoutes() http.Handler {
	router := chi.NewRouter()

	// Authenticated Routes
	router.Group(func(r chi.Router) {
		r.Use(mid.RequireLogin)

		r.Get("/", c.List) // List all Scan Results
	})

	return router
}

func (c *scanResultController) List(w http.ResponseWriter, r *http.Request) {
	targetID := chi.URLParam(r, "targetID")

	user := ctx.Get(r, "user").(models.User)
	target, err := c.targetRepo.ByIdAndCustomerUsername(r.Context(), targetID, user.Username)
	if err != nil {
		c.logger.Error("Error getting target", err)
		c.FlashAndGoBack(w, r, enums.FlashWarning, "Target Not Found")
		return
	}

	scanResults, err := c.scanResultRepo.ListByTarget(r.Context(), target.ID)
	if err != nil {
		c.logger.Error("Error getting alerts", err)
		c.FlashAndGoBack(w, r, enums.FlashDanger, "Issue getting alerts")
		return
	}

	if len(scanResults) == 0 {
		c.FlashAndGoBack(w, r, enums.FlashInfo, "No alerts found")
		return
	}

	vulnerabilityStats := Vulnerability{
		Critical:          0,
		High:              0,
		Medium:            0,
		Low:               0,
		Info:              0,
		NoVulnerabilities: 1,
	}
	totalAlerts := 0

	for _, result := range scanResults {
		totalAlerts++
		switch result.Severity {
		case enums.SeverityCritical:
			vulnerabilityStats.Critical++
		case enums.SeverityHigh:
			vulnerabilityStats.High++
		case enums.SeverityMedium:
			vulnerabilityStats.Medium++
		case enums.SeverityLow:
			vulnerabilityStats.Low++
		case enums.SeverityInfo:
			vulnerabilityStats.Info++
		}
	}

	// Check if no vulnerabilities found
	totalAlertStats := vulnerabilityStats.Critical + vulnerabilityStats.High + vulnerabilityStats.Medium + vulnerabilityStats.Low + vulnerabilityStats.Info
	if totalAlertStats > 0 {
		vulnerabilityStats.NoVulnerabilities = 0
	}

	templatePath := ""

	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "Add Scan"
	if target.TargetType == enums.TargetTypeIPRange {
		ipCount, err := iputils.ConvertIPRangeToIPSize(target.TargetAddress)
		if err != nil {
			c.FlashAndGoBack(w, r, enums.FlashDanger, "Invalid Target Address")
			return
		}

		records := groupResultsByIP(scanResults)

		templatePath = "scan-results/list-ip-range"

		templateData.Data = ScanResultResponse{
			Target:                              target,
			Records:                             records,
			VulnerabilityStats:                  vulnerabilityStats,
			OverallCVSSScore:                    target.OverallCVSSScore,
			TotalAlerts:                         totalAlerts,
			TotalTargets:                        ipCount.Int64(),
			CVSSScoreByHost:                     target.CVSSScoreByHost,
			DefaultRemediation:                  services.DefaultRemediation,
			NumberOfTARowsForDefaultRemediation: services.NumberOfTARowsForDefaultRemediation,
		}

	} else {
		templatePath = "scan-results/list"
		templateData.Data = ScanResultResponse{
			Target:                              target,
			VulnerabilityStats:                  vulnerabilityStats,
			OverallCVSSScore:                    target.OverallCVSSScore,
			ScanResults:                         scanResults,
			TotalAlerts:                         totalAlerts,
			TotalTargets:                        1,
			DefaultRemediation:                  services.DefaultRemediation,
			NumberOfTARowsForDefaultRemediation: services.NumberOfTARowsForDefaultRemediation,
		}
	}
	if err := views.RenderTemplate(w, templatePath, templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}

func groupResultsByIP(scanResults []datamodels.ScanResult) *map[string][]datamodels.ScanResult {
	records := make(map[string][]datamodels.ScanResult)
	if scanResults == nil {
		return &records
	}

	for _, scanResult := range scanResults {
		ip := scanResult.IP // Assuming each ScanResult has an IP field
		if _, ok := records[ip]; !ok {
			records[ip] = []datamodels.ScanResult{}
		}
		records[ip] = append(records[ip], scanResult)
	}

	// Sort the records by IP
	keys := make([]string, 0, len(records))
	for k := range records {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedRecords := make(map[string][]datamodels.ScanResult)
	for _, k := range keys {
		sortedRecords[k] = records[k]
	}

	return &sortedRecords
}
