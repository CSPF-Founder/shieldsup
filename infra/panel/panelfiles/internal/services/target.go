package services

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/repositories/datarepos"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models/datamodels"
)

func DownloadReport(target datamodels.Target, w http.ResponseWriter, r *http.Request) error {
	reportPath, err := target.GetReportPath()
	if err != nil {
		return err
	}
	reportID := target.ID.Hex()

	_, err = os.Stat(reportPath)

	if err == nil {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"ShieldsUP_report_%s.docx\"", reportID))
		http.ServeFile(w, r, reportPath)
		return nil
	} else {
		return err
	}
}

func DeleteTarget(ctx context.Context, target datamodels.Target, targetRepo datarepos.TargetRepository) bool {
	if !target.CanDelete() {
		return false
	}

	reportDir, err := target.GetReportDir()
	if err != nil {
		return false
	}

	if _, err := os.Stat(reportDir); err == nil {
		// Delete report folder
		os.RemoveAll(reportDir)
	}

	// Delete target
	deletedCount, err := targetRepo.DeleteTargetByID(ctx, target.ID)
	if err != nil {
		return false
	}

	return deletedCount > 0
}
