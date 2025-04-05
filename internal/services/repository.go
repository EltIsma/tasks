package services

import (
	"context"
	"task/internal/domain"

	"github.com/google/uuid"
)

type Database interface {
	CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error)
	GetTaskByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetTasks(ctx context.Context) ([]*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
	CreateAssignments(ctx context.Context, taskAssignments *domain.TaskAsignments) (assignments []domain.Assignment, err error)
	GetTaskByClass(ctx context.Context, class string) ([]*domain.LessonTask, error)
	SetTaskResultsByUsers(ctx context.Context, taskResults *domain.TaskResult) error
	DeleteAssignment(ctx context.Context, assignmentID uuid.UUID) error
	CreateTaskWithAssignments(ctx context.Context, assignment *domain.TaskWithAsignment) (uuid.UUID, error)
	UpdateAssignment(ctx context.Context, task *domain.TaskAsignment) error
}
