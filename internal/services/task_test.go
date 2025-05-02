package services_test

import (
	"context"
	"encoding/json"
	"task/internal/app"
	"task/internal/domain"
	"task/internal/services"
	repoMock "task/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	ctx := context.Background()
	mockService := new(repoMock.Database)
	cacheMock := new(repoMock.Cache)
	producerMock := new(repoMock.Producer)
	id := uuid.New()
	task := &domain.Task{
		ID:      id,
		Payload: "5+5 = ?",
	}
	rtask, err := json.Marshal(*task)
	require.NoError(t, err)
	mockService.On("CreateTask", ctx, task).Return(id, nil)
	cacheMock.On("Set", ctx, id, rtask, time.Hour).Return(nil)
	logger := app.InitLogger()
	usecase := services.New(logger, mockService, cacheMock, producerMock)

	taskID, err := usecase.CreateTask(ctx, task)

	assert.Equal(t, id, taskID)

	assert.NoError(t, err)
}

func TestUpdateTask(t *testing.T) {
	ctx := context.Background()
	mockService := new(repoMock.Database)
	cacheMock := new(repoMock.Cache)
	producerMock := new(repoMock.Producer)
	id := uuid.New()
	task := &domain.Task{
		ID:      id,
		Payload: "5+5 = ?",
	}
	rtask, err := json.Marshal(*task)
	require.NoError(t, err)
	mockService.On("UpdateTask", ctx, task).Return(nil)
	cacheMock.On("Set", ctx, id, rtask, time.Hour).Return(nil)
	logger := app.InitLogger()
	usecase := services.New(logger, mockService, cacheMock, producerMock)

	taskID, err := usecase.UpdateTask(ctx, task)

	assert.Equal(t, id, taskID)

	assert.NoError(t, err)
}

func TestDeleteTask(t *testing.T) {
	ctx := context.Background()
	mockService := new(repoMock.Database)
	cacheMock := new(repoMock.Cache)
	producerMock := new(repoMock.Producer)
	id := uuid.New()

	mockService.On("DeleteTask", ctx, id).Return(nil)
	cacheMock.On("Del", ctx, id).Return(nil)
	logger := app.InitLogger()
	usecase := services.New(logger, mockService, cacheMock, producerMock)

	err := usecase.DeleteTask(ctx, id)

	assert.NoError(t, err)
}

func TestGetTaskByClass(t *testing.T) {
	ctx := context.Background()
	mockService := new(repoMock.Database)
	cacheMock := new(repoMock.Cache)
	producerMock := new(repoMock.Producer)
	class := "9A"

	mockService.On("GetTaskByClass", ctx, class).Return([]*domain.LessonTask{
		{
			LessonID:       uuid.New(),
			TaskID:         uuid.New(),
			Payload:        "??",
			TaskTemplateID: uuid.New(),
		},
	}, nil)
	logger := app.InitLogger()
	usecase := services.New(logger, mockService, cacheMock, producerMock)

	_, err := usecase.GetTaskByClass(ctx, class)

	assert.NoError(t, err)
}

func TestCreateTaskWithAssignment(t *testing.T) {
	ctx := context.Background()
	mockService := new(repoMock.Database)
	cacheMock := new(repoMock.Cache)
	producerMock := new(repoMock.Producer)
	class := "9A"

	id := uuid.New()
	assignment := &domain.TaskWithAsignment{
		Class:    class,
		LessonID: uuid.New(),
		TaskID:   uuid.New(),
		Payload:  "what?",
	}

	events := []domain.Assignment{
		{
			AssignmentID: id,
			Class:        assignment.Class,
			LessonID:     assignment.LessonID,
		},
	}

	domainEvents := domain.NewTaskAssignedToUserEvent(events)

	mockService.On("CreateTaskWithAssignments", ctx, assignment).Return(id, nil)
	producerMock.On("Produce", domainEvents[0]).Return(nil)
	logger := app.InitLogger()
	usecase := services.New(logger, mockService, cacheMock, producerMock)

	assignmentID, err := usecase.CreateTaskWithAssignments(ctx, assignment)
	assert.Equal(t, id, assignmentID)
	assert.NoError(t, err)
}

func TestUpdateAssignment(t *testing.T) {
	ctx := context.Background()
	mockService := new(repoMock.Database)
	cacheMock := new(repoMock.Cache)
	producerMock := new(repoMock.Producer)
	class := "9A"

	id := uuid.New()
	assignment := &domain.TaskAsignment{
		AssignmentID: id,
		Class:        class,
		Payload:      "what?",
	}

	mockService.On("UpdateAssignment", ctx, assignment).Return(nil)
	logger := app.InitLogger()
	usecase := services.New(logger, mockService, cacheMock, producerMock)

	err := usecase.UpdateAssignment(ctx, assignment)
	assert.NoError(t, err)
}

func TestDeletAssignment(t *testing.T) {
	ctx := context.Background()
	mockService := new(repoMock.Database)
	cacheMock := new(repoMock.Cache)
	producerMock := new(repoMock.Producer)

	id := uuid.New()
	mockService.On("DeleteAssignment", ctx, id).Return(nil)
	logger := app.InitLogger()
	usecase := services.New(logger, mockService, cacheMock, producerMock)

	err := usecase.DeleteAssignment(ctx, id)
	assert.NoError(t, err)
}
