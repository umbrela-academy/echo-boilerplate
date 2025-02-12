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

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/tasks"
	"github.com/alexferl/echo-boilerplate/handlers/users"
)

func TestHandler_GetTask_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	newTask := tasks.NewTask()
	newTask.Create(user.Id)
	task := newTask.MakeResponse(user, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"Aggregate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskResponse{task},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_GetTask_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_GetTask_404(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"Aggregate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskResponse{},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandler_GetTask_410(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	newTask := tasks.NewTask()
	newTask.Create(user.Id)
	newTask.Delete(user.Id)
	task := newTask.MakeResponse(user, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"Aggregate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskResponse{task},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusGone, resp.Code)
}

func TestHandler_UpdateTask_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.UpdateTaskRequest{
		Title:     "My Edited Task",
		Completed: true,
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	newTask := tasks.NewTask()
	newTask.Create(user.Id)
	newTask.CreatedBy = user.Id
	newTask.Title = payload.Title
	newTask.Complete(user.Id)

	updated := newTask.MakeResponse(user, nil, user)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			newTask,
			nil,
		).
		On(
			"UpdateById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskResponse{updated},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_UpdateTask_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_UpdateTask_403(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.UpdateTaskRequest{
		Title: "My Edited Task",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create("another_id")

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			task,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestHandler_UpdateTask_404(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.UpdateTaskRequest{
		Title: "My Edited Task",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			tasks.ErrTaskNotFound,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandler_UpdateTask_410(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &tasks.UpdateTaskRequest{
		Title: "My Edited Task",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create(user.Id)
	task.Delete(user.Id)
	find := &tasks.Task{
		Model: &data.Model{
			DeletedAt: task.DeletedAt,
		},
	}

	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			find,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusGone, resp.Code)
}

func TestHandler_UpdateTask_422(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	b := bytes.NewBuffer([]byte(`{"invalid": "invalid"}`))
	req := httptest.NewRequest(http.MethodPatch, "/tasks/id", b)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestHandler_DeleteTask_204(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create(user.Id)
	find := &tasks.Task{
		Model: &data.Model{
			CreatedBy: user.Id,
			DeletedAt: nil,
		},
		Title:       "",
		Completed:   false,
		CompletedAt: task.CompletedAt,
		CompletedBy: task.CompletedBy,
	}
	task.Delete(user.Id)
	update := &tasks.TaskResponse{
		Id:        task.Id,
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
		DeletedBy: user.Id,
		DeletedAt: user.DeletedAt,
	}

	req := httptest.NewRequest(http.MethodDelete, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			find,
			nil,
		).
		On(
			"UpdateById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskResponse{update},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestHandler_DeleteTask_403(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	task := tasks.NewTask()
	task.Create("another_id")

	req := httptest.NewRequest(http.MethodDelete, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			task,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestHandler_DeleteTask_404(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			tasks.ErrTaskNotFound,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandler_DeleteTask_410(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	newTask := tasks.NewTask()
	newTask.Create(user.Id)
	find := &tasks.Task{
		Model: &data.Model{
			CreatedBy: user.Id,
		},
	}
	newTask.Delete(user.Id)
	task := newTask.MakeResponse(user, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			find,
			nil,
		).
		On(
			"UpdateById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			[]*tasks.TaskResponse{task},
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}
