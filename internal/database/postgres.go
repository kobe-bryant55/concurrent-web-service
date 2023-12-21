package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
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
	var err error

	storeOnce.Do(func() {
		o := postgresStoreDefaultOpts()
		for _, fn := range opts {
			fn(&o)
		}

		conStr := fmt.Sprintf(`user=%s dbname=%s password=%s port=%s sslmode=disable`, o.user, o.name, o.password, o.port)

		db, err := sql.Open("postgres", conStr)
		if err != nil {
			err = fmt.Errorf("SQL open error: %w", err)
			return
		}

		if err := db.Ping(); err != nil {
			err = fmt.Errorf("SQL Connection error: %w", err)
			return
		}

		singleton = &postgresStore{DB: db}
	})

	return singleton, err
}

func (s postgresStore) GetInstance() *sql.DB {
	return s.DB
}

func (s postgresStore) InitDB() error {
	var err error
	initOnce.Do(func() {
		if err2 := s.createStatusEnums(); err != nil {
			err = err2
			return
		}
		if err2 := s.createTasksTable(); err != nil {
			err = err2
			return
		}

		log.Println("DB started successfully...")
	})

	return err
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
		user:     "development",
		name:     "development",
		password: "development",
		port:     "5432",
	}
}
