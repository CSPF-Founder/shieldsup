package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// TODO: Add custom logger

func SendJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal the JSON response %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func SendError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	payload := errResponse{
		Error: msg,
	}

	SendJSON(w, code, payload)
}
