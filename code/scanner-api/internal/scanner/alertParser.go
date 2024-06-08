package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/CSPF-Founder/shieldsup/scanner-api/logger"
)

type ParsedData map[string]any

type AlertsParser struct {
	Alerts []ParsedData
}

func NewAlertsParser() *AlertsParser {
	return &AlertsParser{
		Alerts: []ParsedData{},
	}
}

func (a *AlertsParser) Parse(inputFile string, logger *logger.Logger) []ParsedData {
	// file, err := os.Open(inputFile)  // this was not working
	file, err := os.OpenFile(inputFile, os.O_RDONLY, 0666)
	if err != nil {
		logger.Error("Error opening file in Parser: ", err)
		return a.Alerts
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		var alertEntry map[string]any
		err := json.Unmarshal(line, &alertEntry)
		if err != nil {
			logger.Error("Error Unmarshalling data in Parser: ", err)
			continue
		}

		info, ok := alertEntry["info"].(map[string]any)
		if !ok {
			continue
		}

		severity, ok := info["severity"].(string)
		fmt.Println(severity)
		if !ok {
			continue
		}

		alertTitle, ok := info["name"].(string)
		fmt.Println(alertTitle)
		if !ok {
			continue
		}

		a.Alerts = append(a.Alerts, alertEntry)
	}

	if err := scanner.Err(); err != nil {
		return a.Alerts
	}
	return a.Alerts
}
