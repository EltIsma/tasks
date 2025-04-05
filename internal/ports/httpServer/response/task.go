package response

import (
	"task/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID       string     `json:"id"`
	Payload  string     `json:"payload" binding:"required"`
	Deadline *time.Time `json:"deadline,omitempty" example:"2025-01-01T13:00:00Z"`
}

func NewTaskResponse(task *domain.Task) *Task {
	response := &Task{
		ID:      task.ID.String(),
		Payload: task.Payload,
	}
	if task.Deadline != nil {
		response.Deadline = task.Deadline
	}

	return response
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

func NewTasksResponse(tasks []*domain.Task) *Tasks {
	var t Tasks
	t.Tasks = make([]Task, 0, len(tasks))
	for _, task := range tasks {
		response := &Task{
			ID:      task.ID.String(),
			Payload: task.Payload,
		}
		if task.Deadline != nil {
			response.Deadline = task.Deadline
		}

		t.Tasks = append(t.Tasks, *response)
	}

	return &t
}

type LessonTask struct {
	LessonID       string     `json:"lesson_id"`
	TaskID         string     `json:"task_id"`
	Payload        string     `json:"payload" binding:"required"`
	Deadline       *time.Time `json:"deadline,omitempty" example:"2025-01-01T13:00:00Z"`
	TaskTemplateID string     `json:"task_template_id"`
}

type TaskID struct {
	ID string `json:"id"`
}

func NewTaskIDResponse(id uuid.UUID) *TaskID {
	return &TaskID{
		ID: id.String(),
	}
}

type UserTask struct {
	UserID   string    `json:"user_id"`
	LessonID string    `json:"lesson_id"`
	TaskID   string    `json:"task_id"`
	Payload  string    `json:"payload"`
	Deadline time.Time `json:"deadline,omitempty" example:"2025-01-01T13:00:00Z"`
}

type ClassTasks struct {
	Class string       `json:"class"`
	Tasks []LessonTask `json:"tasks"`
}

func NewClassTasksResponse(class string, lt []*domain.LessonTask) *ClassTasks {
	lessonTask := make([]LessonTask, 0, len(lt))
	for _, domainLessonTask := range lt {
		lessonTask = append(lessonTask, LessonTask{
			LessonID:       domainLessonTask.LessonID.String(),
			TaskID:         domainLessonTask.TaskID.String(),
			Payload:        domainLessonTask.Payload,
			Deadline:       domainLessonTask.Deadline,
			TaskTemplateID: domainLessonTask.TaskTemplateID.String(),
		})
	}
	return &ClassTasks{
		Class: class,
		Tasks: lessonTask,
	}
}

type AssignmentID struct {
	ID string `json:"class_task_id"`
}

func NewAssignmentIDResponse(assignmentID string) *AssignmentID {
	return &AssignmentID{
		ID: assignmentID,
	}
}

// type TaskAssignments struct {
// 	Users    []string  `json:"users"`
// 	LessonID string    `json:"lesson_id" binding:"required"`
// 	TaskID   string    `json:"task_id"`
// 	Payload  string    `json:"payload"`
// 	Deadline time.Time `json:"deadline,omitempty" example:"2025-01-01T13:00:00Z"`
// }

// type ClassLesson struct {
// 	Class    string `json:"class" binding:"required"`
// 	LessonID string `json:"lesson_id" binding:"required"`
// }

// type TaskAssignments struct {
// 	ToAssign     []ClassLesson `json:"class_lessons" binding:"required"`
// 	AssignmentID string        `json:"assignment_id" binding:"required"`
// }

// func NewTaskAssignmentsResponse(a *domain.TaskAsignments) *TaskAssignments {
// 	classLessons := make([]ClassLesson, 0, len(a.ToAssign))
// 	for _, cl := range a.ToAssign {
// 		classLessons = append(classLessons, ClassLesson{
// 			Class:    cl.Class,
// 			LessonID: cl.LessonID.String(),
// 		})
// 	}

// 	return &TaskAssignments{
// 		AssignmentID: a.TaskID.String(),
// 		ToAssign:     classLessons,
// 	}
// }

type Assignments struct {
	AssignmentID string `json:"class_task_id" binding:"required"`
	Class        string `json:"class" binding:"required"`
	LessonID     string `json:"lesson_id" binding:"required"`
}

type TaskAssignments struct {
	ToAssign []Assignments `json:"assignments" binding:"required"`
	TaskID   string        `json:"task_template_id" binding:"required"`
}

func NewAssignmentsResponse(taskID string, a []domain.Assignment) *TaskAssignments {
	assignments := make([]Assignments, 0, len(a))
	for _, val := range a {
		assignments = append(assignments, Assignments{
			AssignmentID: val.AssignmentID.String(),
			Class:        val.Class,
			LessonID:     val.LessonID.String(),
		})
	}

	return &TaskAssignments{
		ToAssign: assignments,
		TaskID:   taskID,
	}
}
