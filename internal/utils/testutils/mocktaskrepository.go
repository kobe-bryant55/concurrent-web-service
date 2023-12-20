package testutils

import (
	taskdomain "github.com/MehmetTalhaSeker/concurrent-web-service/domain/task"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the IRepository interface for testing purposes.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(task *taskdomain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRepository) Read(id uint64) (*taskdomain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*taskdomain.Task), args.Error(1)
}

func (m *MockRepository) Reads() ([]*taskdomain.Task, error) {
	args := m.Called()
	return args.Get(0).([]*taskdomain.Task), args.Error(1)
}

func (m *MockRepository) Update(task *taskdomain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}
