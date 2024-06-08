package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CSPF-Founder/shieldsup/scanner-api/logger"

	"github.com/joho/godotenv"
)

// Server represents the Server configuration details
type ServerConfig struct {
	ServerAddress string
	APIKey        string
}

// Config represents the configuration information.
type Config struct {
	Logging    *logger.Config `json:"logging"`
	ServerConf ServerConfig   `json:"server"`
	DockerName string
	// ScannerPath        string
	TemplateFolder     string
	TempTemplateFolder string
	LocalTmpDir        string
}

// Version contains the current project version
var Version = "1"

// ServerName is the server type that is returned in the transparency response.
const ServerName = "shieldsup"

func loadEnv() {
	//determin bin directory and load .env from there
	exe, err := os.Executable()
	if err != nil {
		logger.GetFallBackLogger().Fatal("Error loading .env file", err)
	}
	binDir := filepath.Dir(exe)
	envPath := filepath.Join(binDir, ".env")
	if err := godotenv.Load(envPath); err == nil {
		return
	}

	// try to load .env from current directory
	envPath = ".env"
	if err := godotenv.Load(envPath); err == nil {
		return
	}
	logger.GetFallBackLogger().Error("Error loading .env file", err)
	os.Exit(1)

}

func getEnvValueOrError(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.GetFallBackLogger().Error(fmt.Sprintf("Environment variable %s not set", key), nil)
		os.Exit(1)
	}
	return value
}

func getEnvValueOrDefault(key string, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

// LoadConfig loads the configuration from the specified filepath
func LoadConfig() (*Config, error) {
	if os.Getenv("USE_DOTENV") != "false" {
		loadEnv()
	}

	srvConfig := ServerConfig{
		ServerAddress: getEnvValueOrError("SERVER_ADDRESS"),
		APIKey:        getEnvValueOrError("API_KEY"),
	}

	config := &Config{
		ServerConf:         srvConfig,
		DockerName:         getEnvValueOrError("DOCKER_NAME"),
		TemplateFolder:     getEnvValueOrError("TEMPLATE_FOLDER"),
		TempTemplateFolder: getEnvValueOrError("TEMP_TEMPLATE_FOLDER"),
		LocalTmpDir:        getEnvValueOrError("LOCAL_TMP_DIR"),
		Logging: &logger.Config{
			Level: getEnvValueOrError("LOG_LEVEL"),
			// Log to stdout if no file path is provided
			FilePath: getEnvValueOrDefault("LOG_FILE_PATH", ""),
		},
	}
	return config, nil
}
