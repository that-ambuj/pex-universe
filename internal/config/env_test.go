package config_test

import (
	"os"
	"pex-universe/internal/config"
	"testing"

	"gotest.tools/assert"
)

func TestEnv(t *testing.T) {
	os.Chdir("../..")
	config.LoadEnv()

	assert.Equal(t, os.Getenv("APP_ENV"), "test")
	assert.Equal(t, os.Getenv("DB_HOST"), "localhost")
	assert.Equal(t, os.Getenv("DB_PORT"), "3306")
	assert.Equal(t, os.Getenv("DB_DATABASE"), "pexuniverse")
	assert.Equal(t, os.Getenv("DB_USERNAME"), "pex_test_user")
	assert.Equal(t, os.Getenv("DB_PASSWORD"), "somepassword")
}
