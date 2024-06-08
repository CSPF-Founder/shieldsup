package schemas

import (
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner/enums"
)

type APIScanResultResponse struct {
	Data       []InputRecord       `json:"data"`
	ScanStatus enums.APIScanStatus `json:"scan_status"`
	Success    bool                `json:"success"`
}

// !-- Input Record Schema --! //
type InputRecord struct {
	Info          InputInfo      `json:"info"`
	Template      string         `json:"template"`
	TemplateURL   string         `json:"template-url"`
	TemplateID    string         `json:"template-id"`
	TemplatePath  string         `json:"template-path"`
	Type          string         `json:"type"`
	Host          string         `json:"host"`
	IP            string         `json:"ip"`
	CURLCommand   string         `json:"curl-command"`
	ExtractorName string         `json:"extractor-name"`
	Scheme        string         `json:"scheme"`
	URL           string         `json:"url"`
	Path          string         `json:"path"`
	Request       string         `json:"request"`
	Response      string         `json:"response"`
	Metadata      map[string]any `json:"meta"`
	Timestamp     time.Time      `json:"timestamp"`
	MatcherStatus bool           `json:"matcher-status"`
	MatchedLine   []int          `json:"matched-line"`
	MatchedAt     string         `json:"matched-at"`

	// Input only fields
	MatcherName      string   `json:"matcher-name"`
	ExtractedResults []string `json:"extracted-results"`
	// ReqURLPattern string `json:"req_url_pattern"`
	// Error               string         `json:"error"`
}

type InputInfo struct {
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	Severity       string              `json:"severity"`
	Tags           []string            `json:"tags"`
	Reference      []string            `json:"reference"`
	Remediation    string              `json:"remediation"`
	Classification InputClassification `json:"classification"`
	Impact         string              `json:"impact"`
	MetaData       map[string]any      `json:"metadata"`
	Author         []string            `json:"author"`
}

type InputClassification struct {
	CVSSScore   float64  `json:"cvss-score"`
	CVSSMetrics string   `json:"cvss-metrics"`
	CVEID       []string `json:"cve-id"`
	CWEID       []string `json:"cwe-id"`
}
