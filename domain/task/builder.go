package taskdomain

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
)

type TaskBuilder struct {
	Task
}

func NewTaskBuilder() *TaskBuilder {
	return &TaskBuilder{Task{}}
}

func (t *TaskBuilder) SetTitle(ti string) *TaskBuilder {
	t.title = ti
	return t
}

func (t *TaskBuilder) SetDescription(d string) *TaskBuilder {
	t.description = d
	return t
}

func (t *TaskBuilder) SetStatus(s types.Status) *TaskBuilder {
	t.status = s
	return t
}

func (t *TaskBuilder) Build() (*Task, []*errorutils.APIError) {
	if errs := t.validate(); errs != nil {
		return nil, errs
	}

	return &t.Task, nil
}
