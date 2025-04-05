package domain

const TaskAssignedToClassEventType = "TaskAssignedToClass"

type TaskAssignmentToClassEvent struct {
	Class    string `json:"class"`
	LessonID string `json:"lesson_id"`
	TaskID   string `json:"task_id"`
}

func NewTaskAssignedToUserEvent(assignments []Assignment) []*TaskAssignmentToClassEvent {
	events := make([]*TaskAssignmentToClassEvent, 0, len(assignments))
	for _, a := range assignments {
		events = append(events, &TaskAssignmentToClassEvent{
			Class:    a.Class,
			LessonID: a.LessonID.String(),
			TaskID:   a.AssignmentID.String(),
		})
	}

	return events
}

func (s *TaskAssignmentToClassEvent) Type() string {
	return TaskAssignedToClassEventType
}
