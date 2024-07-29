package scan

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/CSPF-Founder/shieldsup/scanner/enums"
	"github.com/CSPF-Founder/shieldsup/scanner/internal/repositories"
	"github.com/CSPF-Founder/shieldsup/scanner/logger"
	"github.com/CSPF-Founder/shieldsup/scanner/models"
	"github.com/CSPF-Founder/shieldsup/scanner/schemas"
)

type ResultsHandlerModule struct {
	logger         *logger.Logger
	Target         models.Target
	reporterBin    string
	scanResultRepo *repositories.ScanResultRepository
}

func ResultsHandler(
	logger logger.Logger,
	reporterBin string,
	target models.Target,
	scanResultRepo *repositories.ScanResultRepository,
) *ResultsHandlerModule {
	return &ResultsHandlerModule{
		logger:         &logger,
		Target:         target,
		reporterBin:    reporterBin,
		scanResultRepo: scanResultRepo,
	}
}

func (r *ResultsHandlerModule) Run(ctx context.Context, results []schemas.InputRecord) error {
	r.AddToDB(results)

	r.logger.Info("Making Report")
	scriptArgs := []string{"-m", "reporter", "-t", r.Target.ID.Hex()}
	cmd := exec.CommandContext(ctx, r.reporterBin, scriptArgs...)

	// capture the stdout and stderr
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	// Run the command
	err := cmd.Run()
	if err != nil {
		errMsg := ""
		errMsg = err.Error()

		if stdOut.String() != "" {
			errMsg = errMsg + "\n" + stdOut.String()
		}

		if stdErr.String() != "" {
			errMsg = errMsg + "\n" + stdErr.String()
		}
		return errors.New(errMsg)
	}

	if stdErr.Len() > 0 {
		return fmt.Errorf("Reporter script returned an error: %s", stdErr.String())
	}

	return nil
}

// addToDB converts results to records and adds them to the database
func (r *ResultsHandlerModule) AddToDB(results []schemas.InputRecord) {

	records := make([]models.ScanResult, 0, len(results))

	for _, entry := range results {
		record, err := parseRecord(entry, r.Target)
		if err != nil {
			log.Println("Error converting entry to record", err)
			continue
		}

		records = append(records, *record)
	}

	if len(records) == 0 {
		r.logger.Error("No records to add to the database", nil)
		return
	}

	r.scanResultRepo.AddMany(records)

}

// parseRecord converts an input record to an output record
func parseRecord(input schemas.InputRecord, target models.Target) (*models.ScanResult, error) {
	var record models.ScanResult

	record.VulnerabilityTitle = input.Info.Name
	if record.VulnerabilityTitle == "" {
		return nil, errors.New("Error converting entry to record: missing or invalid 'name' field")
	}

	record.VulnerabilityDescription = input.Info.Description
	// if record.VulnerabilityDescription == "" {
	// 	return nil, errors.New("Error converting entry to record: missing or invalid 'description' field")
	// }

	record.Severity = enums.SeverityInfo // Default to Info

	if input.Info.Severity != "" {
		if strings.ToLower(input.Info.Severity) == "unknown" {
			record.Severity = enums.SeverityMedium
			record.SeverityText = "Medium"
		} else {
			severity := enums.SeverityFromString(input.Info.Severity)
			severityText, err := enums.SeverityMap.GetText(severity)
			if err != nil {
				return nil, errors.New("Error converting entry to record: invalid 'severity' field")
			}
			record.Severity = severity
			record.SeverityText = severityText
		}
	}

	record.Template = input.Template
	record.TemplateURL = input.TemplateURL
	record.TemplateID = input.TemplateID
	record.TemplatePath = input.TemplatePath
	record.Type = input.Type
	record.Host = input.Host
	record.IP = input.IP
	record.CURLCommand = input.CURLCommand
	record.ExtractorName = input.ExtractorName
	record.Scheme = input.Scheme
	record.URL = input.URL
	record.Path = input.Path
	record.Request = input.Request
	record.Response = input.Response
	record.Timestamp = input.Timestamp
	record.MatcherStatus = input.MatcherStatus
	record.MatchedLine = input.MatchedLine
	record.MatchedAt = input.MatchedAt

	record.Tags = input.Info.Tags
	record.Reference = input.Info.Reference
	record.Remediation = input.Info.Remediation
	record.Classification.CVSSScore = input.Info.Classification.CVSSScore
	record.Classification.CVEID = input.Info.Classification.CVEID
	record.Classification.CWEID = input.Info.Classification.CWEID

	record.CustomerName = target.CustomerUsername

	// Info fields
	record.Info.MetaData = input.Info.MetaData
	record.Info.Author = input.Info.Author

	evidence := ""
	if input.MatcherName != "" {
		evidence += input.MatcherName
	}

	if len(input.ExtractedResults) > 0 {
		if evidence != "" {
			evidence += "\n"
		}

		for _, extractedResult := range input.ExtractedResults {
			evidence += extractedResult + "\n"
		}
	}

	record.Evidence = evidence

	record.TargetID = target.ID
	return &record, nil
}
