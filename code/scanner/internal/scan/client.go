package scan

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/CSPF-Founder/shieldsup/scanner/logger"
	"github.com/CSPF-Founder/shieldsup/scanner/schemas"
)

type APIClient struct {
	APIKey       string
	APIServerURL string
	logger       *logger.Logger
}

type CleanUpResponse struct {
	Success bool `json:"success"`
}

type TargetAddResponse struct {
	Success bool `json:"success"`
}

func NewAPIClient(logger *logger.Logger, apiKey string, apiURL string) *APIClient {
	return &APIClient{
		APIKey:       apiKey,
		APIServerURL: apiURL,
		logger:       logger,
	}
}

func (m *APIClient) SendPostRequest(ctx context.Context, urlPath string, data any) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		m.logger.Error("Error marshalling data: %v", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.APIServerURL+urlPath, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+m.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

func (m *APIClient) StartScan(ctx context.Context, targetAddress string) bool {

	urlPath := "/targets"
	data := map[string]string{"target": targetAddress}

	resp, err := m.SendPostRequest(ctx, urlPath, data)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.Error("Error reading response body:", err)
		return false
	}

	var jsonRes TargetAddResponse
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		m.logger.Error("Error decoding JSON response:", err)
		return false
	}

	if !jsonRes.Success {
		m.logger.Error("Unable to start scan", nil)
		m.logger.Error(targetAddress, nil)
		return false
	}

	return jsonRes.Success
}

func (m *APIClient) GetResults(ctx context.Context, targetAddress string, force bool) *schemas.APIScanResultResponse {
	urlPath := "/targets/results"

	// Prepare data
	data := map[string]any{
		"target": targetAddress,
		"force":  force,
	}

	// Send POST request
	resp, err := m.SendPostRequest(ctx, urlPath, data)
	if err != nil {
		m.logger.Error("Error sending POST request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.Error("Error reading response body:", err)
		return nil
	}

	// Decode JSON response
	var result *schemas.APIScanResultResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		m.logger.Error("Error decoding JSON response:", err)
		return nil
	}

	return result
}

func (m *APIClient) Cleanup(ctx context.Context) bool {
	urlPath := "/cleanup"
	data := make(map[string]any)

	resp, err := m.SendPostRequest(ctx, urlPath, data)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.Error("Error reading response body:", err)
		return false
	}

	var jsonRes CleanUpResponse
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		m.logger.Error("Error decoding JSON response:", err)
		return false
	}

	if jsonRes.Success {
		return true
	}

	m.logger.Error("Unable to cleanup", nil)
	return false
}
