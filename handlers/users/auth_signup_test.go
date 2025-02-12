package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/handlers/users"
)

func TestHandler_Auth_Signup_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	payload := &users.AuthSignUpRequest{
		Email:    "test@example.com",
		Username: "test",
		Name:     "Test",
		Password: "abcdefghijkl",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			nil,
		).
		On(
			"Upsert",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			nil,
		)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_Auth_Signup_409(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	payload := &users.AuthSignUpRequest{
		Email:    "test@example.com",
		Username: "test",
		Name:     "Test",
		Password: "abcdefghijkl",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			&users.User{},
			nil,
		)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusConflict, resp.Code)
}

func TestHandler_Auth_Signup_422(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer([]byte(`{"invalid": "key"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}
