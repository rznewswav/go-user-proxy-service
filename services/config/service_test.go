package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadStruct(t *testing.T) {
	os.Setenv("APP_NAME", "service_test")
	type AppConfig struct {
		AppName  string `env:"APP_NAME"             default:"service_test" printDebug:"true"`
		AppEnv   string `env:"APP_ENV"              default:"development"         printDebug:"true"`
		LogLevel string `env:"LOG_LEVEL"            default:"info"                printDebug:"true"`
	}
	loadedConfig, loadConfigError := Build(AppConfig{})
	if loadConfigError != nil {
		t.Fatal(loadConfigError)
	}

	assert.Equal(t, loadedConfig.AppName, "service_test")
}
