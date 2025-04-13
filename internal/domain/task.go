package domain

type (
	TaskStatus string
	TaskType   string
)

const (
	TaskStatusCreated    TaskStatus = "created"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusSuccess    TaskStatus = "success"
	TaskStatusFailed     TaskStatus = "failed"

	TaskTypeFirst  TaskType = "first"
	TaskTypeSecond TaskType = "second"
	TaskTypeThird  TaskType = "third"
)

func (t TaskStatus) String() string {
	return string(t)
}

func (t TaskType) String() string {
	return string(t)
}

func (t TaskType) IsValid() bool {
	switch t {
	case TaskTypeFirst, TaskTypeSecond, TaskTypeThird:
		return true
	default:
		return false
	}
}
