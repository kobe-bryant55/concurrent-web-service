package taskdomain

type IRepository interface {
	Create(task *Task) error
	Read(id uint64) (*Task, error)
	Reads() ([]*Task, error)
	Update(task *Task) error
	Delete(id uint64) error
}
