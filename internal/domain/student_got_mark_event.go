package domain

const StudentsGotMarkEventType = "StudentsGotMarkEvent"

type UsersMark struct {
	UserID string `json:"user_id"`
	Mark   int    `json:"mark"`
}

type StudentsGotMarkEvent struct {
	UsersMark []UsersMark `json:"users_mark"`
	TaskID    string      `json:"task_id"`
	LessonID  string      `json:"lesson_id"`
}

func NewStudentsGotMarkEvent(taskResults *TaskResult) *StudentsGotMarkEvent {
	usersMark := make([]UsersMark, 0, len(taskResults.UsersResult))
	for _, ur := range taskResults.UsersResult {
		usersMark = append(usersMark, UsersMark{
			UserID: ur.UserID.String(),
			Mark:   ur.Mark,
		})
	}

	return &StudentsGotMarkEvent{
		UsersMark: usersMark,
		TaskID:    taskResults.TaskID.String(),
		LessonID:  taskResults.LessonID.String(),
	}
}

func (s *StudentsGotMarkEvent) Type() string {
	return StudentsGotMarkEventType
}
