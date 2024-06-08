package scan

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/CSPF-Founder/shieldsup/scanner/models"
	"github.com/CSPF-Founder/shieldsup/scanner/schemas"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParseRecord(t *testing.T) {
	// read json from file
	body, err := os.ReadFile("../../testdata/success_scan_results.json")
	if err != nil {
		t.Fatal(err)
	}

	var resp schemas.APIScanResultResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Success != true {
		t.Error("Expected true, got false")
	}

	data := resp.Data
	if data == nil {
		t.Error("Expected data to be non-nil")
	}

	targetID := primitive.NewObjectID()

	// parse the record
	for _, record := range data {
		_, err := parseRecord(record, models.Target{
			ID:               targetID,
			CustomerUsername: "test",
		})
		if err != nil {
			t.Error(err)
		}
	}
}
