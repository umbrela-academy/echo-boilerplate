package tasks_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/handlers/tasks"
	"github.com/alexferl/echo-boilerplate/handlers/users"
)

func TestHandler_CreateTask_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.CreateTaskRequest{
		Title: "My Title",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	newTask := tasks.NewTask()
	newTask.Create(user.Id)
	task := newTask.MakeResponse(user, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"Insert",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskResponse{task},
			nil,
		)

	s.ServeHTTP(resp, req)

	var result tasks.TaskResponse
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, task.Id, result.Id)
	assert.Equal(t, task.CreatedBy.Id, result.CreatedBy.Id)
}

func TestHandler_CreateTask_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_CreateTask_422(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.CreateTaskRequest{
		Title: "",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	var result tasks.ListTasksResponse
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func createTasks(num int, user *users.User) []*tasks.TaskResponse {
	var result []*tasks.TaskResponse

	for i := 1; i <= num; i++ {
		newTask := tasks.NewTask()
		newTask.Create(user.Id)
		task := newTask.MakeResponse(user, nil, nil)
		result = append(result, task)
	}

	return result
}

func TestHandler_ListTasks_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	retTasks := createTasks(10, user)

	req := httptest.NewRequest(http.MethodGet, "/tasks?per_page=1&page=2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"Count",
			mock.Anything,
			mock.Anything,
		).
		Return(
			int64(10),
			nil,
		).
		On(
			"Aggregate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			retTasks,
			nil,
		)

	s.ServeHTTP(resp, req)

	var result tasks.ListTasksResponse
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	h := resp.Header()
	link := `<http://example.com/tasks?per_page=1&page=3>; rel=next, ` +
		`<http://example.com/tasks?per_page=1&page=10>; rel=last, ` +
		`<http://example.com/tasks?per_page=1&page=1>; rel=first, ` +
		`<http://example.com/tasks?per_page=1&page=1>; rel=prev`

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, 10, len(result.Tasks))
	assert.Equal(t, "2", h.Get("X-Page"))
	assert.Equal(t, "1", h.Get("X-Per-Page"))
	assert.Equal(t, "10", h.Get("X-Total"))
	assert.Equal(t, "10", h.Get("X-Total-Pages"))
	assert.Equal(t, "3", h.Get("X-Next-Page"))
	assert.Equal(t, "1", h.Get("X-Prev-Page"))
	assert.Equal(t, link, h.Get("Link"))
}

func TestHandler_ListTasks_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
