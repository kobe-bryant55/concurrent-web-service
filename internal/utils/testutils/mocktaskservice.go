package testutils

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/dto"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) Create(req *dto.TaskCreateRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockService) Read(req *dto.RequestWithID) (*dto.TaskResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.TaskResponse), args.Error(1)
}

func (m *MockService) Reads() ([]*dto.TaskResponse, error) {
	args := m.Called()
	return args.Get(0).([]*dto.TaskResponse), args.Error(1)
}

func (m *MockService) Update(req *dto.TaskUpdateRequest) (*dto.TaskResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.TaskResponse), args.Error(1)
}

func (m *MockService) Delete(req *dto.RequestWithID) (*dto.ResponseWithID, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.ResponseWithID), args.Error(1)
}
