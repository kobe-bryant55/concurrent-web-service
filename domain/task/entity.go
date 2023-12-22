package taskdomain

import (
	"errors"

	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
)

type Task struct {
	id          uint64
	title       string
	description string
	status      types.Status
}

func (t *Task) GetID() uint64 {
	return t.id
}

func (t *Task) GetDescription() string {
	return t.description
}

func (t *Task) GetTitle() string {
	return t.title
}

func (t *Task) GetStatus() types.Status {
	return t.status
}

func (t *Task) SetStatus(s types.Status) []*errorutils.APIError {
	t.status = s
	if errs := t.validate(); errs != nil {
		return errs
	}

	return nil
}

func (t *Task) validate() []*errorutils.APIError {
	var errs []*errorutils.APIError

	if len(t.description) == 0 {
		errs = append(errs, errorutils.New(errors.New("description required"), nil))
	} else if len(t.description) < 5 {
		errs = append(errs, errorutils.New(errors.New("description too short"), nil))
	} else if len(t.description) > 500 {
		errs = append(errs, errorutils.New(errors.New("description too long"), nil))
	}

	if len(t.title) == 0 {
		errs = append(errs, errorutils.New(errors.New("title required"), nil))
	} else if len(t.title) < 3 {
		errs = append(errs, errorutils.New(errors.New("title too short"), nil))
	} else if len(t.description) > 50 {
		errs = append(errs, errorutils.New(errors.New("title too long"), nil))
	}

	if t.status == "" {
		errs = append(errs, errorutils.New(errors.New("status cannot be empty"), nil))
	} else if t.status != types.Active && t.status != types.Passive {
		errs = append(errs, errorutils.New(errors.New("invalid status"), nil))
	}

	return errs
}

func (t *Task) ToDTO() *dto.TaskResponse {
	return &dto.TaskResponse{
		ID:          t.id,
		Description: t.description,
		Title:       t.title,
		Status:      t.status,
	}
}
