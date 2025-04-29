package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"task/internal/domain"
	"time"

	"task/pkg/cache"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Producer interface {
	Produce(event domain.Event) error
}

type TaskService struct {
	logger   *slog.Logger
	db       Database
	cache    cache.Cache
	producer Producer
}

func New(logger *slog.Logger, db Database, cache cache.Cache, producer Producer) *TaskService {
	return &TaskService{
		logger:   logger,
		db:       db,
		cache:    cache,
		producer: producer,
	}
}

func (u *TaskService) CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {

	id, err := u.db.CreateTask(ctx, task)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed create task: %w", err)
	}

	//store in the redis
	rtask, err := json.Marshal(*task)
	if err != nil {
		u.logger.Error("serialize task", slog.String("message", err.Error()))
	}

	err = u.cache.Set(ctx, id, rtask, time.Hour)
	if err != nil {
		u.logger.Error("redis insertion error", slog.String("message", err.Error()))
	}

	return id, nil
}

func (u *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {

	//use trategy cashe aside
	//first check in redis
	redisTask, err := u.cache.Get(ctx, id)
	if err == nil {
		redisTaskBytes, ok := redisTask.(string)
		if !ok {
			u.logger.Error("failed convert redis data to bytes")
		}
		var task domain.Task
		err := json.Unmarshal([]byte(redisTaskBytes), &task)
		if err != nil {
			u.logger.Error("failed convert redis data to domain", slog.String("message", err.Error()))
		}
		if err == nil {
			u.logger.Info("success", slog.String("message", redisTaskBytes))
			return &task, nil
		}
	} else {
		if !errors.Is(err, redis.Nil) {
			u.logger.Error("redis error", slog.String("message", err.Error()))
		}
	}

	task, err := u.db.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, fmt.Errorf("task doesn't exist: %w", err)
		}
		return nil, fmt.Errorf("failed get task: %w", err)
	}

	//store in the redis
	rtask, err := json.Marshal(task)
	if err != nil {
		u.logger.Error("serialize task", slog.String("message", err.Error()))
	}

	if err == nil {
		err = u.cache.Set(ctx, id, rtask, time.Hour)
		if err != nil {
			u.logger.Error("redis insertion error", slog.String("message", err.Error()))
		}
	}

	return task, nil
}

func (u *TaskService) GetTasks(ctx context.Context) ([]*domain.Task, error) {
	task, err := u.db.GetTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get task: %w", err)
	}

	return task, nil
}

func (u *TaskService) UpdateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {

	err := u.db.UpdateTask(ctx, task)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed update task: %w", err)
	}

	rtask, err := json.Marshal(*task)
	if err != nil {
		u.logger.Error("serialize task", slog.String("message", err.Error()))
	}

	err = u.cache.Set(ctx, task.ID, rtask, time.Hour)
	if err != nil {
		u.logger.Error("redis insertion error", slog.String("message", err.Error()))
	}

	return task.ID, nil
}

func (u *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {

	err := u.db.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("failed delete task: %w", err)
	}

	err = u.cache.Del(ctx, id)
	if err != nil {
		u.logger.Error("delete from redis", slog.String("message", err.Error()))
	}

	return nil
}

func (u *TaskService) CreateAssignments(ctx context.Context, taskAssignments *domain.TaskAsignments) ([]domain.Assignment, error) {
	assignments, err := u.db.CreateAssignments(ctx, taskAssignments)
	if err != nil {
		return nil, fmt.Errorf("failed assignment task to users task: %w", err)
	}

	for _, event := range domain.NewTaskAssignedToUserEvent(assignments) {
		err = u.producer.Produce(event)
		if err != nil {
			u.logger.Error("failed to send event:", slog.String("error", err.Error()))
		}
	}

	return assignments, err
}

func (u *TaskService) GetTaskByClass(ctx context.Context, class string) ([]*domain.LessonTask, error) {
	tasks, err := u.db.GetTaskByClass(ctx, class)
	if err != nil {
		return nil, fmt.Errorf("failed get task: %w", err)
	}

	return tasks, err
}

func (u *TaskService) SetTaskResultsByUsers(ctx context.Context, taskResults *domain.TaskResult) error {
	err := u.db.SetTaskResultsByUsers(ctx, taskResults)
	if err != nil {
		return fmt.Errorf("failed assignment task to users task: %w", err)
	}

	err = u.producer.Produce(domain.NewStudentsGotMarkEvent(taskResults))
	if err != nil {
		u.logger.Error("failed to send event:", slog.String("error", err.Error()))
	}

	return nil
}

func (u *TaskService) DeleteAssignment(ctx context.Context, assignmentID uuid.UUID) error {
	err := u.db.DeleteAssignment(ctx, assignmentID)
	if err != nil {
		return fmt.Errorf("failed assignment task to users task: %w", err)
	}

	return nil
}

func (u *TaskService) CreateTaskWithAssignments(ctx context.Context, assignment *domain.TaskWithAsignment) (uuid.UUID, error) {
	id, err := u.db.CreateTaskWithAssignments(ctx, assignment)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed create task with assignment: %w", err)
	}

	events := []domain.Assignment{
		{
			AssignmentID: id,
			Class:        assignment.Class,
			LessonID:     assignment.LessonID,
		},
	}

	domainEvents := domain.NewTaskAssignedToUserEvent(events)

	err = u.producer.Produce(domainEvents[0])
	if err != nil {
		u.logger.Error("failed to send event:", slog.String("error", err.Error()))
	}

	return id, nil
}

func (u *TaskService) UpdateAssignment(ctx context.Context, assignment *domain.TaskAsignment) error {
	err := u.db.UpdateAssignment(ctx, assignment)
	if err != nil {
		return fmt.Errorf("failed create task with assignment: %w", err)
	}

	return nil
}
