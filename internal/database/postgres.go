package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

var (
	singleton *postgresStore
	storeOnce sync.Once
	initOnce  sync.Once
)

type postgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(opts ...StoreOptsFunc) (SQLStore, error) {
	ch := make(chan error, 1)
	defer close(ch)

	storeOnce.Do(func() {
		o := postgresStoreDefaultOpts()
		for _, fn := range opts {
			fn(&o)
		}

		conStr := fmt.Sprintf(`user=%s dbname=%s password=%s port=%s sslmode=disable`, o.user, o.name, o.password, o.port)

		db, err := sql.Open("postgres", conStr)
		if err != nil {
			ch <- fmt.Errorf("SQL open error: %w", err)
			return
		}

		if err := db.Ping(); err != nil {
			ch <- fmt.Errorf("SQL Connection error: %w", err)
			return
		}

		ch <- nil
		singleton = &postgresStore{DB: db}
	})

	return singleton, <-ch
}

func (s postgresStore) GetInstance() *sql.DB {
	return s.DB
}

func (s postgresStore) InitDB() error {
	ch := make(chan error, 1)
	defer close(ch)

	initOnce.Do(func() {
		if err := s.createStatusEnums(); err != nil {
			ch <- err
			return
		}
		if err := s.createTasksTable(); err != nil {
			ch <- err
			return
		}

		ch <- nil
		log.Println("DB started successfully...")
	})

	return <-ch
}

func (s postgresStore) createStatusEnums() error {
	query := `DO $$ BEGIN
		IF to_regtype('status') IS NULL THEN
		CREATE TYPE status AS ENUM('active', 'passive');
		END IF;
		END $$;
	`

	_, err := s.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("tasks status enums creation failed: %w", err)
	}

	return nil
}

func (s postgresStore) createTasksTable() error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
    	id				serial PRIMARY KEY,
    	description 	varchar(500) NOT NULL, 
    	title 		   	varchar(21) NOT NULL, 
		status_state    status
	)
	`

	_, err := s.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("tasks table creation failed: %w", err)
	}

	return nil
}

func postgresStoreDefaultOpts() StoreOpts {
	return StoreOpts{
		user:     "local",
		name:     "local",
		password: "local",
		port:     "5432",
	}
}
