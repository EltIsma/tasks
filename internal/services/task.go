package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"task/internal/domain"

	"github.com/google/uuid"
)

type Producer interface {
	Produce(event domain.Event) error
}

type TaskService struct {
	logger   *slog.Logger
	db       Database
	producer Producer
}

func New(logger *slog.Logger, db Database, producer Producer) *TaskService {
	return &TaskService{
		logger:   logger,
		db:       db,
		producer: producer,
	}
}

func (u *TaskService) CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {

	id, err := u.db.CreateTask(ctx, task)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed create task: %w", err)
	}

	return id, nil
}

func (u *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {

	task, err := u.db.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, fmt.Errorf("task doesn't exist: %w", err)
		}
		return nil, fmt.Errorf("failed get task: %w", err)
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

	return task.ID, nil
}

func (u *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {

	err := u.db.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("failed delete task: %w", err)
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
