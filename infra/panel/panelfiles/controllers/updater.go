package controllers

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	ctx "github.com/CSPF-Founder/shieldsup/onpremise/panel/context"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"
	mid "github.com/CSPF-Founder/shieldsup/onpremise/panel/middlewares"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/views"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

const (
	LAST_UPDATED_DATE_FORMAT = "2-Jan-2006 03:04 PM"
)

type UpdateResponse struct {
	ColorStatus string
	Status      string
	LastUpdate  string
	Updating    bool
}

type updateController struct {
	*App
}

func newUpdateController(app *App) *updateController {
	return &updateController{
		App: app,
	}
}

func (c *updateController) registerRoutes() http.Handler {
	router := chi.NewRouter()

	// Authenticated Routes
	router.Group(func(r chi.Router) {
		r.Use(mid.RequireLogin)

		r.Get("/", c.Index)        // Show update state page
		r.Post("/start", c.Start)  // Start Update
		r.Get("/status", c.Status) // Start Update
	})

	return router
}

// Show Update List
func (c *updateController) Index(w http.ResponseWriter, r *http.Request) {
	user := ctx.Get(r, "user").(models.User)
	data, err := models.FindUpdateStateByUser(user.ID)
	if err != nil {
		c.logger.Error("Error fetching update state", err)
	}

	updateResponse := UpdateResponse{
		ColorStatus: "",
		Status:      "--",
		Updating:    false,
		LastUpdate:  "--",
	}

	if data != nil {
		statusText, err := enums.UpdateStatusMap.GetText(data.Status)
		if err != nil {
			statusText = "Unknown"
			c.logger.Error("Error fetching update status", err)
		}
		updateResponse.ColorStatus = enums.BGColorFromStatus(data.Status)
		updateResponse.Status = statusText
		// Check if the update is in progress
		updateResponse.Updating = data.Status == enums.UpdateStatusUpdating
		if !data.UpdatedAt.IsZero() {
			updateResponse.LastUpdate = data.UpdatedAt.Format(LAST_UPDATED_DATE_FORMAT)
		}
	}

	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "Update State"

	templateData.Data = updateResponse
	if err := views.RenderTemplate(w, "updater/index", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}

func (c *updateController) Start(w http.ResponseWriter, r *http.Request) {
	err := isRsyncFeedUrlReachable(c.config.FeedUrl)
	if err != nil {
		c.logger.Error("Error in rsync feed url", err)
		c.SendJSONError(w, "Unable to reach the Update Server (Ensure all the firewall rules applied as per the manual)", http.StatusInternalServerError)
		return
	}

	user := ctx.Get(r, "user").(models.User)
	existingState, err := models.FindUpdateStateByUser(user.ID)
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			updateState := &models.UpdateState{
				UserID: user.ID,
				Status: enums.UpdateStatusUpdating,
			}
			err = models.SaveUpdateState(updateState)
			if err != nil {
				c.logger.Error("Error in creating update request", err)
				c.SendJSONError(w, "Error in creating update request!", http.StatusInternalServerError)
				return
			}
			c.SendJSONSuccess(w, "Update Requested Successfully")
			return
		}

		c.logger.Error("Error fetching update state", err)
		c.SendJSONError(w, "Error fetching update state")
		return
	}

	if existingState.Status == enums.UpdateStatusUpdating {
		c.SendJSONError(w, "Update already requested!, Please wait for the previous request to finish.")
		return
	}

	if existingState.Status == enums.UpdateStatusUpdated {
		if existingState.GetLastUpdatedDifference() {
			c.SendJSONError(w, fmt.Sprintf("Update request recently completed, Please wait for %d hour before initiating a new request.", models.DELAY_BEFORE_NEW_REQUEST))
			return
		}
	}

	updateState := &models.UpdateState{
		UserID: user.ID,
		Status: enums.UpdateStatusUpdating,
	}
	err = models.SaveUpdateState(updateState)

	if err != nil {
		c.SendJSONError(w, "Error in creating update request!", http.StatusInternalServerError)
		return
	}
	c.SendJSONSuccess(w, "Update Requested Successfully")
}

func (c *updateController) Status(w http.ResponseWriter, r *http.Request) {
	user := ctx.Get(r, "user").(models.User)
	existingState, err := models.FindUpdateStateByUser(user.ID)

	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			c.SendJSONSuccess(w, map[string]any{
				"yet_to_update": true,
			})
			return
		}
		c.logger.Error("Error fetching update state", err)
		c.JSONResponse(w, []string{}, 200)
		return
	}
	c.SendJSONSuccess(w, map[string]any{
		"status":            int(existingState.Status),
		"last_updated_time": existingState.UpdatedAt.Format(LAST_UPDATED_DATE_FORMAT),
	})
}

func parseRsyncURI(uri string) (string, string, error) {
	if len(uri) < 10 || uri[:8] != "rsync://" {
		return "", "", fmt.Errorf("invalid rsync URI format: %s", uri)
	}

	parts := strings.SplitN(uri[8:], ":", 2)
	host := parts[0]

	if len(parts) == 1 {
		// If no port is specified, use the default rsync port
		hostPart := strings.SplitN(host, "/", 2)
		return hostPart[0], "873", nil
	}

	secondParts := strings.SplitN(parts[1], "/", 2)
	var port string
	if len(secondParts) == 2 {
		port = secondParts[0]
	} else {
		port = "873"
	}

	return host, port, nil

}

func isRsyncFeedUrlReachable(feedURL string) error {

	host, port, err := parseRsyncURI(feedURL)
	if err != nil {
		return err
	}

	if host == "" || port == "" {
		return errors.New("invalid rsync URI format")
	}

	// Check if the rsync server is reachable using TCP dial
	// Check TCP connection to the server
	timeout := 5 * time.Second

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
