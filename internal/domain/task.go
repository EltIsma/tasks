package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID       uuid.UUID
	Payload  string
	Deadline *time.Time
}

type TaskWithAsignment struct {
	Class    string
	LessonID uuid.UUID
	TaskID   uuid.UUID
	Payload  string
	Deadline *time.Time
}

type ClassLesson struct {
	Class    string
	LessonID uuid.UUID
}

type TaskAsignments struct {
	ToAssign []ClassLesson
	TaskID   uuid.UUID
}

type Assignment struct {
	AssignmentID uuid.UUID
	Class        string
	LessonID     uuid.UUID
}

type TaskAsignment struct {
	AssignmentID uuid.UUID
	Class        string
	Payload      string
}

type LessonTask struct {
	LessonID       uuid.UUID
	TaskID         uuid.UUID
	Payload        string
	Deadline       *time.Time
	TaskTemplateID uuid.UUID
}

type UserResult struct {
	UserID uuid.UUID
	Mark   int
}

type TaskResult struct {
	UsersResult []UserResult
	TaskID      uuid.UUID
	LessonID    uuid.UUID
}
