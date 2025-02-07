package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	app "github.com/alexferl/echo-boilerplate"
	"github.com/alexferl/echo-boilerplate/handlers"
	_ "github.com/alexferl/echo-boilerplate/testing"
)

func TestHandler_Root(t *testing.T) {
	s := app.NewTestServer(handlers.NewHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Welcome")
}
