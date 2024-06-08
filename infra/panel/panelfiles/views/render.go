package views

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/config"
	ctx "github.com/CSPF-Founder/shieldsup/onpremise/panel/context"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/services"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/sessions"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:embed templates
var templateFS embed.FS

type TemplateData struct {
	Version                string
	ProductTitle           string
	CopyrightFooterCompany string
	Title                  string
	Flashes                []sessions.SessionFlash
	User                   models.User
	Data                   any
	CSRFToken              string
	CSRFName               string

	HideHeaderAndFooter bool
	CurrentYear         int

	PreviousPage string

	// StringMap map[string]string
	// IntMap     map[string]int
	// FloatMap   map[string]float32
	// CSSVersion string
}

var functions = template.FuncMap{
	"assetPath":            AssetPath,
	"formatDate":           formatDate,
	"formatNormalDate":     formatNormalDate,
	"BgClassBySeverity":    services.GetBgClassBySeverity,
	"CalculateTextAreaRow": services.CalculateTextAreaRow,
	"Severity":             enums.SeverityToString,
	"HexID":                getObjectIDString,
	"ConvertJSONToString":  ConvertJSONToString,
}

// views.NewTemplateData returns the default template parameters for a user and
// the CSRF token.
func NewTemplateData(conf *config.Config, session *sessions.SessionManager, r *http.Request) *TemplateData {
	checkUser := ctx.Get(r, "user")
	user := models.User{}
	if checkUser != nil {
		user = ctx.Get(r, "user").(models.User)
	}
	year, _, _ := time.Now().Date()

	return &TemplateData{
		CSRFToken:              session.GetCSRF(r.Context()),
		User:                   user,
		Version:                config.Version,
		Flashes:                session.Flashes(r.Context()),
		ProductTitle:           conf.ProductTitle,
		CopyrightFooterCompany: conf.CopyrightFooterCompany,
		CurrentYear:            year,
		CSRFName:               conf.ServerConf.CSRFName,
		PreviousPage:           utils.GetRelativePath(r),
	}
}

/**
* parseTemplate parses the template and returns it
**/
func parseTemplate(
	page string,
	templateToRender string,
) (t *template.Template, err error) {

	t, err = template.New(page).Funcs(
		functions,
	).ParseFS(
		templateFS,
		"templates/layout/base.tmpl",
		"templates/layout/header.tmpl",
		"templates/layout/footer.tmpl",
		templateToRender,
		"templates/layout/flashes.tmpl",
	)

	if err != nil {
		return nil, fmt.Errorf("Error parsing template %s", err)
	}

	// app.templateCache[templateToRender] = t
	return t, nil
}

/**
* renderTemplate renders the template
**/
func RenderTemplate(
	w http.ResponseWriter,
	page string,
	td *TemplateData,
	// partials ...string,
) (err error) {
	templateToRender := fmt.Sprintf("templates/%s.tmpl", page)

	if td == nil {
		td = &TemplateData{}
	}

	t, err := parseTemplate(
		"base",
		templateToRender,
	)

	if err != nil {
		return fmt.Errorf("Error parsing template %s", err)
	}

	err = t.Execute(w, td)
	if err != nil {
		return fmt.Errorf("Error executing template %s", err)
	}

	return nil
}

// Custom function to format Mongo date
func formatDate(date primitive.DateTime, layout string) string {
	return date.Time().Format(layout)
}

// Custom function to format tim.Time
func formatNormalDate(date time.Time, layout string) string {
	return date.Format(layout)
}

// Custom function to format Mongo date
func getObjectIDString(id primitive.ObjectID) string {
	return id.Hex()
}

// ConvertJSONToString converts JSON to string
func ConvertJSONToString(data any) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonBytes)
}
