package datamodels

import (
	"time"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Info struct {
	MetaData map[string]any `json:"metadata" bson:"metadata"`
	Author   []string       `json:"author" bson:"author"`
}

type Classification struct {
	CVSSScore   float64  `json:"cvss_score" bson:"cvss_score"`
	CVEID       []string `json:"cve_id" bson:"cve_id"`
	CWEID       []string `json:"cwe_id" bson:"cwe_id"`
	CVSSMetrics string   `json:"cvss_metrics" bson:"cvss_metrics"`
}

type ScanResult struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CustomerName  string             `bson:"customer_username"`
	Info          Info               `json:"info" bson:"info"`
	Template      string             `json:"template" bson:"template"`
	TemplateURL   string             `json:"template_url" bson:"template_url"`
	TemplateID    string             `json:"template_id" bson:"template_id"`
	TemplatePath  string             `json:"template_path" bson:"template_path"`
	Type          string             `json:"type" bson:"type"`
	Host          string             `json:"host" bson:"host"`
	IP            string             `json:"ip" bson:"ip"`
	CURLCommand   string             `json:"curl_command" bson:"curl_command"`
	ExtractorName string             `json:"extractor_name" bson:"extractor_name"`
	Scheme        string             `json:"scheme" bson:"scheme"`
	URL           string             `json:"url" bson:"url"`
	Path          string             `json:"path" bson:"path"`
	Request       string             `json:"request" bson:"request"`
	Response      string             `json:"response" bson:"response"`
	Metadata      map[string]any     `json:"meta" bson:"meta"`
	Timestamp     time.Time          `json:"timestamp" bson:"timestamp"`
	MatcherStatus bool               `json:"matcher_status" bson:"matcher_status"`
	MatchedLine   []int              `json:"matched_line" bson:"matched_line"`
	MatchedAt     string             `json:"matched_at" bson:"matched_at"`

	// Parsed fields
	Severity                 enums.Severity     `json:"severity" bson:"severity"`
	SeverityText             string             `json:"severity_text" bson:"severity_text"`
	VulnerabilityTitle       string             `json:"vulnerability_title" bson:"vulnerability_title"`
	VulnerabilityDescription string             `json:"vulnerability_description" bson:"vulnerability_description"`
	Tags                     []string           `json:"tags" bson:"tags"`
	TargetID                 primitive.ObjectID `json:"target_id" bson:"target_id"`
	Reference                []string           `json:"reference" bson:"reference"`
	Remediation              string             `json:"remediation" bson:"remediation"`
	Classification           Classification     `json:"classification" bson:"classification"`
	Evidence                 string             `json:"evidence" bson:"evidence"`
}
