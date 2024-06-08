package controllers

import (
	"net/http"

	"github.com/CSPF-Founder/shieldsup/scanner-api/helpers"
)

func HealthChecker(w http.ResponseWriter, r *http.Request) {
	type isAliveResponse struct {
		IsAlive bool `json:"is_alive"`
	}

	helpers.SendJSON(w, http.StatusOK, &isAliveResponse{IsAlive: true})
}
