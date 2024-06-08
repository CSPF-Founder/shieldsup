package config

import (
	"os"
	"strconv"
	"testing"
)

func prepareTestEnv() {
	os.Setenv("PRODUCT_TITLE", "Shields Up")
	os.Setenv("SERVER_ADDRESS", "0.0.0.0:8080")
	os.Setenv("DATABASE_URI", "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local")
	os.Setenv("DBMS_TYPE", "mysql")
	os.Setenv("COPYRIGHT_FOOTER_COMPANY", "CySecurity Pte Ltd")
	os.Setenv("WORK_DIR", "/app/data/")
	os.Setenv("TEMP_UPLOADS_DIR", "/app/data/temp_uploads/")
	os.Setenv("MIGRATIONS_PREFIX", "db")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("USE_TLS", "true")
	os.Setenv("CERT_PATH", "/app/certs/panel.crt")
	os.Setenv("KEY_PATH", "/app/certs/panel.key")
	os.Setenv("FEED_URL", "http://localhost")
	os.Setenv("FEED_CONNECTION_CHECK_TIMEOUT", "1200")
	os.Setenv("MONGO_DATABASE_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DATABASE_NAME", "shieldsup")
	os.Setenv("SHIELDSUP_REPORT_DIR", "/tmp/")

}

func TestLoadConfig(t *testing.T) {
	prepareTestEnv()

	os.Setenv("USE_DOTENV", "false")
	config, err := LoadConfig()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if config.ProductTitle != "Shields Up" {
		t.Errorf("Expected ProductTitle to be 'Shields Up', got %s", config.ProductTitle)
	}

	if config.ServerConf.ServerAddress != "0.0.0.0:8080" {
		t.Errorf("Expected ServerAddress to be '0.0.0.0:8080', got %s", config.ServerConf.ServerAddress)
	}

	if config.DatabaseURI != "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local" {
		t.Errorf("Expected DatabaseURI to be 'root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local', got %s", config.DatabaseURI)
	}

	if config.DBMSType != "mysql" {
		t.Errorf("Expected DBMSType to be 'mysql', got %s", config.DBMSType)
	}

	if config.CopyrightFooterCompany != " Cyber Security and Privacy Foundation " {
		t.Errorf("Expected CopyrightFooterCompany to be ' Cyber Security and Privacy Foundation ', got %s", config.CopyrightFooterCompany)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected LogLevel to be 'info', got %s", config.Logging.Level)
	}

	useTLS, err := strconv.ParseBool(os.Getenv("USE_TLS"))
	if err != nil {
		t.Errorf("Error parsing USE_TLS: %v", err)
	}
	if config.ServerConf.UseTLS != useTLS {
		t.Errorf("Expected UseTLS to be %t, got %t", useTLS, config.ServerConf.UseTLS)
	}

	if config.ServerConf.CertPath != "/app/certs/panel.crt" {
		t.Errorf("Expected CertPath to be '/app/certs/panel.crt', got %s", config.ServerConf.CertPath)
	}

	if config.ServerConf.KeyPath != "/app/certs/panel.key" {
		t.Errorf("Expected KeyPath to be '/app/certs/panel.key', got %s", config.ServerConf.KeyPath)
	}

	if config.SheildsUpReportDir != "/var/www/html" {
		t.Errorf("Expected SheildsUpReportDir to be '/var/www/html', got %s", config.SheildsUpReportDir)
	}

	// if config.FeedUrl != "http://localhost" {
	// 	t.Errorf("Expected FeedUrl to be 'http://localhost', got %s", config.FeedUrl)
	// }

	// if config.FeedConnectionCheckTimeout != "1200" {
	// 	t.Errorf("Expected FeedUrl to be '1200', got %s", config.FeedConnectionCheckTimeout)
	// }
}
