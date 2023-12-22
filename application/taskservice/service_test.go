package taskservice_test

import (
	"testing"

	"github.com/MehmetTalhaSeker/concurrent-web-service/application/taskservice"
	taskdomain "github.com/MehmetTalhaSeker/concurrent-web-service/domain/task"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Create(t *testing.T) {
	repo := new(testutils.MockRepository)
	service := taskservice.NewService(repo)

	request := &dto.TaskCreateRequest{
		Title:       "Test Task",
		Description: "Test Description",
	}

	repo.On("Create", mock.Anything).Return(nil)

	err := service.Create(request)

	repo.AssertCalled(t, "Create", mock.Anything)

	assert.NoError(t, err)
}

func TestService_Read(t *testing.T) {
	repo := new(testutils.MockRepository)
	service := taskservice.NewService(repo)

	request := &dto.RequestWithID{
		ID: "123",
	}

	repo.On("Read", uint64(123)).Return(&taskdomain.Task{}, nil)

	_, err := service.Read(request)

	repo.AssertCalled(t, "Read", uint64(123))

	assert.NoError(t, err)
}

func TestService_Reads(t *testing.T) {
	repo := new(testutils.MockRepository)
	service := taskservice.NewService(repo)

	repo.On("Reads").Return([]*taskdomain.Task{}, nil)

	_, err := service.Reads()

	repo.AssertCalled(t, "Reads")

	assert.NoError(t, err)
}

func TestService_Update(t *testing.T) {
	repo := new(testutils.MockRepository)
	service := taskservice.NewService(repo)

	request := &dto.TaskUpdateRequest{
		ID:     "123",
		Status: types.Active,
	}

	mt, _ := taskdomain.NewTaskBuilder().SetTitle("title").SetDescription("description").SetStatus(types.Passive).Build()
	repo.On("Read", uint64(123)).Return(mt, nil)
	repo.On("Update", mock.Anything).Return(nil)

	_, err := service.Update(request)

	repo.AssertCalled(t, "Read", uint64(123))
	repo.AssertCalled(t, "Update", mock.Anything)

	assert.NoError(t, err)
}

func TestService_Delete(t *testing.T) {
	repo := new(testutils.MockRepository)
	service := taskservice.NewService(repo)

	request := &dto.RequestWithID{
		ID: "123",
	}

	repo.On("Read", uint64(123)).Return(&taskdomain.Task{}, nil)
	repo.On("Delete", uint64(123)).Return(nil)

	_, err := service.Delete(request)

	repo.AssertCalled(t, "Read", uint64(123))
	repo.AssertCalled(t, "Delete", uint64(123))

	assert.NoError(t, err)
}
