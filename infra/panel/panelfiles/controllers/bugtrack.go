package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	ctx "github.com/CSPF-Founder/shieldsup/onpremise/panel/context"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/repositories/datarepos"
	mid "github.com/CSPF-Founder/shieldsup/onpremise/panel/middlewares"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/utils"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/views"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BugTrackDetailReponse struct {
	Data                 models.BugTrack
	BugTrackSeverity     enums.BTSeverityMapType
	BugTrackStatus       enums.BTStatusMapType
	PrioritizationStatus map[enums.Prioritization]string
}

type BugTrackAddReponse struct {
	Target               string
	AlertTitle           string
	Details              string
	Severity             enums.Severity
	Poc                  string
	Remediation          string
	BugTrackSeverity     enums.BTSeverityMapType
	BugTrackStatus       enums.BTStatusMapType
	PrioritizationStatus enums.PrioritizationMapType
	TodayDate            time.Time
}

type bugTrackController struct {
	*App
	scanResultRepo datarepos.ScanResultRepository
	targetRepo     datarepos.TargetRepository
}

func newBugTrackController(
	app *App,
	scanResultRepo datarepos.ScanResultRepository,
	targetRepo datarepos.TargetRepository,
) *bugTrackController {
	return &bugTrackController{
		App:            app,
		scanResultRepo: scanResultRepo,
		targetRepo:     targetRepo,
	}
}

func (c *bugTrackController) registerRoutes() http.Handler {
	router := chi.NewRouter()

	// Authenticated Routes
	router.Group(func(r chi.Router) {
		r.Use(mid.RequireLogin)

		r.Get("/", c.List)                                  // Display list
		r.Post("/export-bugs", c.ExportBugsAction)          // Export bug track entries
		r.Get("/add-from-scanresult", c.AddFromScanResult)  // Display add from scan result form
		r.Post("/add-from-scanresult", c.AddFromScanResult) // Add from scan result action

		r.Route("/{bugID:[0-9]+}", func(r chi.Router) {
			r.Get("/", c.Details)           // Display detail page
			r.Delete("/", c.DeleteBug)      // Delete Bug action
			r.Patch("/", c.UpdateBugAction) // Update bug action

		})
	})

	return router
}

// Show List
func (c *bugTrackController) List(w http.ResponseWriter, r *http.Request) {
	user := ctx.Get(r, "user").(models.User)
	data, _ := models.GetOverviewListByUser(models.OverViewParameters{
		User: user,
	})
	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "BugTrack"
	templateData.Data = data

	if err := views.RenderTemplate(w, "bug-track/index", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}

}

// Show Details
func (c *bugTrackController) Details(w http.ResponseWriter, r *http.Request) {
	bugID := chi.URLParam(r, "bugID")
	user := ctx.Get(r, "user").(models.User)
	data, err := models.FindBugTrackByIdAndUser(bugID, user)
	if err != nil {
		c.FlashAndGoBack(w, r, enums.FlashDanger, "Invalid Request")
		return
	}
	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "BugTrack"
	templateData.Data = BugTrackDetailReponse{
		Data:                 data,
		BugTrackSeverity:     enums.BTSeverityMap,
		BugTrackStatus:       enums.BTStatusMap,
		PrioritizationStatus: enums.PrioritizationMap,
	}

	if err := views.RenderTemplate(w, "bug-track/details", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}

func (c *bugTrackController) UpdateBugAction(w http.ResponseWriter, r *http.Request) {
	bugID := chi.URLParam(r, "bugID")

	user := ctx.Get(r, "user").(models.User)
	bug, err := models.FindBugTrackByIdAndUser(bugID, user)
	if err != nil {
		c.SendJSONError(w, "Invalid bug id")
		return
	}

	requiresClarification, _ := strconv.Atoi(r.FormValue("requires_clarification"))
	newRemarks := r.FormValue("new_remarks")

	// Handle requires_clarification
	if requiresClarification == 1 {
		bug.ClarificationStatus = requiresClarification

		if newRemarks == "" {
			c.SendJSONError(w, "Please fill the remarks column with clarification query")
			return
		}
	}

	// Handle new_remarks
	if newRemarks := newRemarks; newRemarks != "" {
		remarks := "User: " + newRemarks
		if bug.Remarks != "" {
			bug.Remarks = fmt.Sprintf("%s\n%s", bug.Remarks, remarks)
		} else {
			bug.Remarks = remarks
		}
	}

	if status, ok := r.Form["status"]; ok {
		st, _ := strconv.Atoi(status[0])

		bug.Status, err = enums.BTStatusMap.ByIndex(st)
		if err != nil {
			c.SendJSONError(w, "Invalid status")
			return
		}
	}

	if toBeFixedBy, ok := r.Form["to_be_fixed_by"]; ok {
		bug.ToBeFixedBy = toBeFixedBy[0]
	}

	if severity, ok := r.Form["severity"]; ok {
		bug.Severity, err = enums.BTSeverityMap.ByIndex(severity[0])
		if err != nil {
			c.SendJSONError(w, "Invalid severity")
			return
		}
	}

	if prioritization, ok := r.Form["prioritization"]; ok {
		bug.Prioritization, err = enums.PrioritizationMap.ByIndex(prioritization[0])
		if err != nil {
			c.SendJSONError(w, "Invalid prioritization")
			return
		}
	}

	if details, ok := r.Form["details"]; ok {
		bug.Details = details[0]
	}

	if poc, ok := r.Form["poc"]; ok {
		bug.Poc = poc[0]
	}

	if remediation, ok := r.Form["remediation"]; ok {
		bug.Remediation = remediation[0]
	}

	// Update the bug
	if err := models.SaveBugTrack(&bug); err != nil {
		c.SendJSONSuccess(w, "No changes made")
		return
	}
	c.SendJSONSuccess(w, "BugTrack Updated")
}

func (c *bugTrackController) DeleteBug(w http.ResponseWriter, r *http.Request) {
	bugID := chi.URLParam(r, "bugID")

	// Get the bug object by ID and user
	user := ctx.Get(r, "user").(models.User)
	bug, err := models.FindBugTrackByIdAndUser(bugID, user)
	if err != nil {
		c.SendJSONError(w, "Invalid Request")
		return
	}

	err = models.DeleteBugTrack(&bug)
	if err != nil {
		c.SendJSONError(w, "Unable to delete the entry")
		return
	}
	c.SendJSONSuccess(w, "Successfully deleted the entry")
}

func (c *bugTrackController) ExportBugsAction(w http.ResponseWriter, r *http.Request) {
	user := ctx.Get(r, "user").(models.User)
	bugList, err := models.GetDetailBugTrackByUser(user)

	if err != nil {
		c.FlashAndGoBack(w, r, enums.FlashWarning, "No Bugs Found")
		return
	}

	if len(bugList) > 0 {
		f := excelize.NewFile()

		// Set column widths
		f.SetColWidth("Sheet1", "A", "A", 42)
		f.SetColWidth("Sheet1", "B", "B", 40)
		f.SetColWidth("Sheet1", "C", "C", 15) // Auto size
		f.SetColWidth("Sheet1", "D", "D", 40)
		f.SetColWidth("Sheet1", "E", "E", 40)
		f.SetColWidth("Sheet1", "F", "F", 40)
		f.SetColWidth("Sheet1", "G", "G", 30)
		f.SetColWidth("Sheet1", "H", "H", 20)
		f.SetColWidth("Sheet1", "I", "I", 20) // Auto size
		f.SetColWidth("Sheet1", "J", "J", 20) // Auto size
		f.SetColWidth("Sheet1", "K", "K", 20) // Auto size
		f.SetColWidth("Sheet1", "L", "L", 20) // Auto size

		// Set header values
		headers := []string{
			"Url/IP/Application",
			"Alert",
			"Severity",
			"Details/Impact",
			"Replication/Proof",
			"Remediation",
			"Remarks",
			"To Be Fixed By",
			"Found Date",
			"Revalidated Date",
			"Prioritization",
			"Status",
		}
		for i, header := range headers {
			cell := fmt.Sprintf("%s1", string(rune(65+i)))
			f.SetCellValue("Sheet1", cell, header)
		}

		// Style header
		styleHeader, _ := f.NewStyle(`{
			"font": {"bold": true, "color": "#ffffff"},
			"fill": {"type": "pattern", "color": ["#FAA61A"], "pattern":1},
			"alignment": {"horizontal": "center", "vertical": "center"}
		}`)
		f.SetCellStyle("Sheet1", "A1", "L1", styleHeader)

		x := 2
		for _, bug := range bugList {
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", x), bug.Target)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", x), bug.AlertTitle)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", x), bug.SeverityText)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", x), bug.Details)
			f.SetCellValue("Sheet1", fmt.Sprintf("E%d", x), bug.Poc)
			f.SetCellValue("Sheet1", fmt.Sprintf("F%d", x), bug.Remediation)
			f.SetCellValue("Sheet1", fmt.Sprintf("G%d", x), bug.Remarks)
			f.SetCellValue("Sheet1", fmt.Sprintf("H%d", x), bug.ToBeFixedBy)
			f.SetCellValue("Sheet1", fmt.Sprintf("I%d", x), bug.FormatedFoundDate)
			f.SetCellValue("Sheet1", fmt.Sprintf("J%d", x), bug.FormatedRevalidatedDate)
			f.SetCellValue("Sheet1", fmt.Sprintf("K%d", x), bug.PrioritizationText)
			f.SetCellValue("Sheet1", fmt.Sprintf("L%d", x), bug.StatusText)

			// Style severity
			severityColor := "#FF0000"
			switch bug.Severity {
			case enums.BTSeverityCritical:
				severityColor = "#E83123"
			case enums.BTSeverityHigh:
				severityColor = "#E77f34"
			case enums.BTSeverityMedium:
				severityColor = "#e6ac30"
			case enums.BTSeverityLow:
				severityColor = "#2fa84d"
			case enums.BtSeverityRecommendation:
				severityColor = "#0773b8"
			}
			styleSeverity, _ := f.NewStyle(`{
				"font": {"bold": true, "color": "#FFFFFF"},
				"fill": {"type": "pattern", "color": ["` + severityColor + `"], "pattern":1},
				"alignment":{"horizontal":"center","vertical":"center","wrap_text":true}
			}`)

			f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", x), fmt.Sprintf("C%d", x), styleSeverity)

			// Style Row I to L
			styleRow, _ := f.NewStyle(`{
				"alignment": {"horizontal": "center", "vertical": "center","wrap_text":true}
			}`)
			f.SetCellStyle("Sheet1", fmt.Sprintf("I%d", x), fmt.Sprintf("L%d", x), styleRow)

			// Style Row
			styleRow2, _ := f.NewStyle(`{
				"alignment": {"wrap_text":true, "vertical": "center"}
			}`)
			f.SetCellStyle("Sheet1", fmt.Sprintf("A%d", x), fmt.Sprintf("B%d", x), styleRow2)
			f.SetCellStyle("Sheet1", fmt.Sprintf("D%d", x), fmt.Sprintf("H%d", x), styleRow2)

			// Status validation & style
			f.AddDataValidation("Sheet1", &excelize.DataValidation{
				Type:         "list",
				Formula1:     `"Unfixed,Fixed"`,
				ShowDropDown: true,
				Sqref:        fmt.Sprintf("L%d", x),
			})
			// f.SetCellDataValidation("Sheet1", fmt.Sprintf("L%d", x), "list", `"Unfixed,Fixed"`)
			// f.SetCellStyle("Sheet1", fmt.Sprintf("L%d", x), fmt.Sprintf("L%d", x), styleHeader)

			f.SetRowHeight("Sheet1", x, 120)
			x++
		}

		// Set sheet title
		f.SetSheetName("Sheet1", "Vulnerabilities")

		// Write to response
		w.Header().Set("Content-Disposition", "attachment; filename=Vulnerabilities.xlsx")
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		if err := f.Write(w); err != nil {
			fmt.Println(err)
		}
	} else {
		c.FlashAndGoBack(w, r, enums.FlashWarning, "No Bugs Found")
		return
	}
}

// addFromScanResultAction handles the logic for adding from scan result action
func (b *bugTrackController) AddFromScanResult(w http.ResponseWriter, r *http.Request) {
	user := ctx.Get(r, "user").(models.User)
	if r.Method == "GET" {
		requiredParams := []string{"target_id", "alert_id"}

		if !utils.CheckAllParamsExist(r, requiredParams) {
			b.FlashAndGoBack(w, r, enums.FlashDanger, "Invalid Request")
			return
		}

		targetID := r.FormValue("target_id")
		alertID := r.FormValue("alert_id")

		var err error
		targetObjectID, err := primitive.ObjectIDFromHex(targetID)
		if err != nil {
			b.FlashAndGoBack(w, r, enums.FlashDanger, "Invalid Target Id")
			return
		}

		alertObjectID, err := primitive.ObjectIDFromHex(alertID)
		if err != nil {
			b.FlashAndGoBack(w, r, enums.FlashDanger, "Invalid Alert Id")
			return
		}

		scanResult, err := b.scanResultRepo.ByID(r.Context(), alertObjectID)
		if err != nil {

			b.FlashAndGoBack(w, r, enums.FlashDanger, "Invalid Alert Id")
			return
		}

		alertTitle := scanResult.VulnerabilityTitle
		target, err := b.targetRepo.ByIdAndCustomerUsername(r.Context(), targetID, user.Username)
		if err != nil {

			b.FlashAndGoBack(w, r, enums.FlashDanger, "Invalid Target Id")
			return
		}

		targetType := target.TargetType
		bugtrackTarget := target.TargetAddress

		var ipList []string
		groupByAlert := r.FormValue("group_by_alert") == "1"

		if targetType == enums.TargetTypeIPRange {
			if groupByAlert {
				ipList, _ = b.scanResultRepo.IPListByAlert(r.Context(), alertTitle, targetObjectID)
			}

			if len(ipList) > 1 {
				bugtrackTarget = "Multiple IPs"
			} else if host := scanResult.Host; host != "" {
				bugtrackTarget = host
			} else if ip := scanResult.IP; ip != "" {
				bugtrackTarget = ip
			}
		}

		var poc, details, remediation string

		details = scanResult.VulnerabilityDescription

		if groupByAlert && len(ipList) > 1 {
			details += "\nAffected IPs:\n" + strings.Join(ipList, "\n")
		}

		severity := scanResult.Severity
		remediation = scanResult.Remediation

		if evidence := scanResult.Evidence; evidence != "" {
			poc = "Evidence:\n" + evidence
		}

		if matchedAt := scanResult.MatchedAt; matchedAt != "" {
			poc += "\nEvidence Found At:\n" + matchedAt
		}

		// Get the current date and time
		now := time.Now()

		templateData := views.NewTemplateData(b.config, b.session, r)
		templateData.Title = "Add BugTrack"
		templateData.Data = BugTrackAddReponse{
			Target:               bugtrackTarget,
			AlertTitle:           alertTitle,
			Details:              details,
			Severity:             severity,
			Poc:                  poc,
			Remediation:          remediation,
			BugTrackSeverity:     enums.BTSeverityMap,
			BugTrackStatus:       enums.BTStatusMap,
			PrioritizationStatus: enums.PrioritizationMap,
			TodayDate:            now,
		}

		if err := views.RenderTemplate(w, "bug-track/add", templateData); err != nil {
			b.logger.Error("Error rendering template: ", err)
		}

	} else if r.Method == "POST" {
		requiredParams := []string{"target", "alert_title", "status", "severity", "prioritization"}

		if !utils.CheckAllParamsExist(r, requiredParams) {
			b.FlashAndGoBack(w, r, enums.FlashDanger, "Please fill all the inputs")
			return
		}

		target := r.PostFormValue("target")
		alertTitle := r.PostFormValue("alert_title")
		details := r.PostFormValue("details")
		poc := r.PostFormValue("poc")
		toBeFixedBy := r.PostFormValue("to_be_fixed_by")

		remediation := r.PostFormValue("remediation")
		remarks := r.PostFormValue("remarks")
		foundDate := r.PostFormValue("found_date")
		revalidatedDate := r.PostFormValue("revalidated_date")

		severity, err := enums.BTSeverityMap.ByIndex(r.PostFormValue("severity"))
		if err != nil {
			b.SendJSONError(w, "Invalid severity")
			return
		}

		status, err := enums.BTStatusMap.ByIndex(r.PostFormValue("status"))
		if err != nil {
			b.SendJSONError(w, "Invalid status")
			return
		}

		prioritization, err := enums.PrioritizationMap.ByIndex(r.PostFormValue("prioritization"))
		if err != nil {
			b.SendJSONError(w, "Invalid prioritization")
			return
		}

		count, err := models.CheckBugTrackAlreadyExists(user.ID, target, severity, alertTitle, details, poc)

		if err != nil {
			fmt.Println(err)
			b.SendJSONError(w, "Error checking bugtrack!")
			return
		}

		if count > 0 {
			b.SendJSONError(w, "Already Exists!")
			return
		}

		ToBeFixedBy := ""

		if toBeFixedBy != "" {
			ToBeFixedBy = toBeFixedBy
		}

		bug := models.BugTrack{
			UserID:           user.ID,
			AlertTitle:       alertTitle,
			Target:           target,
			Status:           status,
			Severity:         severity,
			Details:          details,
			Prioritization:   prioritization,
			TestingMethod:    enums.TestingMethodAutomatic,
			EffortsToExploit: enums.EffortsToExploitNotApplicable,
			DataLeakage:      enums.DataLeakageNotApplicable,
			CanWafStop:       enums.CanWafStopNotApplicable,
			Likelihood:       enums.LikelihoodNotApplicable,
			ToBeFixedBy:      ToBeFixedBy,
		}

		if poc != "" {
			bug.Poc = poc
		} else {
			bug.Poc = ""
		}

		if remediation != "" {
			bug.Remediation = remediation
		} else {
			bug.Remediation = ""
		}

		if remarks != "" {
			bug.Remarks = remarks
		} else {
			bug.Remarks = ""
		}

		if foundDate != "" {
			date, _ := time.Parse("2006-01-02", foundDate)
			bug.FoundDate = date
		} else {
			now := time.Now()
			bug.FoundDate = now
		}

		if revalidatedDate != "" {
			date, _ := time.Parse("2006-01-02", revalidatedDate)
			bug.RevalidatedDate = date
		}

		err = models.SaveBugTrack(&bug)
		if err != nil {
			b.SendJSONError(w, "Unable to add bugtrack")
			return
		}
		b.SendJSONSuccess(w, "BugTrack Added")
	}
}
