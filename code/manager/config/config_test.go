package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DATABASE_URI", "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local")
	os.Setenv("LOG_FILENAME", "logfile")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("SCANNER_CMD", "scanner")
	os.Setenv("MAIN_DOMAIN", "localhost")

	os.Setenv("USE_DOTENV", "false")
	config := LoadConfig()

	if config.DatabaseURI != "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local" {
		t.Errorf("Expected DatabaseURI to be 'root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local', got %s", config.DatabaseURI)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected LogLevel to be 'info', got %s", config.Logging.Level)
	}

	if config.ScannerCmd != "scanner" {
		t.Errorf("Expected ScannerCmd to be 'scanner', got %s", config.ScannerCmd)
	}

	if config.ScannerDocker != "scanner" {
		t.Errorf("Expected ScannerDocker to be 'scanner', got %s", config.ScannerDocker)
	}

	if config.DatabaseName != "bas_db" {
		t.Errorf("Expected ServerURL to be 'localhost', got %s", config.DatabaseName)
	}

}
