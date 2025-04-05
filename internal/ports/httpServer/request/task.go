package request

import (
	"fmt"
	"task/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Payload  string    `json:"payload" binding:"required"`
	Deadline time.Time `json:"deadline,omitempty" example:"2025-01-01T13:00:00Z"`
}

func (t Task) ToDomain() *domain.Task {
	task := &domain.Task{
		ID:      uuid.New(),
		Payload: t.Payload,
	}

	if !t.Deadline.IsZero() {
		task.Deadline = &t.Deadline
	}

	return task
}

func (t Task) ToDomainWithID(id uuid.UUID) *domain.Task {
	task := &domain.Task{
		ID:      id,
		Payload: t.Payload,
	}

	if !t.Deadline.IsZero() {
		task.Deadline = &t.Deadline
	}

	return task
}

type TaskAsignment struct {
	AssignmentID string `json:"class_task_id" binding:"required"`
	Class        string `json:"class" binding:"required"`
	Payload      string `json:"payload" binding:"required"`
}

func (t TaskAsignment) ToDomain() (*domain.TaskAsignment, error) {
	assignmentID, err := uuid.Parse(t.AssignmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid assignment id = %s with error: %w", t.AssignmentID, err)
	}

	result := &domain.TaskAsignment{
		AssignmentID: assignmentID,
		Class:        t.Class,
		Payload:      t.Payload,
	}

	return result, nil
}

type ClassLesson struct {
	Class    string `json:"class" binding:"required"`
	LessonID string `json:"lesson_id" binding:"required"`
}

type TaskAsignments struct {
	ToAssign []ClassLesson `json:"assign_to" binding:"required"`
	TaskID   string        `json:"template_task_id" binding:"required"`
}

func (t TaskAsignments) ToDomain() (*domain.TaskAsignments, error) {
	taskID, err := uuid.Parse(t.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task id = %s with error: %w", t.TaskID, err)
	}

	classLesson := make([]domain.ClassLesson, 0, len(t.ToAssign))
	for _, cl := range t.ToAssign {
		lessonID, err := uuid.Parse(cl.LessonID)
		if err != nil {
			return nil, fmt.Errorf("invalid lesson id = %s with error: %w", cl.LessonID, err)
		}
		classLesson = append(classLesson, domain.ClassLesson{
			Class:    cl.Class,
			LessonID: lessonID,
		})
	}

	return &domain.TaskAsignments{
		TaskID:   taskID,
		ToAssign: classLesson,
	}, nil
}

type TaskWithAsignment struct {
	Class    string    `json:"class" binding:"required"`
	LessonID string    `json:"lesson_id" binding:"required"`
	Payload  string    `json:"payload" binding:"required"`
	Deadline time.Time `json:"deadline,omitempty" example:"2025-01-01T13:00:00Z"`
}

func (t TaskWithAsignment) ToDomain() (*domain.TaskWithAsignment, error) {

	lessonID, err := uuid.Parse(t.LessonID)
	if err != nil {
		return nil, fmt.Errorf("invalid lesson id = %s with error: %w", t.LessonID, err)
	}

	taskWithAssignment := &domain.TaskWithAsignment{
		Class:    t.Class,
		LessonID: lessonID,
		TaskID:   uuid.New(),
		Payload:  t.Payload,
	}

	if !t.Deadline.IsZero() {
		taskWithAssignment.Deadline = &t.Deadline
	}

	return taskWithAssignment, nil
}

type TaskResult struct {
	UsersResult []UserResult `json:"users_result" binding:"required"`
	TaskID      string       `json:"task_id" binding:"required"`
	LessonID    string       `json:"lesson_id" binding:"required"`
}

func (t TaskResult) ToDomain() (*domain.TaskResult, error) {
	usersResult := make([]domain.UserResult, 0, len(t.UsersResult))

	taskID, err := uuid.Parse(t.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task id = %s with error: %w", t.TaskID, err)
	}

	lessonID, err := uuid.Parse(t.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid lesson id = %s with error: %w", t.TaskID, err)
	}

	for _, ur := range t.UsersResult {
		userID, err := uuid.Parse(ur.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user id = %s with error: %w", ur.UserID, err)
		}
		usersResult = append(usersResult, domain.UserResult{
			UserID: userID,
			Mark:   ur.Mark,
		})
	}

	return &domain.TaskResult{
		UsersResult: usersResult,
		TaskID:      taskID,
		LessonID:    lessonID,
	}, nil
}

type UserResult struct {
	UserID string `json:"user_id" binding:"required"`
	Mark   int    `json:"mark" binding:"required"`
}

type Class struct {
	Class string `form:"class"`
}

type TaskAsignmentID struct {
	AssignmentID string `json:"class_task_id" binding:"required"`
}

func (t TaskAsignmentID) ToUUID() (uuid.UUID, error) {
	assignmentID, err := uuid.Parse(t.AssignmentID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid assignment id = %s with error: %w", t.AssignmentID, err)
	}

	return assignmentID, nil
}
