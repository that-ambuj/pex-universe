package routes_test

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"pex-universe/internal/config"
	"pex-universe/internal/server"
	"pex-universe/model/user"
	"pex-universe/routes/v1"

	json "github.com/mixcode/golib-json-snake"
	"gotest.tools/assert"
)

var s *routes.Controller

func TestLogin(t *testing.T) {
	s = Setup(t)
}

func Setup(t *testing.T) *routes.Controller {
	os.Chdir("../..")
	config.LoadEnv()

	host := os.Getenv("DB_HOST")

	// Fail if database host is not localhost
	assert.Assert(t, host == "localhost" || host == "127.0.0.1")

	s := routes.Controller(*server.New())

	s.RegisterAuthRoutes()

	return &s
}

func TestSignup(t *testing.T) {
	body := user.UserSignUpDto{
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "testingPassword#123",
	}

	jsonBytes, err := json.MarshalSnakeCase(body)
	if err != nil {
		t.Fatal(err)
	}

	r := bytes.NewReader(jsonBytes)

	req, err := http.NewRequest("POST", "/signup", r)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := s.App.Test(req)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)

}
