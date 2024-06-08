package scanner

import (
	"os"
	"testing"
)

func prepareTestEnv() {
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("LOCAL_TMP_DIR", "/tmp")
	os.Setenv("SERVER_ADDRESS", "0.0.0.0:5000")
	os.Setenv("API_KEY", "test")
	os.Setenv("DOCKER_NAME", "test")
	os.Setenv("TEMPLATE_FOLDER", "test")
	os.Setenv("TEMP_TEMPLATE_FOLDER", "test")
}

func TestParse(t *testing.T) {
	prepareTestEnv()

	// TODO: implement this test

}
