package taskdomain

import (
	"database/sql"

	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
)

type taskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) IRepository {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) Create(t *Task) error {
	query := `INSERT INTO tasks (description, title, status_state) VALUES ($1, $2, $3)`

	_, err := r.db.Query(query, t.GetDescription(), t.GetTitle(), t.GetStatus())
	if err != nil {
		return errorutils.New(errorutils.ErrTaskCreate, err)
	}

	return nil
}

func (r *taskRepository) Read(i uint64) (*Task, error) {
	rows, err := r.db.Query("SELECT * FROM tasks WHERE id = $1", i)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrTaskRead, err)
	}

	for rows.Next() {
		task, err := scanIntoTask(rows)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrTaskNotFound, err)
		}

		return task, nil
	}

	return nil, errorutils.New(errorutils.ErrTaskNotFound, errorutils.ErrTaskRead)
}

func (r *taskRepository) Reads() ([]*Task, error) {
	fq := `SELECT * FROM tasks`

	var tasks []*Task

	rows, err := r.db.Query(fq)
	if err != nil {
		return nil, errorutils.New(errorutils.ErrTaskReads, err)
	}

	for rows.Next() {
		t, err := scanIntoTask(rows)
		if err != nil {
			return nil, errorutils.New(errorutils.ErrTaskReads, err)
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *taskRepository) Update(t *Task) error {
	_, err := r.db.Query("UPDATE tasks SET description = $1, title = $2, status_state = $3 WHERE id = $4;", t.GetDescription(), t.GetTitle(), t.GetStatus(), t.GetID())
	if err != nil {
		return errorutils.New(errorutils.ErrTaskUpdate, err)
	}

	return nil
}

func (r *taskRepository) Delete(i uint64) error {
	_, err := r.db.Query("DELETE FROM tasks WHERE id = $1", i)
	if err != nil {
		return errorutils.New(errorutils.ErrTaskDelete, err)
	}

	return nil
}

func scanIntoTask(rows *sql.Rows) (*Task, error) {
	t := new(Task)

	err := rows.Scan(&t.id, &t.description, &t.title, &t.status)

	return t, err
}
