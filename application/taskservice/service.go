package taskservice

import (
	taskdomain "github.com/MehmetTalhaSeker/concurrent-web-service/domain/task"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
)

type Service interface {
	Create(*dto.TaskCreateRequest) error
	Read(*dto.RequestWithID) (*dto.TaskResponse, error)
	Reads() ([]*dto.TaskResponse, error)
	Update(*dto.TaskUpdateRequest) (*dto.TaskResponse, error)
	Delete(*dto.RequestWithID) (*dto.ResponseWithID, error)
}

type service struct {
	repository taskdomain.IRepository
}

func NewService(repository taskdomain.IRepository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(req *dto.TaskCreateRequest) error {
	t, errs := taskdomain.NewTaskBuilder().SetDescription(req.Description).SetTitle(req.Title).SetStatus(types.Active).Build()
	if errs != nil {
		return errorutils.ValidationError(errs)
	}

	err := s.repository.Create(t)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Read(req *dto.RequestWithID) (*dto.TaskResponse, error) {
	tid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	t, err := s.repository.Read(*tid)
	if err != nil {
		return nil, err
	}

	return t.ToDTO(), nil
}

func (s *service) Reads() ([]*dto.TaskResponse, error) {
	tasks, err := s.repository.Reads()
	if err != nil {
		return nil, err
	}

	var usr []*dto.TaskResponse

	for _, u := range tasks {
		usr = append(usr, u.ToDTO())
	}

	return usr, nil
}

func (s *service) Update(req *dto.TaskUpdateRequest) (*dto.TaskResponse, error) {
	tid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	t, err := s.repository.Read(*tid)
	if err != nil {
		return nil, err
	}

	if errs := t.SetStatus(req.Status); errs != nil {
		return nil, errorutils.ValidationError(errs)
	}

	if err = s.repository.Update(t); err != nil {
		return nil, err
	}

	return t.ToDTO(), nil
}

func (s *service) Delete(req *dto.RequestWithID) (*dto.ResponseWithID, error) {
	tid, err := apputils.StringToUINT64(req.ID)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrInvalidID, err)
	}

	_, err = s.repository.Read(*tid)
	if err != nil {
		return nil, err
	}

	if err = s.repository.Delete(*tid); err != nil {
		return nil, err
	}

	return &dto.ResponseWithID{ID: req.ID}, nil
}
