package config

import (
	"os"
	"testing"
)

func prepareTestEnv() {
	os.Setenv("LOG_LEVEL", "info")
}

func TestLoadConfig(t *testing.T) {
	prepareTestEnv()

	os.Setenv("USE_DOTENV", "false")
	config, err := LoadConfig()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected LogLevel to be 'info', got %s", config.Logging.Level)
	}
}
